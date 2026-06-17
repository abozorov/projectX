package models

type Stats struct {
	TotalUsers    int     `json:"total_users"`
	ActiveUsers   int     `json:"active_users"`
	InactiveUsers int     `json:"inactive_users"`
	AverageAge    float64 `json:"average_age"`
}

func NewStats() *Stats {
	return &Stats{}
}
