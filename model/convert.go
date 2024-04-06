package model

import (
	"context"
	"github.com/6qhtsk/sonolus-test-server/errors"
	"sort"
	"strconv"
)

// from sonolus-core
// ref: https://github.com/Sonolus/sonolus-core/blob/master/src/common/core/engine/engine-archetype-data-name.ts
const (
	ARCHETYPE_BEAT      = "#BEAT"
	ARCHETYPE_BPM       = "#BPM"
	ARCHETYPE_BPMCHANGE = "#BPM_CHANGE"
	ARCHTYPE_TIMESCALE  = "#TIMESCALE"
	ARCHTYPE_JUDGMENT   = "#JUDGMENT"
	ARCHTYPE_ACCURACY   = "#ACCURACY"
)

func (bdNote *BestdoriNote) Convert(ctx context.Context) error {
	var archetype string
	if bdNote.Flick {
		archetype = "FlickNote"
	} else {
		archetype = "TapNote"
	}
	return appendIntermediate(ctx, &Intermediate{
		Archetype: archetype,
		Data: map[string]interface{}{
			ARCHETYPE_BEAT: bdNote.Beat,
			"lane":         bdNote.Lane - 3.0,
		},
		Sim: true,
	})
}

func (bdNote *BestdoriDirectioalNote) Convert(ctx context.Context) error {
	var direction float64
	if bdNote.Direction == "Left" {
		direction = -1
	} else {
		direction = 1
	}
	return appendIntermediate(ctx, &Intermediate{
		Archetype: "DirectionalFlickNote",
		Data: map[string]interface{}{
			ARCHETYPE_BEAT: bdNote.Beat,
			"lane":         bdNote.Lane - 3.0,
			"direction":    direction,
			"size":         bdNote.Width,
		},
		Sim: true,
	})
}

func SlideLongConvertor(ctx context.Context, connections []BestdoriConnectionNote) error {
	var first, start, head *Intermediate
	var connectors []*Intermediate
	var appends []*Intermediate

	connectorArchetype := "StraightSlideConnector"
	for _, connection := range connections {
		if connection.Hidden {
			connectorArchetype = "CurvedSlideConnector"
			break
		}
	}

	for i, connection := range connections {
		if i == 0 { // Start
			start = &Intermediate{
				Archetype: "SlideStartNote",
				Data: map[string]interface{}{
					ARCHETYPE_BEAT: connection.Beat,
					"lane":         connection.Lane - 3.0,
				},
				Sim: true,
			}
			first = start
			head = start
			appends = append(appends, first)
		} else if i == len(connections)-1 { // Tail
			var archetype string
			if connection.Flick {
				archetype = "SlideEndFlickNote"
			} else {
				archetype = "SlideEndNote"
			}
			tail := &Intermediate{
				Archetype: archetype,
				Data: map[string]interface{}{
					ARCHETYPE_BEAT: connection.Beat,
					"lane":         connection.Lane - 3.0,
					"first":        first,
					"prev":         start,
				},
				Sim: true,
			}
			if connection.Flick {
				if len(connections) == 2 && head.Data["lane"].(float64) == tail.Data["lane"].(float64) {
					tail.Data["long"] = 1.0
				} else {
					tail.Data["long"] = 0.0
				}
			}
			appends = append(appends, tail)
			first.Data["last"] = tail
			connectors = append(connectors, &Intermediate{
				Archetype: connectorArchetype,
				Data: map[string]interface{}{
					"first": first,
					"start": start,
					"head":  head,
					"tail":  tail,
				},
				Sim: false,
			})
			for _, connector := range connectors {
				connector.Data["end"] = tail
				appends = append(appends, connector)
			}
			connectors = []*Intermediate{}
		} else if connection.Hidden { // 隐藏音符
			tail := &Intermediate{
				Archetype: "IgnoredNote",
				Data: map[string]interface{}{
					ARCHETYPE_BEAT: connection.Beat,
					"lane":         connection.Lane - 3.0,
				},
				Sim: false,
			}
			appends = append(appends, tail)
			connectors = append(connectors, &Intermediate{
				Archetype: connectorArchetype,
				Data: map[string]interface{}{
					"first": first,
					"start": start,
					"head":  head,
					"tail":  tail,
				},
				Sim: false,
			})
			head = tail
		} else { //普通节点音符
			tail := &Intermediate{
				Archetype: "SlideTickNote",
				Data: map[string]interface{}{
					ARCHETYPE_BEAT: connection.Beat,
					"lane":         connection.Lane - 3.0,
					"first":        first,
					"prev":         start,
				},
				Sim: false,
			}
			appends = append(appends, tail)
			connectors = append(connectors, &Intermediate{
				Archetype: connectorArchetype,
				Data: map[string]interface{}{
					"first": first,
					"start": start,
					"head":  head,
					"tail":  tail,
				},
				Sim: false,
			})
			for _, connector := range connectors {
				connector.Data["end"] = tail
				appends = append(appends, connector)
			}
			connectors = []*Intermediate{}
			start = tail
			head = tail
		}
	}
	for _, intermediate := range appends {
		err := appendIntermediate(ctx, intermediate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bdNote *BestdoriSlideNote) Convert(ctx context.Context) error {
	return SlideLongConvertor(ctx, bdNote.Connections)
}

func (bdNote *BestdoriLongNote) Convert(ctx context.Context) error {
	return SlideLongConvertor(ctx, bdNote.Connections)
}

func (bdNote *BestdoriBpmObject) Convert(ctx context.Context) error {
	return appendIntermediate(ctx, &Intermediate{
		Archetype: ARCHETYPE_BPMCHANGE,
		Data: map[string]interface{}{
			ARCHETYPE_BEAT: bdNote.Beat,
			ARCHETYPE_BPM:  bdNote.Bpm,
		},
		Sim: false,
	})
}

type convertContextValue struct {
	Entities             []*SonolusLevelDataEntity
	BeatToIntermediates  map[float64][]*Intermediate
	IntermediateToRef    map[*Intermediate]string
	IntermediateToEntity map[*Intermediate]*SonolusLevelDataEntity
	refCounter           int64
}

func getRef(ctx context.Context, intermediate *Intermediate) (ref string) {
	ctxValues := ctx.Value("values").(*convertContextValue)
	ref, ok := ctxValues.IntermediateToRef[intermediate]
	if !ok {
		ref = strconv.FormatInt(ctxValues.refCounter, 36)
		ctxValues.refCounter++
		ctxValues.IntermediateToRef[intermediate] = ref
		entity, ok := ctxValues.IntermediateToEntity[intermediate]
		if ok {
			entity.Name = ref
			ctxValues.IntermediateToEntity[intermediate] = entity
		}
	}
	return ref
}

func appendIntermediate(ctx context.Context, intermediate *Intermediate) error {
	ctxValues := ctx.Value("values").(*convertContextValue)
	entity := SonolusLevelDataEntity{
		Archetype: intermediate.Archetype,
		Data:      []SonolusLevelDataEntityData{},
	}

	if intermediate.Sim {
		data := (*intermediate).Data
		beat, ok := data[ARCHETYPE_BEAT].(float64)
		if !ok {
			return errors.ConvertorUnexpectedBeat
		}
		intermediates, ok := ctxValues.BeatToIntermediates[beat]
		if ok {
			ctxValues.BeatToIntermediates[beat] = append(intermediates, intermediate)
		} else {
			ctxValues.BeatToIntermediates[beat] = []*Intermediate{intermediate}
		}
	}

	ref, ok := ctxValues.IntermediateToRef[intermediate]
	if ok {
		entity.Name = ref
	}

	ctxValues.IntermediateToEntity[intermediate] = &entity
	ctxValues.Entities = append(ctxValues.Entities, &entity)

	intermediateDataKeys := make([]string, 0, len((*intermediate).Data))
	for k := range (*intermediate).Data {
		intermediateDataKeys = append(intermediateDataKeys, k)
	}
	sort.Strings(intermediateDataKeys)

	for _, key := range intermediateDataKeys {
		valueNumber, ok := (*intermediate).Data[key].(float64)
		if ok {
			entity.Data = append(entity.Data, SonolusLevelDataEntityData{
				Name:  key,
				Value: &valueNumber,
			})
		} else {
			valuePIntermediate, ok := (*intermediate).Data[key].(*Intermediate)
			if ok {
				entity.Data = append(entity.Data, SonolusLevelDataEntityData{
					Name: key,
					Ref:  getRef(ctx, valuePIntermediate),
				})
			}
		}
	}
	return nil
}

func repairSlideNote(connection []BestdoriConnectionNote) BestdoriObject {
	sort.Slice(connection, func(i, j int) bool {
		return connection[i].Beat < connection[j].Beat
	})
	var visibleConnection []BestdoriConnectionNote
	for _, connectionItem := range connection {
		if !connectionItem.Hidden {
			visibleConnection = append(visibleConnection, connectionItem)
		}
	}
	switch len(visibleConnection) {
	case 0:
		return nil
	case 1:
		return &BestdoriNote{
			BaseBestdoriNote: BaseBestdoriNote{
				BaseBestdoriObject: BaseBestdoriObject{
					Beat: visibleConnection[0].Beat,
				},
				Lane:  visibleConnection[0].Lane,
				Flick: visibleConnection[0].Flick,
			},
		}
	default:
		// 删除之前的hidden音符
		start := 0
		for connection[start].Hidden {
			start++
		}
		end := len(connection)
		for connection[end-1].Hidden {
			end--
		}
		return &BestdoriSlideNote{
			Connections: connection[start:end],
		}
	}
}

func repair(chart *BestdoriChart) (repairedChart BestdoriChart) {
	for _, objectInterface := range *chart {
		switch object := objectInterface.(type) {
		case *BestdoriLongNote:
			slideObject := repairSlideNote(object.Connections)
			repairedChart = append(repairedChart, slideObject)
		case *BestdoriSlideNote:
			slideObject := repairSlideNote(object.Connections)
			repairedChart = append(repairedChart, slideObject)
		default:
			repairedChart = append(repairedChart, object)
		}
	}
	return repairedChart
}

func (bdChart *BestdoriChart) ConvertToSonnolus() (levelData SonolusLevelData, err error) {
	ctxValues := convertContextValue{
		Entities:             []*SonolusLevelDataEntity{},
		BeatToIntermediates:  map[float64][]*Intermediate{},
		IntermediateToRef:    map[*Intermediate]string{},
		IntermediateToEntity: map[*Intermediate]*SonolusLevelDataEntity{},
		refCounter:           0,
	}
	ctx := context.WithValue(context.Background(), "values", &ctxValues)
	err = appendIntermediate(ctx, &Intermediate{
		Archetype: "Initialization",
		Data:      map[string]interface{}{},
		Sim:       false,
	})
	if err != nil {
		return levelData, err
	}
	err = appendIntermediate(ctx, &Intermediate{
		Archetype: "Stage",
		Data:      map[string]interface{}{},
		Sim:       false,
	})
	if err != nil {
		return levelData, err
	}

	objects := repair(bdChart)

	for _, object := range objects {
		err := object.Convert(ctx)
		if err != nil {
			return SonolusLevelData{}, err
		}
	}

	beatToIntermediatesKeys := make([]float64, 0, len(ctxValues.BeatToIntermediates))
	for k := range ctxValues.BeatToIntermediates {
		beatToIntermediatesKeys = append(beatToIntermediatesKeys, k)
	}

	sort.Float64s(beatToIntermediatesKeys)

	for _, beat := range beatToIntermediatesKeys {
		intermediate := ctxValues.BeatToIntermediates[beat]
		for i := 1; i < len(intermediate); i++ {
			err := appendIntermediate(ctx, &Intermediate{
				Archetype: "SimLine",
				Data: map[string]interface{}{
					"a": intermediate[i-1],
					"b": intermediate[i],
				},
				Sim: false,
			})
			if err != nil {
				return SonolusLevelData{}, err
			}
		}
	}

	return SonolusLevelData{
		BgmOffset: 0,
		Entities:  ctxValues.Entities,
	}, nil
}
