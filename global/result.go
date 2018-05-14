package global

import (
	"fmt"
	"time"
)

// Result - the query result as struct
type Result struct {
	TimeStamp      time.Time `json:"@timestamp"`
	Country        string    `json:"Country"`
	Store          string    `json:"Store"`
	UUID           string    `json:"UUID"`
	RouteID        string    `json:"RouteId"`
	Service        string    `json:"Service"`
	TargetFileName string    `json:"TargetFileName"`
	Key1           string    `json:"Key1"`
	Key2           string    `json:"Key2"`
	Key3           string    `json:"Key3"`
	Key4           string    `json:"Key4"`
	Status         string    `json:"Status"`
	Message        string    `json:"Message"`
}

func (result *Result) String() string {
	return fmt.Sprintf("%s | %s | %s | %s | %s | %s | %s | %s | %s | %s | %s",
		result.TimeStamp,
		result.Country,
		result.Store,
		result.UUID,
		result.RouteID,
		result.Key1,
		result.Key2,
		result.Key3,
		result.Key4,
		result.Status,
		result.Message,
	)
}
