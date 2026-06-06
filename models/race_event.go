package models

import "time"

type RaceEvent struct {
	EventID     string    `dynamodbav:"eventId" json:"eventId"`
	RaceID      string    `dynamodbav:"raceId" json:"raceId"`
	DriverID    string    `dynamodbav:"driverId" json:"driverId"`
	Lap         int       `dynamodbav:"lap" json:"lap"`
	EventType   EventType `dynamodbav:"eventType" json:"eventType"`
	Description string    `dynamodbav:"description" json:"description"`
	CreatedAt   time.Time `dynamodbav:"createdAt" json:"createdAt"`
}
