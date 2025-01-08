package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ognick/word_of_wisdom/pkg/shutdown"
)

var UnexpectedCloseComponent = errors.New("unexpected close component")

type Component interface {
	Run(ctx context.Context, ready chan struct{}) error
}

type Lifecycle interface {
	Register(component Component)
	RunAllComponents(
		runner shutdown.Runner,
		gracefulCtx context.Context,
	)
}

type lifecycle struct {
	components []Component
}

func NewLifecycle() Lifecycle {
	return &lifecycle{}
}

func (lc *lifecycle) Register(component Component) {
	lc.components = append(lc.components, component)
}

func (lc *lifecycle) RunAllComponents(
	runner shutdown.Runner,
	gracefulCtx context.Context,
) {
	componentCount := len(lc.components)
	cancelCtxFuncs := make([]context.CancelFunc, componentCount)
	cancelPreviousCtx := func(i int) {
		if i > 0 {
			cancelCtxFuncs[i-1]()
		}
	}

	lcCtx, lcCtxCancel := context.WithCancelCause(gracefulCtx)
	for i, component := range lc.components {
		ctx, cancelCtx := context.WithCancel(context.Background())
		cancelCtxFuncs[i] = cancelCtx
		ready := make(chan struct{})
		runner.Go(func() error {
			defer cancelPreviousCtx(i)
			err := component.Run(ctx, ready)
			if err == nil && lcCtx.Err() == nil {
				componentName := reflect.TypeOf(component).String()
				err = fmt.Errorf("%s: %w", componentName, UnexpectedCloseComponent)
				lcCtxCancel(err)
			}
			return err
		})
		select {
		case <-ready:
		case <-lcCtx.Done():
		}
	}

	runner.Go(func() error {
		<-lcCtx.Done()
		cancelPreviousCtx(componentCount)
		lcCtxCancel(lcCtx.Err())
		return nil
	})
}
