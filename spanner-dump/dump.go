//
// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package spanner_dump

import (
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"io"
	"regexp"
	"strings"
	"time"

	adminapi "cloud.google.com/go/spanner/admin/database/apiv1"
	adminpb "cloud.google.com/go/spanner/admin/database/apiv1/databasepb"

	schenerate_spanner "github.com/Jumpaku/schenerate/spanner"
)

// This is an ad hoc value, but considering mutations limit (20,000),
// 100 rows/statement would be safe in most cases.
// https://cloud.google.com/spanner/quotas#limits_for_creating_reading_updating_and_deleting_data
const defaultBulkSize = 100

// Dumper is a dumper to export a database.
type Dumper struct {
	project   string
	instance  string
	database  string
	query     map[string]string
	out       io.Writer
	timestamp *time.Time
	bulkSize  uint
	upsert    bool
	tables    []string

	client      *spanner.Client
	adminClient *adminapi.DatabaseAdminClient
}

// NewDumper creates Dumper with specified configurations.
func NewDumper(ctx context.Context, project, instance, database string, out io.Writer, timestamp *time.Time, bulkSize uint, query map[string]string, sort bool, upsert bool) (*Dumper, error) {
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instance, database)
	client, err := spanner.NewClientWithConfig(ctx, dbPath, spanner.ClientConfig{
		SessionPoolConfig: spanner.SessionPoolConfig{
			MinOpened: 1,
			MaxOpened: 1,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create spanner client: %v", err)
	}

	adminClient, err := adminapi.NewDatabaseAdminClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create spanner admin client: %v", err)
	}

	if bulkSize == 0 {
		bulkSize = defaultBulkSize
	}

	tables := []string{}
	dumperQuery := map[string]string{}
	for table, where := range query {
		t := strings.Trim(table, "`")
		dumperQuery[t] = where
		tables = append(tables, t)
	}
	if sort {
		q, err := schenerate_spanner.Open(ctx, project, instance, database)
		if err != nil {
			return nil, fmt.Errorf("failed to open schema: %v", err)
		}
		s, err := schenerate_spanner.ListSchemas(ctx, q, tables)
		if err != nil {
			return nil, fmt.Errorf("failed to list schemas: %v", err)
		}
		g := s.BuildGraph()
		idx, cyclic := g.TopologicalSort()
		if cyclic {
			return nil, fmt.Errorf("cyclic dependency detected in tables")
		}
		var ts []string
		for _, i := range idx {
			ts = append(ts, strings.Trim(g.Get(i).Name, "`"))
		}
		// reverse tables
		for i, j := 0, len(ts)-1; i < j; i, j = i+1, j-1 {
			ts[i], ts[j] = ts[j], ts[i]
		}
		tables = ts

	}

	d := &Dumper{
		project:     project,
		instance:    instance,
		database:    database,
		query:       dumperQuery,
		tables:      tables,
		out:         out,
		bulkSize:    bulkSize,
		timestamp:   timestamp,
		upsert:      upsert,
		client:      client,
		adminClient: adminClient,
	}

	return d, nil
}

// Cleanup cleans up hold resources.
func (d *Dumper) Cleanup() {
	d.client.Close()
	d.adminClient.Close()
}

// DumpDDLs dumps all DDLs in the database.
func (d *Dumper) DumpDDLs(ctx context.Context) error {
	dbPath := fmt.Sprintf("projects/%s/instances/%s/databases/%s", d.project, d.instance, d.database)
	resp, err := d.adminClient.GetDatabaseDdl(ctx, &adminpb.GetDatabaseDdlRequest{
		Database: dbPath,
	})
	if err != nil {
		return err
	}

	for _, ddl := range resp.Statements {
		tableName := strings.Trim(parseTableNameFromDDL(ddl), "`")
		if _, ok := d.query[tableName]; len(d.query) > 0 && !ok {
			continue
		}
		fmt.Fprintf(d.out, "%s;\n", ddl)
	}

	return nil
}

func parseTableNameFromDDL(ddl string) string {
	if indexRegexp.MatchString(ddl) {
		match := indexRegexp.FindStringSubmatch(ddl)
		return match[1]
	}
	if tableRegexp.MatchString(ddl) {
		match := tableRegexp.FindStringSubmatch(ddl)
		return match[1]
	}
	if alterRegexp.MatchString(ddl) {
		match := alterRegexp.FindStringSubmatch(ddl)
		return match[1]
	}
	return ""
}

var indexRegexp = regexp.MustCompile("^\\s*CREATE\\s+(?:UNIQUE\\s+|NULL_FILTERED\\s+)?INDEX\\s+(?:[a-zA-Z0-9_`]+)\\s+ON\\s+`?([a-zA-Z0-9_]+)`?")
var tableRegexp = regexp.MustCompile("^\\s*CREATE\\s+TABLE\\s+`?([a-zA-Z0-9_]+)`?")
var alterRegexp = regexp.MustCompile("^\\s*ALTER\\s+TABLE\\s+`?([a-zA-Z0-9_]+)`?")

// DumpTables dumps all table records in the database.
func (d *Dumper) DumpTables(ctx context.Context) error {
	txn := d.client.ReadOnlyTransaction()
	if d.timestamp != nil {
		txn = txn.WithTimestampBound(spanner.ReadTimestamp(*d.timestamp))
	}
	defer txn.Close()

	tables, err := FetchTables(ctx, txn, d.tables)
	if err != nil {
		return fmt.Errorf("failed to fetch tables: %v", err)
	}
	for _, t := range tables {
		if err := d.dumpTable(ctx, t, txn); err != nil {
			return fmt.Errorf("failed to dump table %s: %v", t.Name, err)
		}
	}
	return nil
}

func (d *Dumper) dumpTable(ctx context.Context, table *Table, txn *spanner.ReadOnlyTransaction) error {
	queryCondition := d.query[table.Name]
	if queryCondition == "" {
		queryCondition = "TRUE"
	}
	stmt := fmt.Sprintf("SELECT %s FROM `%s` WHERE %s", table.quotedColumnList(), table.Name, queryCondition)
	iter := txn.Query(ctx, spanner.NewStatement(stmt))
	defer iter.Stop()

	writer := NewBufferedWriter(table, d.out, d.bulkSize, d.upsert)
	defer writer.Flush()
	for {
		row, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return err
		}

		values, err := DecodeRow(row)
		if err != nil {
			return err
		}
		writer.Write(values)
	}

	return nil
}
