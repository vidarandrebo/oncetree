package eventbus

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"

	"github.com/vidarandrebo/oncetree/consts"
)

type EventBus struct {
	pendingTasks    chan func()
	taskChanClosed  bool
	pendingEvents   chan any
	eventChanClosed bool
	handlers        map[reflect.Type][]func(any)
	mut             sync.RWMutex
	logger          *slog.Logger
}

func New(logger *slog.Logger) *EventBus {
	return &EventBus{
		pendingTasks:    make(chan func(), consts.EventBusQueueLength),
		taskChanClosed:  false,
		pendingEvents:   make(chan any, consts.EventBusQueueLength),
		eventChanClosed: false,
		handlers:        make(map[reflect.Type][]func(any)),
		logger:          logger.With(slog.Group("node", slog.String("module", "eventbus"))),
	}
}

func (eb *EventBus) Run(ctx context.Context, wg *sync.WaitGroup) {
	var handlerWaitGroup sync.WaitGroup
	defer wg.Done()
	for i := 0; i < consts.NumTaskHandlers; i++ {
		handlerWaitGroup.Add(1)
		go eb.taskHandler(ctx, &handlerWaitGroup)
	}

	for i := 0; i < consts.NumEventHandlers; i++ {
		handlerWaitGroup.Add(1)
		go eb.eventHandler(ctx, &handlerWaitGroup)
	}
	handlerWaitGroup.Wait()
}

func (eb *EventBus) eventHandler(ctx context.Context, wg *sync.WaitGroup) {
loop:
	for {
		select {
		case event := <-eb.pendingEvents:
			if event != nil {
				eb.handle(event)
			} else {
				// break by closing channel, should be test only behaviour
				eb.logger.Warn("event channel closed, exiting event handler loop")
				break loop
			}
		case <-ctx.Done():
			break loop
		}
	}
	wg.Done()
}

func (eb *EventBus) taskHandler(ctx context.Context, wg *sync.WaitGroup) {
loop:
	for {
		select {
		case task := <-eb.pendingTasks:
			// task can be nil if chan is closed
			if task != nil {
				task()
			} else {
				// break by closing channel, should be test only behaviour
				eb.logger.Warn("task channel closed, exiting task handler loop")
				break loop
			}
		case <-ctx.Done():
			break loop
		}
	}
	wg.Done()
}

func (eb *EventBus) handle(event any) {
	eventType := reflect.TypeOf(event)
	handlers, err := eb.eventHandlers(eventType)
	if err != nil {
		eb.logger.Error(
			"handling of event failed",
			slog.Any("event-type", eventType),
			slog.Any("err", err),
		)
		return
	}
	for _, handler := range handlers {
		handler(event)
	}
}

func (eb *EventBus) PushEvent(event any) {
	eb.mut.RLock()
	defer eb.mut.RUnlock()
	if eb.eventChanClosed {
		eb.logger.Error("event discarded due to closed channel")
		return
	}
	eb.pendingEvents <- event
}

func (eb *EventBus) PushTask(task func()) {
	if len(eb.pendingTasks) > 200 {
		eb.logger.Warn("over 200 entries in task queue")
	}
	eb.mut.RLock()
	defer eb.mut.RUnlock()
	if eb.taskChanClosed {
		eb.logger.Error("task discarded due to closed channel")
		return
	}
	eb.pendingTasks <- task
}

func (eb *EventBus) eventHandlers(eventType reflect.Type) ([]func(any), error) {
	eb.mut.RLock()
	defer eb.mut.RUnlock()
	handlers, ok := eb.handlers[eventType]
	if !ok {
		return nil, fmt.Errorf("no handlers for event of type %v", eventType)
	}
	return handlers, nil
}

func (eb *EventBus) RegisterHandler(eventType reflect.Type, fn func(any)) {
	eb.mut.Lock()
	defer eb.mut.Unlock()
	if _, ok := eb.handlers[eventType]; !ok {
		eb.handlers[eventType] = make([]func(any), 0)
	}
	eb.handlers[eventType] = append(eb.handlers[eventType], fn)
}

func (eb *EventBus) Execute(event any) {
	eventType := reflect.TypeOf(event)
	handlers, err := eb.eventHandlers(eventType)
	if err != nil {
		eb.logger.Error(
			"no handler for execution of event",
			slog.Any("event-type", eventType),
			slog.Any("err", err),
		)
		return
	}
	for _, handler := range handlers {
		handler(event)
	}
}
