package lifecycle

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

type TestComponent struct {
	Component
	id          int
	startOutput chan int
	stopOutput  chan int
	close       chan error
}

func NewBaseComponent(
	lc Lifecycle,
	id int,
	startOutput chan int,
	stopOutput chan int,
) *TestComponent {
	return RegisterComponent(lc,
		&TestComponent{
			id:          id,
			startOutput: startOutput,
			stopOutput:  stopOutput,
			close:       make(chan error),
		})
}

func (c *TestComponent) Close(err error) {
	c.close <- err
}

func (c *TestComponent) Run(ctx context.Context, ready chan struct{}) error {
	c.startOutput <- c.id

	close(ready)
	var err error
	select {
	case err = <-c.close:
	case <-ctx.Done():
	}

	c.stopOutput <- c.id
	return err
}

func StartAllComponents(
	ctx context.Context,
	t *testing.T,
	componentCount int,
) (
	lc Lifecycle,
	stopOutput chan int,
	testComponents []*TestComponent,
	errGroupWait func() error,
) {
	lc = NewLifecycle()
	stopOutput = make(chan int, componentCount)
	testComponents = make([]*TestComponent, componentCount)
	startOutput := make(chan int, componentCount)
	components := make([]Component, componentCount)
	for i := range componentCount {
		testComponent := NewBaseComponent(lc, i, startOutput, stopOutput)
		testComponents[i] = testComponent
		components[i] = testComponent
	}

	errGroup, errCtx := errgroup.WithContext(ctx)
	lc.RunAllComponents(errGroup, errCtx)
	for i := range componentCount {
		id := <-startOutput
		require.Equal(t, i, id)
	}

	return lc, stopOutput, testComponents, errGroup.Wait
}

func Test_StartAllComponents_Graceful_Stop(t *testing.T) {
	const componentCount int = 10
	ctx, cancel := context.WithCancel(context.Background())
	_, stopOutput, _, errGroupWait := StartAllComponents(ctx, t, componentCount)
	cancel()
	err := errGroupWait()
	require.NoError(t, err)
	for i := range componentCount {
		j := componentCount - i - 1
		id := <-stopOutput
		require.Equal(t, j, id)
	}
}

func Test_StartAllComponents_Stop_With_Error(t *testing.T) {
	const componentCount int = 10
	var internalApplicationError = errors.New("internal application error")
	ctx, cancel := context.WithCancelCause(context.Background())
	_, stopOutput, _, errGroupWait := StartAllComponents(ctx, t, componentCount)
	cancel(internalApplicationError)
	err := errGroupWait()
	require.NoError(t, err)
	for i := range componentCount {
		j := componentCount - i - 1
		id := <-stopOutput
		require.Equal(t, j, id)
	}
}

func Test_StartAllComponents_UnexpectedCloseComponent(t *testing.T) {
	const (
		componentCount     int = 10
		stoppedComponentID int = 5
	)
	_, stopOutput, testComponents, errGroupWait := StartAllComponents(context.Background(), t, componentCount)

	testComponents[stoppedComponentID].Close(nil)
	require.Equal(t, stoppedComponentID, <-stopOutput)
	err := errGroupWait()
	require.ErrorAs(t, err, &UnexpectedCloseComponent)
	stoppedComponentCount := 1
	for range componentCount - 1 {
		<-stopOutput
		stoppedComponentCount++
	}
	require.Equal(t, componentCount, stoppedComponentCount)
}

func Test_StartAllComponents_CloseComponent_With_Error(t *testing.T) {
	const (
		componentCount     int = 10
		stoppedComponentID int = 5
	)

	var internalComponentError = errors.New("internal component error")

	_, stopOutput, testComponents, errGroupWait := StartAllComponents(context.Background(), t, componentCount)

	testComponents[stoppedComponentID].Close(internalComponentError)
	require.Equal(t, stoppedComponentID, <-stopOutput)
	err := errGroupWait()
	require.ErrorIs(t, err, internalComponentError)
	stoppedComponentCount := 1
	for range componentCount - 1 {
		<-stopOutput
		stoppedComponentCount++
	}
	require.Equal(t, componentCount, stoppedComponentCount)
}
