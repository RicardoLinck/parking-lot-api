package barrier

import (
	"errors"
	"time"
)

var ErrBarrierNotFound error = errors.New("barrier not found")

type BarrierConfig struct {
	barriers       map[string]bool
	accessLogsPath string
}

func (b *BarrierConfig) In(id string, carRegistration string) error {
	return SaveAccessLog(AccessLog{
		Timestamp:      time.Now(),
		BarrierID:      id,
		CarRegisration: carRegistration,
		Direction:      "in",
	}, b.accessLogsPath)
}

func (b *BarrierConfig) Out(id string, carRegistration string) error {
	return SaveAccessLog(AccessLog{
		Timestamp:      time.Now(),
		BarrierID:      id,
		CarRegisration: carRegistration,
		Direction:      "out",
	}, b.accessLogsPath)
}

func (b *BarrierConfig) Validate(id string) error {
	if _, ok := b.barriers[id]; !ok {
		return ErrBarrierNotFound
	}
	return nil
}

func NewBarrierConfig(accessLogsPath string) *BarrierConfig {
	return &BarrierConfig{
		barriers: map[string]bool{
			"east":  false,
			"west":  false,
			"north": false,
			"south": false,
		},
		accessLogsPath: accessLogsPath,
	}
}
