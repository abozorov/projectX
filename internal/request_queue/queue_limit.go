package requestQueue

import (
	"sync"
	"time"
)

type request struct {
	clientId  int
	requestId int
}

type QueueLimit struct {
	mu           sync.Mutex
	queue        map[int][]ClientRequest
	doneRequests chan request
}

func NewQueueLimit(buffer int) *QueueLimit {
	return &QueueLimit{
		queue:        make(map[int][]ClientRequest),
		doneRequests: make(chan request, buffer),
	}
}

func (q *QueueLimit) GetEndTime(cliID, reqID int) time.Time {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, v := range q.queue[cliID] {
		if v.RequestID == reqID {
			return v.EndTime
		}
	}
	return time.Now().Add(time.Minute)
}

func (q *QueueLimit) GetClientRequests(cliID int) []ClientRequest {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.queue[cliID]
}

func (q *QueueLimit) QueueLen(cliID int) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue[cliID])
}

func (q *QueueLimit) DeleteRequest(cliID, reqID int) {
	q.mu.Lock()
	defer q.mu.Unlock()

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
	q.queue[cliID] = append(q.queue[cliID], newCliRequest)
}

func (q *QueueLimit) CheckQueue(cliID int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.queue[cliID]) > 5
}

func (q *QueueLimit) Publish(cliID, reqID int) {
	q.doneRequests <- request{
		clientId:  cliID,
		requestId: reqID,
	}
}

func (q *QueueLimit) Subscribe() <-chan request {
	return q.doneRequests
}
