package barrier

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

type AccessLog struct {
	Timestamp      time.Time `json:"timestamp"`
	BarrierID      string    `json:"barrier_id"`
	CarRegisration string    `json:"car_registration"`
	Direction      string    `json:"direction"`
}

func SaveAccessLog(a AccessLog, accessLogsPath string) error {
	if _, err := os.Stat(accessLogsPath); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(accessLogsPath, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/access_logs_%s.csv", accessLogsPath, a.BarrierID),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()
	return w.Write([]string{a.Timestamp.UTC().Format(time.RFC3339Nano), a.BarrierID, a.CarRegisration, a.Direction})
}

func ReadAccessLog(barrierId, accessLogsPath string) ([]AccessLog, error) {
	f, err := os.Open(fmt.Sprintf("%s/access_logs_%s.csv", accessLogsPath, barrierId))
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	logs := make([]AccessLog, len(records))
	for i, rec := range records {
		t, err := time.Parse(time.RFC3339Nano, rec[0])
		if err != nil {
			return nil, err
		}
		logs[i] = AccessLog{
			Timestamp:      t,
			BarrierID:      rec[1],
			CarRegisration: rec[2],
			Direction:      rec[3],
		}
	}
	return logs, nil
}
