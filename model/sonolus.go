package model

type Intermediate struct {
	Archetype string                 `json:"archetype"`
	Data      map[string]interface{} `json:"data"` // Can be float64 or Intermediate
	Sim       bool                   `json:"sim"`
}

type SonolusLevelDataEntityData struct {
	Name  string   `json:"name"`
	Value *float64 `json:"value,omitempty"` // Can be float64 or string
	Ref   string   `json:"ref,omitempty"`
}

type SonolusLevelDataEntity struct {
	Archetype string                       `json:"archetype"`
	Data      []SonolusLevelDataEntityData `json:"data"`
	Ref       string                       `json:"ref,omitempty"`
}

type SonolusLevelData struct {
	BgmOffset float64                   `json:"bgmOffset"`
	Entities  []*SonolusLevelDataEntity `json:"entities"`
}
