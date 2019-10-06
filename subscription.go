package main

import "fmt"

type SubscriptionState int

type Subscription struct {
	Channel     chan bool
	EventSource *EventSource
	Cancelled   bool
}

func NewSubscription(e *EventSource) *Subscription {
	return &Subscription{
		EventSource: e,
		Channel:     make(chan bool),
	}
}

// should only be called by subscribing goroutine
func (s *Subscription) cancel() {
	fmt.Println("any: cancelling subscription")
	if s.Cancelled {
		return
	}
	for i, sub := range s.EventSource.Subscriptions {
		if sub == s {
			s.EventSource.Subscriptions = append(s.EventSource.Subscriptions[:i], s.EventSource.Subscriptions[i+1:]...)
			break
		}
	}
	s.Cancelled = true
	// someones gotta be waiting, so let 'em know
	s.Channel <- false
}

func (s *Subscription) firstWait() {
	if s.Cancelled {
		return
	}
	// Wait for the event
	fmt.Println("sub: first: listening for source")
	<-s.Channel
	fmt.Printf("sub: first: source responded\n")
}

func (s *Subscription) wait() {
	if s.Cancelled {
		return
	}
	// Wait for the event
	fmt.Println("sub: listening for source")
	s.Channel <- true
	<-s.Channel
	fmt.Printf("sub: source responded\n")
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
	s.firstWait()
	return s
}

func (e *EventSource) emit() {
	for _, s := range e.Subscriptions {
		if !s.Cancelled {
			fmt.Println("source: listening for response from subscriber")
			s.Channel <- true
			<-s.Channel
			fmt.Println("source: someone responded")
		}
	}
}
