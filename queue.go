package main

// An extensible fixed size blocking queue based on channels.
// Internally we store a list of channels with fixed size. When pushing an item
// we always add to the last channel (i.e. the newest one). When popping an item
// we use the first channel. We remove channels from the list when they are
// emptied.
type Queue interface {
	Push()
	Pop()
	Size() int
	Capacity() int
	SetCapacity(newCapacity int)
}

func CreateQueue(initialCapacity int) *queueImpl {
	var channels []chan struct{}
	ret := &queueImpl{
		channels: channels,
	}
	ret.SetCapacity(initialCapacity)
	return ret
}

type queueImpl struct {
	channels []chan struct{}
}

func (q *queueImpl) Push() {
	q.channels[len(q.channels)-1] <- struct{}{}
}

func (q *queueImpl) Pop() {
	<-q.channels[0]
	q.cleanupChannels()
}

func (q *queueImpl) cleanupChannels() {
	if len(q.channels) > 1 && len(q.channels[0]) == 0 {
		close(q.channels[0])
		q.channels = q.channels[1:]
	}
}

func (q *queueImpl) Size() int {
	size := 0
	for _, ch := range q.channels {
		size += len(ch)
	}
	return size
}

func (q *queueImpl) Capacity() int {
	return cap(q.channels[len(q.channels)-1])
}

func (q *queueImpl) SetCapacity(newCapacity int) {
	if len(q.channels) == 0 || q.Capacity() != newCapacity {
		q.channels = append(q.channels, make(chan struct{}, newCapacity))
	}
	q.cleanupChannels()
}
