package funcs

type Fixture struct {
	Entities map[string]Entity `yaml:",inline"`
}

type Entity map[string]map[string]interface{}
