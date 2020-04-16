package amqp

import (
	"context"
	"io"
	"io/ioutil"

	"github.com/Azure/go-amqp"

	"github.com/ian-mi/sdk-go/v2/binding"
	"github.com/ian-mi/sdk-go/v2/binding/format"
	"github.com/ian-mi/sdk-go/v2/binding/spec"
	"github.com/ian-mi/sdk-go/v2/types"
)

// WriteMessage fills the provided amqpMessage with the message m.
// Using context you can tweak the encoding processing (more details on binding.Write documentation).
func WriteMessage(ctx context.Context, m binding.Message, amqpMessage *amqp.Message, transformers ...binding.Transformer) error {
	structuredWriter := (*amqpMessageWriter)(amqpMessage)
	binaryWriter := (*amqpMessageWriter)(amqpMessage)

	_, err := binding.Write(
		ctx,
		m,
		structuredWriter,
		binaryWriter,
		transformers...,
	)
	return err
}

type amqpMessageWriter amqp.Message

func (b *amqpMessageWriter) SetStructuredEvent(ctx context.Context, format format.Format, event io.Reader) error {
	val, err := ioutil.ReadAll(event)
	if err != nil {
		return err
	}
	b.Data = [][]byte{val}
	b.Properties = &amqp.MessageProperties{ContentType: format.MediaType()}
	return nil
}

func (b *amqpMessageWriter) Start(ctx context.Context) error {
	b.Properties = &amqp.MessageProperties{}
	b.ApplicationProperties = make(map[string]interface{})
	return nil
}

func (b *amqpMessageWriter) End(ctx context.Context) error {
	return nil
}

func (b *amqpMessageWriter) SetData(reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	b.Data = [][]byte{data}
	return nil
}

func (b *amqpMessageWriter) SetAttribute(attribute spec.Attribute, value interface{}) error {
	if attribute.Kind() == spec.DataContentType {
		if value == nil {
			b.Properties.ContentType = ""
			return nil
		}
		s, err := types.Format(value)
		if err != nil {
			return err
		}
		b.Properties.ContentType = s
	} else {
		if value == nil {
			delete(b.ApplicationProperties, prefix+attribute.Name())
			return nil
		}
		v, err := safeAMQPPropertiesUnwrap(value)
		if err != nil {
			return err
		}
		b.ApplicationProperties[prefix+attribute.Name()] = v
	}
	return nil
}

func (b *amqpMessageWriter) SetExtension(name string, value interface{}) error {
	v, err := safeAMQPPropertiesUnwrap(value)
	if err != nil {
		return err
	}
	b.ApplicationProperties[prefix+name] = v
	return nil
}

var _ binding.BinaryWriter = (*amqpMessageWriter)(nil)     // Test it conforms to the interface
var _ binding.StructuredWriter = (*amqpMessageWriter)(nil) // Test it conforms to the interface
