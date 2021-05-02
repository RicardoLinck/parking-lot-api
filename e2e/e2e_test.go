package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
	"gotest.tools/v3/assert/opt"

	"github.com/RicardoLinck/parking-lot-api/api"
	"github.com/RicardoLinck/parking-lot-api/barrier"
)

type AccessLog struct {
	Timestamp      time.Time `json:"timestamp"`
	BarrierID      string    `json:"barrier_id"`
	CarRegisration string    `json:"car_registration"`
	Direction      string    `json:"direction"`
}

type AccessLogResponse []AccessLog

var srv *httptest.Server

func TestMain(m *testing.M) {
	bc := barrier.NewBarrierConfig("./logs/barriers")
	defer os.RemoveAll("./logs/barriers")
	srv = httptest.NewServer(api.ConfigureEndpoints(bc))
	defer srv.Close()
	m.Run()
}

func Test_BarriersEndpoints(t *testing.T) {
	t.Run("/barrier/:barrierID/in", func(t *testing.T) {
		t.Run("returns 404 for barrier invalid barrier", func(t *testing.T) {
			checkBarrierNotFound(t, srv.URL+"/barrier/invalid-barrier/in/valid-car-reg", http.MethodPost)
		})

		t.Run("returns 200 for barrier valid barrier", func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/barrier/east/in/valid-car-reg", "", nil)
			assert.Check(t, err)
			assert.Check(t, cmp.Equal(resp.StatusCode, 200))

			var body map[string]string
			err = json.NewDecoder(resp.Body).Decode(&body)
			defer resp.Body.Close()
			assert.Check(t, err)

			assert.Check(t, cmp.DeepEqual(map[string]string{
				"message": "registration-id: valid-car-reg entered the parking lot using barrier east",
			}, body))
		})

		t.Run("/barrier/:barrierID/out", func(t *testing.T) {
			t.Run("returns 404 for barrier invalid barrier", func(t *testing.T) {
				checkBarrierNotFound(t, srv.URL+"/barrier/invalid-barrier/out/valid-car-reg", http.MethodPost)
			})

			t.Run("returns 200 for barrier valid barrier", func(t *testing.T) {
				resp, err := http.Post(srv.URL+"/barrier/east/out/valid-car-reg", "", nil)
				assert.Check(t, err)
				assert.Check(t, cmp.Equal(resp.StatusCode, 200))

				var body map[string]string
				err = json.NewDecoder(resp.Body).Decode(&body)
				defer resp.Body.Close()
				assert.Check(t, err)

				assert.Check(t, cmp.DeepEqual(map[string]string{
					"message": "registration-id: valid-car-reg exited the parking lot using barrier east",
				}, body))
			})
		})

		t.Run("/barrier/:barrierID/logs", func(t *testing.T) {
			t.Run("returns 404 for barrier invalid barrier", func(t *testing.T) {
				checkBarrierNotFound(t, srv.URL+"/barrier/invalid-barrier/logs", http.MethodGet)
			})

			t.Run("returns 200 with logs for in and out calls", func(t *testing.T) {
				resp, err := http.Get(srv.URL + "/barrier/east/logs")
				assert.Check(t, err)
				assert.Check(t, cmp.Equal(resp.StatusCode, 200))

				var body AccessLogResponse
				err = json.NewDecoder(resp.Body).Decode(&body)
				defer resp.Body.Close()
				assert.Check(t, err)

				assert.Check(t, cmp.DeepEqual(AccessLogResponse{
					{
						Timestamp:      time.Now(),
						BarrierID:      "east",
						CarRegisration: "valid-car-reg",
						Direction:      "in",
					},
					{
						Timestamp:      time.Now(),
						BarrierID:      "east",
						CarRegisration: "valid-car-reg",
						Direction:      "out",
					},
				}, body, opt.TimeWithThreshold(time.Second)))
			})
		})
	})
}

func checkBarrierNotFound(t *testing.T, url, method string) {
	t.Helper()

	req, err := http.NewRequest(method, url, nil)
	assert.Check(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.Check(t, err)
	assert.Check(t, cmp.Equal(resp.StatusCode, 404))

	var body map[string]string
	err = json.NewDecoder(resp.Body).Decode(&body)
	defer resp.Body.Close()
	assert.Check(t, err)

	assert.Check(t, cmp.DeepEqual(map[string]string{"error": "barrier not found"}, body))
}
