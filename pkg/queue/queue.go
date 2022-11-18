package queue

type queue struct {
	elements chan Element
	name     string
}

func NewQueue(name string) IQueue {
	return &queue{
		elements: make(chan Element),
		name:     name,
	}
}

func (q *queue) EnQueue(e Element) error {
	go func(e Element) {
		q.elements <- e
	}(e)
	return nil
}

func (q *queue) DeQueue() (Element, error) {
	return <-q.elements, nil
}

func (q *queue) GetName() string {
	return q.name
}

func (q *queue) Close() error {
	close(q.elements)
	return nil
}
