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
	"testing"
)

func TestQuotedColumnList(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		table *Table
		want  string
	}{
		{
			desc:  "No columns",
			table: &Table{Columns: []string{}},
			want:  "",
		},
		{
			desc:  "Single column",
			table: &Table{Columns: []string{"C1"}},
			want:  "`C1`",
		},
		{
			desc:  "Multiple columns",
			table: &Table{Columns: []string{"C1", "C2"}},
			want:  "`C1`, `C2`",
		},
	} {
		t.Run(tt.desc, func(t *testing.T) {
			if got := tt.table.quotedColumnList(); got != tt.want {
				t.Errorf("quotedColumnList() of %v: got = %v, want = %v", tt.table, got, tt.want)
			}
		})
	}

}
