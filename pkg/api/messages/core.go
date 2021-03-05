package messages

const (
	// step 1
	CodeExpose  = 1
	CodeRecord  = 2
	CodeAnalyze = 3

	// step 2
	CodeConstraints = 10

	// step 3
	CodeFormat = 100

	// step 4
	CodeExposeAgreement  = 1000
	CodeRecordAgreement  = 2000
	CodeAnalyzeAgreement = 3000
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
	GetCode() int
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

func (e *Expose) GetCode() int {
	return CodeExpose
}

func (e *Record) GetCode() int {
	return CodeRecord
}

func (e *Analyze) GetCode() int {
	return CodeAnalyze
}

func (e *Constraints) GetCode() int {
	return CodeConstraints
}

func (e *Format) GetCode() int {
	return CodeFormat
}

func (e *ExposeAgreement) GetCode() int {
	return CodeExposeAgreement
}

func (e *RecordAgreement) GetCode() int {
	return CodeRecordAgreement
}

func (e *AnalyzeAgreement) GetCode() int {
	return CodeAnalyzeAgreement
}
