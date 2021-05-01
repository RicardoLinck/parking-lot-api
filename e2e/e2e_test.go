package e2e

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/RicardoLinck/parking-lot-api/api"
	"github.com/RicardoLinck/parking-lot-api/barrier"
)

var srv *httptest.Server

func TestMain(m *testing.M) {
	bc := barrier.NewBarrierConfig("./logs/barriers")
	srv = httptest.NewServer(api.ConfigureEndpoints(bc))
	defer srv.Close()
	os.Exit(m.Run())
}

func Test_BarriersEndpoints(t *testing.T) {
	t.Run("/barriers/in", func(t *testing.T) {
		t.Run("returns 404 for barrier invalid barrier", func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/barrier/invalid-barrier/in/valid-car-reg", "", nil)
			assert.Check(t, err)
			assert.Check(t, cmp.Equal(resp.StatusCode, 404))

			var mbody map[string]string
			err = json.NewDecoder(resp.Body).Decode(&mbody)
			defer resp.Body.Close()
			assert.Check(t, err)

			assert.Check(t, cmp.DeepEqual(map[string]string{"error": "barrier not found"}, mbody))
		})
	})

}
