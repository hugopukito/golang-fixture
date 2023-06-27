package database

var goSQLTypeMap = map[string]string{
	"int":       "INT",
	"int8":      "TINYINT",
	"int16":     "SMALLINT",
	"int32":     "INT",
	"int64":     "BIGINT",
	"uint":      "INT UNSIGNED",
	"uint8":     "TINYINT UNSIGNED",
	"uint16":    "SMALLINT UNSIGNED",
	"uint32":    "INT UNSIGNED",
	"uint64":    "BIGINT UNSIGNED",
	"float32":   "FLOAT",
	"float64":   "DOUBLE",
	"bool":      "BOOL",
	"string":    "VARCHAR(255)",
	"time.Time": "DATETIME",
	"uuid.UUID": "VARCHAR(36)",
}
