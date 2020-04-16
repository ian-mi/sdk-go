package transformer

import (
	"time"

	"github.com/ian-mi/sdk-go/v2/binding"
	"github.com/ian-mi/sdk-go/v2/binding/spec"
	"github.com/ian-mi/sdk-go/v2/types"
)

var (
	// Add the cloudevents time attribute, if missing, to time.Now()
	AddTimeNow binding.Transformer = addTimeNow{}
)

type addTimeNow struct{}

func (a addTimeNow) Transform(reader binding.MessageMetadataReader, writer binding.MessageMetadataWriter) error {
	attr, ti := reader.GetAttribute(spec.Time)
	if ti == nil {
		return writer.SetAttribute(attr, types.Timestamp{Time: time.Now()})
	}
	return nil
}
