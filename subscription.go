package main

type Subscription struct {
	cancelled   bool
	Channel     chan bool
	EventSource *EventSource
}

func NewSubscription(e *EventSource) *Subscription {
	return &Subscription{Channel: make(chan bool)}
}

func (s *Subscription) cancel() {
	s.cancelled = true
	for i, sub := range s.EventSource.Subscriptions {
		if sub == s {
			s.EventSource.Subscriptions = append(s.EventSource.Subscriptions[:i], s.EventSource.Subscriptions[i+1:]...)
			break
		}
	}
	s.Channel <- false
}

func (s *Subscription) wait() {
}

type EventSource struct {
	Subscriptions []*Subscription
}

func NewEventSource() *EventSource {
	return &EventSource{Subscriptions: make([]*Subscription, 0)}
}

func (e *EventSource) subscribe() *Subscription {
	s := &Subscription{Channel: make(chan bool)}
	e.Subscriptions = append(e.Subscriptions, s)
	return s
}
