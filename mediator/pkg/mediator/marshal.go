package mediator

import (
	"cognitivexr.at/cogstream/api/messages"
	"encoding/json"
	"fmt"
)

func UnmarshalOperationSpec(message messages.Message) (messages.OperationSpec, error) {
	if message.Type != messages.MessageTypeOperationSpec {
		return messages.OperationSpec{}, fmt.Errorf("message is not of type %d (%s)", messages.MessageTypeOperationSpec, "OperationSpec")
	}
	var spec messages.OperationSpec
	err := json.Unmarshal(message.Content, &spec)
	if err != nil {
		return messages.OperationSpec{}, fmt.Errorf("cannot unmarshal message: %v", err)
	}
	return spec, nil
}

func UnmarshalClientFormatSpec(message messages.Message) (messages.ClientFormatSpec, error) {
	if message.Type != messages.MessageTypeClientFormatSpec {
		return messages.ClientFormatSpec{}, fmt.Errorf("message is not of type %d (%s)", messages.MessageTypeClientFormatSpec, "OperationSpec")
	}
	var spec messages.ClientFormatSpec
	err := json.Unmarshal(message.Content, &spec)
	if err != nil {
		return messages.ClientFormatSpec{}, fmt.Errorf("cannot unmarshal message: %v", err)
	}
	return spec, nil
}

func MarshalEngineFormatSpec(spec *messages.EngineFormatSpec) ([]byte, error) {
	content, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal engine format spec: %v", err)
	}
	return wrapInMessage(messages.MessageTypeEngineFormatSpec, string(content)), nil
}

func MarshalStreamSpec(spec *messages.StreamSpec) ([]byte, error) {
	content, err := json.Marshal(spec)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal engine format spec: %v", err)
	}
	return wrapInMessage(messages.MessageTypeStreamSpec, string(content)), nil
}

func MarshalAlert(alert *messages.Alert) ([]byte, error) {
	content, err := json.Marshal(alert)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal engine format spec: %v", err)
	}
	return wrapInMessage(messages.MessageTypeAlert, string(content)), nil
}

func wrapInMessage(type_ messages.MessageType, content string) []byte {
	return []byte(fmt.Sprintf("{\"type\":%d,\"content\":%s}", type_, content))
}
