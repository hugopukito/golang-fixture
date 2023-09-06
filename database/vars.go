package database

var GoSQLTypeMap = map[string]string{
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
	"time.Time": "TIMESTAMP",
	"uuid.UUID": "VARCHAR(36)",
}

var RandomWords = []string{
	"Lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "Ut", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "ut", "aliquip", "ex", "ea", "commodo",
	"consequat", "Duis", "aute", "irure", "dolor", "in", "reprehenderit", "in", "voluptate",
	"velit", "esse", "cillum", "dolore", "eu", "fugiat", "nulla", "pariatur", "Excepteur",
	"sint", "occaecat", "cupidatat", "non", "proident", "sunt", "in", "culpa", "qui",
	"officia", "deserunt", "mollit", "anim", "id", "est", "laborum",
}
