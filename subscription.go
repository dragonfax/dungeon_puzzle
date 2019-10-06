package main

import "fmt"

type Subscription struct {
	SubscriberChannel   chan bool
	EventChannel        chan bool
	EventSource         *EventSource
	SubscriberListening bool
	EventListening      bool
	Cancelled           bool
}

func NewSubscription(e *EventSource) *Subscription {
	return &Subscription{
		EventSource:       e,
		SubscriberChannel: make(chan bool),
		EventChannel:      make(chan bool),
	}
}

func (s *Subscription) cancel() {
	fmt.Println("any: cancelling subscription")
	if s.Cancelled {
		return
	}
	s.Cancelled = true
	for i, sub := range s.EventSource.Subscriptions {
		if sub == s {
			s.EventSource.Subscriptions = append(s.EventSource.Subscriptions[:i], s.EventSource.Subscriptions[i+1:]...)
			break
		}
	}
	if s.SubscriberListening {
		fmt.Println("any: notifying subscriber of cancellation")
		s.SubscriberChannel <- true
	}
	if s.EventListening {
		fmt.Println("any: notifying source of cancellation")
		s.EventChannel <- true
	}
}

func (s *Subscription) wait() {
	if s.Cancelled {
		return
	}
	if s.EventListening {
		fmt.Println("sub: notify source")
		s.EventChannel <- true
		fmt.Println("sub: source received notification")
	}

	// Wait for the event
	s.SubscriberListening = true
	fmt.Println("sub: listening for source")
	<-s.SubscriberChannel
	fmt.Println("sub: source responded")
	s.SubscriberListening = false
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
		if !s.Cancelled {
			fmt.Println("source: notifying subscriber")
			s.SubscriberChannel <- true
			fmt.Println("source: subscriber received notification")
			s.EventListening = true
			fmt.Println("source: listening for response from subscriber")
			<-s.EventChannel
			fmt.Println("source: subscriber responded")
			s.EventListening = false
		}
	}
}
