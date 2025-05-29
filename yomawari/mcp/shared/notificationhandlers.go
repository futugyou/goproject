package shared

import (
	"context"
	"fmt"
	"sync"

	"github.com/futugyou/yomawari/mcp/protocol"
)

type NotificationHandlers struct {
	mu       sync.RWMutex
	handlers map[string]*registration
}

type registration struct {
	handler   protocol.NotificationHandler
	temporary bool
	prev      *registration
	next      *registration
	parent    *NotificationHandlers
	method    string

	activeMu sync.Mutex
	activeWg sync.WaitGroup
	disposed bool
}

func NewNotificationHandlers() *NotificationHandlers {
	return &NotificationHandlers{
		handlers: make(map[string]*registration),
	}
}

func (nh *NotificationHandlers) RegisterRange(handlers map[string]protocol.NotificationHandler) {
	nh.mu.Lock()
	defer nh.mu.Unlock()

	for method, handler := range handlers {
		nh.registerLocked(method, handler, false)
	}
}

func (nh *NotificationHandlers) Register(method string, handler protocol.NotificationHandler, temporary bool) *RegistrationHandle {
	nh.mu.Lock()
	defer nh.mu.Unlock()

	reg := nh.registerLocked(method, handler, temporary)
	return &RegistrationHandle{reg: reg}
}

func (nh *NotificationHandlers) registerLocked(method string, handler protocol.NotificationHandler, temporary bool) *registration {
	reg := &registration{
		handler:   handler,
		temporary: temporary,
		parent:    nh,
		method:    method,
	}

	if head, exists := nh.handlers[method]; exists {
		reg.next = head
		head.prev = reg
	}
	nh.handlers[method] = reg

	return reg
}

func (nh *NotificationHandlers) InvokeHandlers(ctx context.Context, method string, notification *protocol.JsonRpcNotification) error {
	nh.mu.RLock()
	head, exists := nh.handlers[method]
	if !exists {
		nh.mu.RUnlock()
		return nil
	}

	var handlers []*registration
	for reg := head; reg != nil; reg = reg.next {
		if !reg.isDisposed() {
			handlers = append(handlers, reg)
		}
	}
	nh.mu.RUnlock()

	var errs []error
	for _, reg := range handlers {
		if err := reg.invoke(ctx, notification); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("handler errors: %v", errs)
	}
	return nil
}

type RegistrationHandle struct {
	reg *registration
}

func (rh *RegistrationHandle) Unregister() {
	if rh.reg == nil {
		return
	}
	rh.reg.dispose()
	rh.reg = nil
}

func (r *registration) isDisposed() bool {
	r.activeMu.Lock()
	defer r.activeMu.Unlock()
	return r.disposed
}

func (r *registration) invoke(ctx context.Context, notification *protocol.JsonRpcNotification) error {
	if !r.temporary {
		return r.handler(ctx, notification)
	}

	r.activeMu.Lock()
	if r.disposed {
		r.activeMu.Unlock()
		return nil
	}
	r.activeWg.Add(1)
	r.activeMu.Unlock()

	defer r.activeWg.Done()
	return r.handler(ctx, notification)
}

func (r *registration) dispose() {
	if !r.temporary {
		return
	}

	r.activeMu.Lock()
	if r.disposed {
		r.activeMu.Unlock()
		return
	}
	r.disposed = true
	r.activeMu.Unlock()

	r.activeWg.Wait()

	r.parent.mu.Lock()
	defer r.parent.mu.Unlock()

	if r.prev != nil {
		r.prev.next = r.next
	} else if head, exists := r.parent.handlers[r.method]; exists && head == r {
		r.parent.handlers[r.method] = r.next
	}

	if r.next != nil {
		r.next.prev = r.prev
	}

	if head, exists := r.parent.handlers[r.method]; exists && head == nil {
		delete(r.parent.handlers, r.method)
	}
}
