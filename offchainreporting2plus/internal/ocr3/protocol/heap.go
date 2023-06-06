package protocol

import (
	"container/heap"
	"time"
)

// Type safe wrapper around MinHeapTimeToContractReportInternal
type MinHeapTimeToPendingTransmission[RI any] struct {
	internal MinHeapTimeToPendingTransmissionInternal[RI]
}

func (h *MinHeapTimeToPendingTransmission[RI]) Push(item MinHeapTimeToPendingTransmissionItem[RI]) {
	heap.Push(&h.internal, item)
}

func (h *MinHeapTimeToPendingTransmission[RI]) Pop() MinHeapTimeToPendingTransmissionItem[RI] {
	return heap.Pop(&h.internal).(MinHeapTimeToPendingTransmissionItem[RI])
}

func (h *MinHeapTimeToPendingTransmission[RI]) Peek() MinHeapTimeToPendingTransmissionItem[RI] {
	return h.internal[0]
}

func (h *MinHeapTimeToPendingTransmission[RI]) Len() int {
	return h.internal.Len()
}

type MinHeapTimeToPendingTransmissionItem[RI any] struct {
	Time           time.Time
	SeqNr          uint64
	Index          int
	AttestedReport AttestedReportMany[RI]
}

// Implements heap.Interface and uses interface{} all over the place.
type MinHeapTimeToPendingTransmissionInternal[RI any] []MinHeapTimeToPendingTransmissionItem[RI]

func (pq MinHeapTimeToPendingTransmissionInternal[RI]) Len() int { return len(pq) }

func (pq MinHeapTimeToPendingTransmissionInternal[RI]) Less(i, j int) bool {
	return pq[i].Time.Before(pq[j].Time)
}

func (pq MinHeapTimeToPendingTransmissionInternal[RI]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *MinHeapTimeToPendingTransmissionInternal[RI]) Push(x interface{}) {
	item := x.(MinHeapTimeToPendingTransmissionItem[RI])
	*pq = append(*pq, item)
}

func (pq *MinHeapTimeToPendingTransmissionInternal[RI]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
