package messages

import (
	"encoding/json"
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
