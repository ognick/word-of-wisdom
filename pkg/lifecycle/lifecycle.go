package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ognick/word_of_wisdom/pkg/logger"
	"github.com/ognick/word_of_wisdom/pkg/shutdown"
)

var UnexpectedCloseComponent = errors.New("unexpected close component")

type Component interface {
	Run(ctx context.Context, readinessProbe chan error) error
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
	log        logger.Logger
}

func NewLifecycle(log logger.Logger) Lifecycle {
	return &lifecycle{log: log}
}

func (lc *lifecycle) Register(component Component) {
	lc.components = append(lc.components, component)
}

func (lc *lifecycle) RunAllComponents(
	runner shutdown.Runner,
	gracefulCtx context.Context,
) {
	cancelCtxFuncs := make([]context.CancelFunc, 0, len(lc.components))
	cancelPreviousCtx := func(i int) {
		if i > 0 {
			cancelCtxFuncs[i-1]()
		}
	}

	lcCtx, lcCtxCancel := context.WithCancelCause(gracefulCtx)
	defer runner.Go(func() error {
		startedComponentCount := len(cancelCtxFuncs)
		defer cancelPreviousCtx(startedComponentCount)
		lc.log.Infof("lifecycle started %d components", startedComponentCount)
		<-lcCtx.Done()
		cause := context.Cause(lcCtx)
		if errors.Is(cause, context.Canceled) {
			lc.log.Infof("lifecycle gracefully stopped")
			return nil
		}
		lc.log.Infof("lifecycle shutdown with error: %v", cause)
		return cause
	})

	for i, component := range lc.components {
		ctx, cancelCtx := context.WithCancel(context.Background())
		cancelCtxFuncs = append(cancelCtxFuncs, cancelCtx)
		readinessProbe := make(chan error)
		componentName := reflect.TypeOf(component).String()
		runner.Go(func() error {
			defer cancelPreviousCtx(i)
			err := component.Run(ctx, readinessProbe)
			if err == nil && lcCtx.Err() == nil {
				componentName := reflect.TypeOf(component).String()
				err = fmt.Errorf("%s: %w", componentName, UnexpectedCloseComponent)
				lc.log.Errorf("[%d] Unexpected close component %s", i, componentName)
				lcCtxCancel(err)
			}
			if err != nil {
				lc.log.Errorf("[%d] Component %s closed with error %v", i, componentName, err)
				return err
			}

			lc.log.Infof("[%d] Component %s closed", i, componentName)
			return nil
		})
		select {
		case err := <-readinessProbe:
			if err != nil {
				lc.log.Errorf("[%d] readiness probe component %s error: %v", i, componentName, err)
				lcCtxCancel(err)
				return
			}
			lc.log.Infof("[%d] readiness probe component %s", i, componentName)
		case <-lcCtx.Done():
			return
		}
	}
}
