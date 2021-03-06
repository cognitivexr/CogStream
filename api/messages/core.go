package messages

import (
	"encoding/json"
)

type OperationCode string

const (
	OperationExpose  OperationCode = "expose"
	OperationRecord  OperationCode = "record"
	OperationAnalyze OperationCode = "analyze"
)

type MessageType int32

const (
	MessageTypeAlert MessageType = iota + 1
	MessageTypeOperationSpec
	MessageTypeEngineFormatSpec
	MessageTypeClientFormatSpec
	MessageTypeStreamSpec
)

type AlertCode int32

type Attributes map[string][]string

func (a Attributes) Set(key string, value string) {
	a[key] = []string{value}
}

func (a Attributes) Get(key string) string {
	if a == nil {
		return ""
	}
	v := a[key]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

type Message struct {
	Type    MessageType     `json:"type"`
	Content json.RawMessage `json:"content"`
}

type Alert struct {
	AlertCode `json:"alertCode"`
}

type OperationSpec struct {
	Code       OperationCode `json:"code"`
	Attributes Attributes    `json:"attributes"`
}

type EngineFormatSpec struct {
	Attributes Attributes `json:"attributes"`
}

type ClientFormatSpec struct {
	Attributes Attributes `json:"attributes"`
}

type EngineAddress string

type StreamSpec struct {
	EngineAddress EngineAddress `json:"engineAddress"`
	Attributes    Attributes    `json:"attributes"`
}
