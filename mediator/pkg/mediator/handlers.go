package mediator

import (
	"fmt"
)

func DefaultOperationHandler(hs *HandshakeContext, platform Platform) error {
	engineFormat, err := platform.GetEngineFormatSpec(hs)
	if err != nil {
		return fmt.Errorf("cannot get engine format spec: %v", err)
	}
	hs.EngineFormatSpec = engineFormat
	return nil
}

func DefaultFormatHandler(hs *HandshakeContext, platform Platform) error {
	streamSpec, err := platform.GetStreamSpec(hs)
	if err != nil {
		return fmt.Errorf("cannot get engine format spec: %v", err)
	}
	hs.StreamSpec = streamSpec
	return nil
}
