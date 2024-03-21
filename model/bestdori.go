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

func (chart BestdoriChart) MarshalJSON() ([]byte, error) {
	// 创建一个包含结构体字段的映射
	var fields []map[string]interface{}

	for _, note := range chart {
		// 获取结构体的 JSON 序列化结果
		var mapNote map[string]interface{}
		jsonData, err := json.Marshal(note)
		if err != nil {
			return nil, err
		}

		// 将原始 JSON 数据解码为一个 map
		if err := json.Unmarshal(jsonData, &mapNote); err != nil {
			return nil, err
		}

		// 添加type字段
		var structType string
		switch note.(type) {
		case *BestdoriNote:
			structType = "Single"
		case *BestdoriDirectioalNote:
			structType = "Directional"
		case *BestdoriSlideNote:
			structType = "Slide"
		case *BestdoriLongNote:
			structType = "Slide" // No Long Note in Fan-made Charts
		case *BestdoriBpmObject:
			structType = "BPM"
		}
		mapNote["type"] = structType

		// 重新序列化为 JSON
		fields = append(fields, mapNote)
	}
	return json.Marshal(fields)
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
	Hidden bool `json:"hidden,omitempty"`
}

type BestdoriNote struct {
	BaseBestdoriNote
}

func (note *BestdoriNote) getType() string {
	return "Single"
}

type BestdoriDirectioalNote struct {
	BaseBestdoriNote
	Direction string  `json:"direction"`
	Width     float64 `json:"width"`
}

func (note *BestdoriDirectioalNote) getType() string {
	return "Direction"
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
