package models

type EventType string

const (
	Overtake      EventType = "OVERTAKE"
	Crash         EventType = "CRASH"
	Penalty       EventType = "PENALTY"
	Investigation EventType = "INVESTIGATION"
	FastestLap    EventType = "FASTEST_LAP"
	PitStop       EventType = "PIT_STOP"
	SafetyCar     EventType = "SAFETY_CAR"
	Retirement    EventType = "RETIREMENT"
)

func (e EventType) IsValid() bool {
	switch e {
	case
		Overtake,
		Crash,
		Penalty,
		Investigation,
		FastestLap,
		PitStop,
		SafetyCar,
		Retirement:
		return true
	}

	return false
}
