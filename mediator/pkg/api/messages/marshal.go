package messages

import (
	"encoding/json"
	"fmt"
	"io"
)

// TODO maybe prettier
func UnmarshalJSONMessage(code MessageCode, r io.Reader) (Message, error) {
	decoder := json.NewDecoder(r)
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
		return nil, fmt.Errorf("could not parse message %v with code %v", r, code)
	}
}

func WriteJSONMessage(msg Message, w io.WriteCloser) error {
	encoder := json.NewEncoder(w)
	defer w.Close()
	switch msg.GetCode() {
	case CodeConstraints:
		reply := msg.(*Constraints)
		err := encoder.Encode(reply)
		if err != nil {
			return fmt.Errorf("couldn't write message: %v", err)
		}
	case CodeRecordAgreement:
		reply := msg.(*RecordAgreement)
		err := encoder.Encode(reply)
		if err != nil {
			return fmt.Errorf("couldn't write message: %v", err)
		}
	default:
		return fmt.Errorf("couldn't establish message type for code: %v", msg.GetCode())
	}
	return nil
}
