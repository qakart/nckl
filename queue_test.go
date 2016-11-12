package main

import (
	. "github.com/aandryashin/matchers"
	"testing"
	"time"
)

const defaultTimeout = 100 * time.Millisecond

func TestSize(t *testing.T) {
	queue := CreateQueue(1)
	AssertThat(t, queue.Size(), EqualTo{0})
	queue.Push()
	AssertThat(t, queue.Size(), EqualTo{1})
	queue.Pop()
	AssertThat(t, queue.Size(), EqualTo{0})
}

func TestSetCapacity(t *testing.T) {
	queue := CreateQueue(1)
	queue.Push()
	queue.SetCapacity(2)
	AssertThat(t, queue.Capacity(), EqualTo{2})
	queue.Push()
	queue.Push()
	AssertThat(t, queue.Size(), EqualTo{3})
	AssertThat(t, actionTimeouts(queue.Push), EqualTo{true})
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false})
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false})
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false})
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false}) //This one is the last push data
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{true})
}

func TestSetCapacityZeroLength(t *testing.T) {
	queue := CreateQueue(1)
	queue.Push()
	queue.Pop() //There's only one channel in slice but it's already empty and should be deleted
	queue.SetCapacity(2)
	queue.Push()
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false})
}

func TestPop(t *testing.T) {
	queue := CreateQueue(2)
	queue.Push()
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{false})
	AssertThat(t, actionTimeouts(queue.Pop), EqualTo{true})
}

func actionTimeouts(action func()) bool {
	timeout := make(chan bool, 1)
	ch := make(chan bool)
	go func() {
		action()
		ch <- true
	}()
	go func() {
		time.Sleep(defaultTimeout)
		timeout <- true
	}()
	select {
	case <-ch:
		return false
	case <-timeout:
		return true
	}
}
