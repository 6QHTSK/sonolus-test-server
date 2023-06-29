package model

import (
	"context"
	"encoding/json"
)

type BestdoriChart []BestdoriObject

func (chart *BestdoriChart) UnmarshalJSON(data []byte) error {
	var temp []map[string]interface{}
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	for _, item := range temp {
		itemType, ok := item["type"].(string)
		if !ok {
			continue
		}
		itemByte, err := json.Marshal(item)
		if err != nil {
			return err
		}
		switch itemType {
		case "Single":
			var note BestdoriNote
			err := json.Unmarshal(itemByte, &note)
			if err != nil {
				return err
			}
			*chart = append(*chart, &note)
		case "Directional":
			var note BestdoriDirectioalNote
			err := json.Unmarshal(itemByte, &note)
			if err != nil {
				return err
			}
			*chart = append(*chart, &note)
		case "Slide":
			var note BestdoriSlideNote
			err := json.Unmarshal(itemByte, &note)
			if err != nil {
				return err
			}
			*chart = append(*chart, &note)
		case "Long":
			var note BestdoriLongNote
			err := json.Unmarshal(itemByte, &note)
			if err != nil {
				return err
			}
			*chart = append(*chart, &note)
		case "BPM":
			var note BestdoriBpmObject
			err := json.Unmarshal(itemByte, &note)
			if err != nil {
				return err
			}
			*chart = append(*chart, &note)
		}
	}
	return nil
}

type BestdoriObject interface {
	getType() string
	Convert(ctx context.Context) error
}

type BaseBestdoriObject struct {
	Beat float64 `json:"beat"`
}

type BaseBestdoriNote struct {
	BaseBestdoriObject
	Lane  float64 `json:"lane"`
	Flick bool    `json:"flick,omitempty"`
}

type BestdoriConnectionNote struct {
	BaseBestdoriNote
	Hidden bool `json:"hidden"`
}

type BestdoriNote struct {
	BaseBestdoriNote
}

func (note *BestdoriNote) getType() string {
	return "Single"
}

type BestdoriDirectioalNote struct {
	BaseBestdoriNote
	Directional string `json:"directional"`
	Width       int    `json:"width"`
}

func (note *BestdoriDirectioalNote) getType() string {
	return "Directional"
}

type BestdoriSlideNote struct {
	Connections []BestdoriConnectionNote `json:"connections"`
}

func (note *BestdoriSlideNote) getType() string {
	return "Slide"
}

type BestdoriLongNote struct {
	Connections []BestdoriConnectionNote `json:"connections"`
}

func (note *BestdoriLongNote) getType() string {
	return "Long"
}

type BestdoriBpmObject struct {
	BaseBestdoriObject
	Bpm float64 `json:"bpm"`
}

func (note *BestdoriBpmObject) getType() string {
	return "BPM"
}
