package eventbus

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sync"
)

type EventBus struct {
	pendingTasks    chan func()
	taskChanClosed  bool
	pendingEvents   chan any
	eventChanClosed bool
	handlers        map[reflect.Type][]func(any)
	mut             sync.RWMutex
	logger          *log.Logger
}

func New(logger *log.Logger) *EventBus {
	return &EventBus{
		pendingTasks:    make(chan func(), 64),
		taskChanClosed:  false,
		pendingEvents:   make(chan any, 100),
		eventChanClosed: false,
		handlers:        make(map[reflect.Type][]func(any)),
		logger:          logger,
	}
}

func (eb *EventBus) Run(ctx context.Context, wg *sync.WaitGroup) {
	var taskHandlerWg sync.WaitGroup
	defer wg.Done()
	numTaskHandlers := 1
	for i := 0; i < numTaskHandlers; i++ {
		taskHandlerWg.Add(1)
		go eb.taskHandler(ctx, &taskHandlerWg)
	}
loop:
	for {
		select {
		case event := <-eb.pendingEvents:
			if event != nil {
				eb.handle(event)
			} else {
				// break by closing channel, should be test only behaviour
				fmt.Println("break after closed chan")
				break loop
			}
		case <-ctx.Done():
			break loop
		}
	}
	taskHandlerWg.Wait()
}

func (eb *EventBus) taskHandler(ctx context.Context, wg *sync.WaitGroup) {
loop:
	for {
		select {
		case task := <-eb.pendingTasks:
			// task can be nil if chan is closed
			if task != nil {
				select {
				case <-ctx.Done():
					break loop
				default:
					eb.logger.Println("[EventBus] handling task")
					task()
				}
			} else {
				// break by closing channel, should be test only behaviour
				fmt.Println("break after closed chan")
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
		eb.logger.Println(err)
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
		return
	}
	eb.pendingEvents <- event
}

func (eb *EventBus) PushTask(task func()) {
	eb.mut.RLock()
	defer eb.mut.RUnlock()
	if eb.taskChanClosed {
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
