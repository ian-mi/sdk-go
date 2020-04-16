// Package test provides re-usable functions for binding tests.
package test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ian-mi/sdk-go/v2/pkg/binding"
	"github.com/ian-mi/sdk-go/v2/pkg/protocol"
)

// SendReceive does s.Send(in), then it receives the message in r.Receive() and executes outAssert
// Halt test on error.
func SendReceive(t *testing.T, ctx context.Context, in binding.Message, s protocol.Sender, r protocol.Receiver, outAssert func(binding.Message)) {
	t.Helper()
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		out, recvErr := r.Receive(ctx)
		require.NoError(t, recvErr)
		outAssert(out)
		finishErr := out.Finish(nil)
		require.NoError(t, finishErr)
	}()

	go func() {
		defer wg.Done()
		err := s.Send(ctx, in)
		require.NoError(t, err)
	}()

	wg.Wait()
}
