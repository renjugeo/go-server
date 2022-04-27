package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renjugeo/go-server/config"
	"github.com/renjugeo/go-server/util"
)

func TestFetchUrlSuccess(t *testing.T) {
	c := &config.Configuration{}
	config.SetDefaults(c)
	logger, _ := util.GetLogger(c)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := `{
			"data": [
			  {
				"url": "www.testurl1.com/",
				"views": 10,
				"relevanceScore": 0.1
			  },
			  {
				"url": "www.testurl2.com",
				"views": 7,
				"relevanceScore": 1
			  },
			  {
				"url": "www.testurl3.com",
				"views": 11,
				"relevanceScore": 1
			  },
			  {
				"url": "www.testurl4.com",
				"views": 20,
				"relevanceScore": 2
			  }
			]
		  }`
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	defer ts.Close()
	api := NewV1API(c, logger)
	stats, err := api.fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
	if len(stats.Data) != 4 {
		t.Fatalf("expected 4 items, got %d", len(stats.Data))
	}
}
