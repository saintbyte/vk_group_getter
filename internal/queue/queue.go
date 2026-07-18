package queue

type Handler func(id uint)

type Queue struct {
	ch      chan uint
	handler Handler
}

func New(bufferSize int, handler Handler) *Queue {
	return &Queue{
		ch:      make(chan uint, bufferSize),
		handler: handler,
	}
}

func (q *Queue) Push(groupID uint) {
	q.ch <- groupID
}

func (q *Queue) Start() {
	go func() {
		for id := range q.ch {
			q.handler(id)
		}
	}()
}
