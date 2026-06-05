package requestQueue

import "time"

type ClientRequest struct {
	RequestID int
	ClientID  int
	EndTime   time.Time
}

func NewClientRequest(clientID, requesID int) *ClientRequest {
	return &ClientRequest{
		ClientID:  clientID,
		RequestID: requesID,
		EndTime:   time.Now().Add(time.Minute),
		// EndTime: time.Now().Add(time.Second * 10),
	}
}
