package structs

type Fixture struct {
	Entities map[string]Entity `yaml:",inline"`
}

type Entity map[string]interface{}
