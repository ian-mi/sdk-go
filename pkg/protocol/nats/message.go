package nats

import (
	"bytes"
	"context"

	"github.com/ian-mi/sdk-go/v2/pkg/binding"
	"github.com/ian-mi/sdk-go/v2/pkg/binding/format"
	"github.com/nats-io/nats.go"
)

const prefix = "cloudEvents:" // Name prefix for AMQP properties that hold CE attributes.

// Message implements binding.Message by wrapping an *nats.Msg.
// This message *can* be read several times safely
type Message struct {
	Msg      *nats.Msg
	encoding binding.Encoding
}

// Wrap an *nats.Msg in a binding.Message.
// The returned message *can* be read several times safely
func NewMessage(msg *nats.Msg) *Message {
	return &Message{Msg: msg, encoding: binding.EncodingStructured}
}

var _ binding.Message = (*Message)(nil)

func (m *Message) ReadEncoding() binding.Encoding {
	return m.encoding
}

func (m *Message) ReadStructured(ctx context.Context, encoder binding.StructuredWriter) error {
	return encoder.SetStructuredEvent(ctx, format.JSON, bytes.NewReader(m.Msg.Data))
}

func (m *Message) ReadBinary(ctx context.Context, encoder binding.BinaryWriter) error {
	return binding.ErrNotBinary
}

func (m *Message) Finish(err error) error {
	return nil
}
