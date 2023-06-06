package protocol

// MessageBuffer implements a fixed capacity ringbuffer for items of type
// MessageToReportGeneration
type MessageBuffer[RI any] struct {
	start  int
	length int
	buffer []*MessageToReportGeneration[RI]
}

func NewMessageBuffer[RI any](cap int) *MessageBuffer[RI] {
	return &MessageBuffer[RI]{
		0,
		0,
		make([]*MessageToReportGeneration[RI], cap),
	}
}

// Peek at the front item
func (rb *MessageBuffer[RI]) Peek() *MessageToReportGeneration[RI] {
	if rb.length == 0 {
		return nil
	} else {
		return rb.buffer[rb.start]
	}
}

// Pop front item
func (rb *MessageBuffer[RI]) Pop() *MessageToReportGeneration[RI] {
	result := rb.Peek()
	if result != nil {
		rb.buffer[rb.start] = nil
		rb.start = (rb.start + 1) % len(rb.buffer)
		rb.length--
	}
	return result
}

// Push new item to back. If the additional item would lead
// to the capacity being exceeded, remove the front item first
func (rb *MessageBuffer[RI]) Push(msg MessageToReportGeneration[RI]) {
	if rb.length == len(rb.buffer) {
		rb.Pop()
	}
	rb.buffer[(rb.start+rb.length)%len(rb.buffer)] = &msg
	rb.length++
}
