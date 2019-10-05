package main

type SubState int

const (
	NEW SubState = iota
	WAITING
	RUNNING
	CANCELLED
)

type Subscription struct {
	State       SubState
	Channel     chan bool
	EventSource *EventSource
}

func NewSubscription(e *EventSource) *Subscription {
	return &Subscription{EventSource: e, Channel: make(chan bool), State: NEW}
}

func (s *Subscription) cancel() {
	if s.State == RUNNING {
		// let the source know we're done.
		s.State = CANCELLED
		s.Channel <- false
	} else {
		s.State = CANCELLED
	}
	for i, sub := range s.EventSource.Subscriptions {
		if sub == s {
			s.EventSource.Subscriptions = append(s.EventSource.Subscriptions[:i], s.EventSource.Subscriptions[i+1:]...)
			break
		}
	}
}

func (s *Subscription) wait() {
	if s.State == CANCELLED {
		return
	}
	if s.State == RUNNING {
		// let the source know we're done
		s.Channel <- true
	}

	// Wait for the event
	s.State = WAITING
	c := <-s.Channel
	if !c {
		s.State = CANCELLED
	}
	if s.State != CANCELLED {
		s.State = RUNNING
	}
}

func (s *Subscription) isCancelled() bool {
	return s.State == CANCELLED
}

type EventSource struct {
	Subscriptions []*Subscription
}

func NewEventSource() *EventSource {
	return &EventSource{Subscriptions: make([]*Subscription, 0)}
}

func (e *EventSource) subscribe() *Subscription {
	s := NewSubscription(e)
	e.Subscriptions = append(e.Subscriptions, s)
	return s
}

func (e *EventSource) emit() {
	for _, s := range e.Subscriptions {
		if s.State != CANCELLED {
			s.Channel <- true
			<-s.Channel
		}
	}
}
