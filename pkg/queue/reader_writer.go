package queue

type readerWriter struct {
	queue IQueue
}

func (r *readerWriter) Read() (Element, error) {
	return r.queue.DeQueue()
}

func (r *readerWriter) Close() error {
	return r.queue.Close()
}

func (w *readerWriter) Write(e Element) error {
	return w.queue.EnQueue(e)
}

func newReaderWriter(queue IQueue) *readerWriter {
	return &readerWriter{
		queue: queue,
	}
}

func NewReader(queue IQueue) IReader {
	return newReaderWriter(queue)
}

func NewWriter(queue IQueue) IWriter {
	return newReaderWriter(queue)
}
