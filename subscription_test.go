package main

import (
	"fmt"
	"testing"
)

type sequence []int

func (s *sequence) add(i int) {
	fmt.Printf("add %d\n", i)
	*s = append(*s, i)
}

func (s *sequence) verify(t *testing.T, last int) {
	for i, x := range *s {
		if i != x {
			t.Errorf("%v %d != %d", *s, i, x)
			break
		}
		if i > last {
			t.Errorf("%v %d > %d", *s, i, last)
			break
		}
	}
}

func TestSubscriptionBasicOneLoop(t *testing.T) {

	seq := make(sequence, 0)

	event := NewEventSource()

	seq.add(0)
	subscription := event.subscribe()

	go func() {
		subscription.wait()
		seq.add(2)

		subscription.cancel()
	}()
	seq.add(1)
	event.emit()
	seq.add(3)

	seq.verify(t, 3)
}

func TestSubscriptionMultipleLoops(t *testing.T) {

	seq := make(sequence, 0)

	event := NewEventSource()

	seq.add(0)
	subscription := event.subscribe()

	go func() {
		subscription.wait()
		seq.add(2)

		subscription.wait()
		seq.add(4)

		subscription.wait()
		seq.add(6)

		subscription.cancel()
	}()

	seq.add(1)
	event.emit()

	seq.add(3)
	event.emit()

	seq.add(5)
	event.emit()

	seq.add(7)

	seq.verify(t, 7)
}

func TestSubscriptionExternalCancel(t *testing.T) {

	seq := make(sequence, 0)

	event := NewEventSource()

	seq.add(0)
	subscription := event.subscribe()

	go func() {
		for {
			subscription.wait()

			if subscription.Cancelled {
				break
			} else {
				seq.add(2)
			}
		}
	}()

	seq.add(1)
	event.emit()
	seq.add(3)

	subscription.cancel()
	seq.add(4)

	seq.verify(t, 4)
}

func TestSubscriptionMultipleListeners(t *testing.T) {

	seq := make(sequence, 0)

	event := NewEventSource()

	seq.add(0)
	subscription1 := event.subscribe()
	subscription2 := event.subscribe()

	go func() {
		subscription1.wait()
		seq.add(2)

		subscription1.cancel()
	}()

	go func() {
		subscription2.wait()
		seq.add(3)

		subscription2.cancel()
	}()

	seq.add(1)
	event.emit()

	seq.add(4)

	seq.verify(t, 4)
}
