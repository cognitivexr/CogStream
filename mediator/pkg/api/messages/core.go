package messages

import (
	"encoding/json"
	"fmt"
	"io"
)

type MessageCode int

// step 1
const (
	CodeExpose MessageCode = iota + 1
	CodeRecord
	CodeAnalyze
)

// step 2
const (
	CodeConstraints MessageCode = (iota + 1) * 10
)

// step 3
const (
	CodeFormat MessageCode = (iota + 1) * 100
)

// step 4
const (
	CodeExposeAgreement MessageCode = (iota + 1) * 1000
	CodeRecordAgreement
	CodeAnalyzeAgreement
)

type Messages struct {
	messages []Message
}

func NewMessages() *Messages {
	return &Messages{messages: make([]Message, 4)}
}

func (m *Messages) Add(message Message) {
	m.messages = append(m.messages, message)
}

type Message interface {
	GetCode() MessageCode
}

type Expose struct {
}

type Record struct {
	Format string `json:"format"`
	Name   string `json:"name"`
}

type Analyze struct {
	Model string
	// this is where a proper query for a model would go, like we defined in
	// https://www.usenix.org/system/files/hotedge19-paper-rausch.pdf Listing 3
}

type Constraints struct {
}

type Format struct {
	Resolution string
}

type ExposeAgreement struct {
	Format Format `json:"format"`
	URI    string `json:"uri"`
}

type RecordAgreement struct {
	Format Format `json:"format"`
	URI    string `json:"uri"`
}

type AnalyzeAgreement struct {
	Format Format `json:"format"`
}

type Envelope struct {
	Type MessageCode     `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TODO maybe prettier
func UnmarshalJSONMessage(code MessageCode, data io.Reader) (Message, error) {
	decoder := json.NewDecoder(data)
	switch code {
	case CodeRecord:
		var decoded Record
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeAnalyze:
		var decoded Analyze
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeExpose:
		var decoded Expose
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeConstraints:
		var decoded Constraints
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeFormat:
		var decoded Format
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeRecordAgreement:
		var decoded RecordAgreement
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeAnalyzeAgreement:
		var decoded AnalyzeAgreement
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	case CodeExposeAgreement:
		var decoded ExposeAgreement
		err := decoder.Decode(&decoded)
		if err != nil {
			return nil, fmt.Errorf("could not parse message: %v", err)
		}
		return &decoded, nil
	default:
		return nil, fmt.Errorf("could not parse message %v with code %v", data, code)
	}
}

func (e *Expose) GetCode() MessageCode {
	return CodeExpose
}

func (e *Record) GetCode() MessageCode {
	return CodeRecord
}

func (e *Analyze) GetCode() MessageCode {
	return CodeAnalyze
}

func (e *Constraints) GetCode() MessageCode {
	return CodeConstraints
}

func (e *Format) GetCode() MessageCode {
	return CodeFormat
}

func (e *ExposeAgreement) GetCode() MessageCode {
	return CodeExposeAgreement
}

func (e *RecordAgreement) GetCode() MessageCode {
	return CodeRecordAgreement
}

func (e *AnalyzeAgreement) GetCode() MessageCode {
	return CodeAnalyzeAgreement
}
