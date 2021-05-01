package barrier

import (
	"os"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func Test_SaveAccessLog(t *testing.T) {
	const logsPath = "./test-logs"

	defer os.RemoveAll(logsPath)

	t.Run("creates one file for each barrier id", func(t *testing.T) {
		accessLogs := []AccessLog{
			{
				Timestamp:      time.Now().UTC(),
				BarrierID:      "test-barrier-id",
				CarRegisration: "test-car-reg",
				Direction:      "in",
			},
			{
				Timestamp:      time.Now().UTC(),
				BarrierID:      "another-barrier-id",
				CarRegisration: "test-car-reg",
				Direction:      "out",
			},
		}

		t.Run("creates path when it does not exist", func(t *testing.T) {
			if _, err := os.Stat(logsPath); err == nil {
				_ = os.RemoveAll(logsPath)
			}

			err := SaveAccessLog(accessLogs[0], logsPath)
			assert.Check(t, err)

			f, err := os.Stat(logsPath)
			assert.Check(t, err)
			assert.Check(t, f.IsDir())
		})

		err := SaveAccessLog(accessLogs[1], logsPath)
		assert.Check(t, err)

		files, err := os.ReadDir(logsPath)
		assert.Check(t, err)
		assert.Check(t, cmp.Len(files, 2))

		for _, v := range accessLogs {
			fromFile, err := ReadAccessLog(v.BarrierID, logsPath)
			assert.Check(t, err)
			assert.Check(t, cmp.DeepEqual(fromFile[0], v))
		}
	})
}
