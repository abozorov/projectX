package requestQueue

import (
	"strconv"
	"sync"
	"time"
)

type QueueLimit struct {
	mu           *sync.Mutex
	queue        map[int][]ClientRequest
	doneRequests chan string
}

func NewQueueLimit(buffer int) *QueueLimit {
	return &QueueLimit{
		mu:           &sync.Mutex{},
		queue:        make(map[int][]ClientRequest),
		doneRequests: make(chan string, buffer),
	}
}

func (q *QueueLimit) GetEndTime(cliID, reqID int) time.Time {
	for _, v := range q.queue[cliID] {
		if v.RequestID == reqID {
			return v.EndTime
		}
	}
	return time.Now().Add(time.Minute)
}

func (q *QueueLimit) GetClientRequests(cliID int) []ClientRequest {
	return q.queue[cliID]
}

func (q *QueueLimit) QueueLen(cliID int) int {
	return len(q.queue[cliID])
}

func (q *QueueLimit) DeleteRequest(cliID, reqID int) {
	for i, v := range q.queue[cliID] {
		if v.RequestID == reqID {
			q.queue[cliID] = append(q.queue[cliID][:i], q.queue[cliID][i+1:]...)
			break
		}
	}
}

func (q *QueueLimit) AddClientRequest(cliID, reqID int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	newCliRequest := *NewClientRequest(cliID, reqID)
	// if q.queue[cliID] == nil {
	// 	q.queue[cliID] = make([]ClientRequest, 0, 10)
	// }
	q.queue[cliID] = append(q.queue[cliID], newCliRequest)
}

func (q *QueueLimit) CheckQueue(cliID int) bool {
	return len(q.queue[cliID]) > 5
}

func (q *QueueLimit) Publish(cliID, reqID int) {
	q.doneRequests <- strconv.Itoa(cliID) + ":" + strconv.Itoa(reqID)
}

func (q *QueueLimit) Subscribe() <-chan string {
	return q.doneRequests
}
