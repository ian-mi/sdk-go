package binding_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ian-mi/sdk-go/v2/pkg/binding"
	"github.com/ian-mi/sdk-go/v2/pkg/event"
	"github.com/ian-mi/sdk-go/v2/pkg/types"
)

func TestWithFinish(t *testing.T) {
	var testEvent = event.Event{
		DataEncoded: []byte(`"data"`),
		Context: event.EventContextV1{
			DataContentType: event.StringOfApplicationJSON(),
			Source:          types.URIRef{URL: url.URL{Path: "source"}},
			ID:              "id",
			Type:            "type"}.AsV1(),
	}

	done := make(chan error, 1)
	m := binding.WithFinish((*binding.EventMessage)(&testEvent), func(err error) {
		done <- err
	})
	select {
	case <-done:
		assert.Fail(t, "done early")
	default:
	}
	ch := make(chan binding.Message, 1)
	assert.NoError(t, binding.ChanSender(ch).Send(context.Background(), m))
	assert.NoError(t, <-done)
}
