CREATE TABLE A (
    PK INT64 NOT NULL,
    Col_01 BOOL,
    Col_02 BOOL NOT NULL,
    Col_03 BYTES(50),
    Col_04 BYTES(50) NOT NULL,
    Col_05 DATE,
    Col_06 DATE NOT NULL,
    Col_07 FLOAT64,
    Col_08 FLOAT64 NOT NULL,
    Col_09 INT64,
    Col_10 INT64 NOT NULL,
    Col_11 JSON,
    Col_12 JSON NOT NULL,
    Col_13 NUMERIC,
    Col_14 NUMERIC NOT NULL,
    Col_15 STRING(50),
    Col_16 STRING(50) NOT NULL,
    Col_17 TIMESTAMP,
    Col_18 TIMESTAMP NOT NULL,
) PRIMARY KEY (PK);

INSERT INTO A (PK, Col_01, Col_02, Col_03, Col_04, Col_05, Col_06, Col_07, Col_08, Col_09, Col_10, Col_11, Col_12, Col_13, Col_14, Col_15, Col_16, Col_17, Col_18) VALUES
(1, true, true, b"abc", b"abc", "2025-06-04", "2025-06-04", -123.45, -123.45, 1, 1, JSON "{}", JSON "{}", NUMERIC "10", NUMERIC "10", "abc", "abc", "2025-06-04T12:34:56Z","2025-06-04T12:34:56Z"),
(2, true, true, b"abc", b"abc", "2025-06-04", "2025-06-04", -123.45, -123.45, 1, 1, JSON "{}", JSON "{}", NUMERIC "10", NUMERIC "10", "abc", "abc", "2025-06-04T12:34:56Z","2025-06-04T12:34:56Z"),
(3, true, true, b"abc", b"abc", "2025-06-04", "2025-06-04", -123.45, -123.45, 1, 1, JSON "{}", JSON "{}", NUMERIC "10", NUMERIC "10", "abc", "abc", "2025-06-04T12:34:56Z","2025-06-04T12:34:56Z"),
(4, true, true, b"abc", b"abc", "2025-06-04", "2025-06-04", -123.45, -123.45, 1, 1, JSON "{}", JSON "{}", NUMERIC "10", NUMERIC "10", "abc", "abc", "2025-06-04T12:34:56Z","2025-06-04T12:34:56Z");
