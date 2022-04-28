package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/renjugeo/go-server/config"
	"github.com/renjugeo/go-server/util"
	"github.com/stretchr/testify/assert"
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
	api := NewV1API(c, nil, logger)
	stats, err := api.fetchUrl(ts.URL)
	if err != nil {
		t.Fatalf("expected err to be nil, got %v", err)
	}
	if len(stats.Data) != 4 {
		t.Fatalf("expected 4 items, got %d", len(stats.Data))
	}
}

func TestGetStatsSuccess(t *testing.T) {
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
				"views": 2,
				"relevanceScore": 2
			  }
			]
		  }`
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	c.StatsEndpoints = []string{ts.URL}
	defer ts.Close()
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=relevanceScore&limit=2", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(apiResp.Data), 2)
	assert.Equal(t, apiResp.Count, 2)
	assert.Equal(t, apiResp.Data[0].Score, float32(2))
	assert.Equal(t, apiResp.Data[0].Url, "www.testurl4.com")

}

func TestGetStatsByViewsSuccess(t *testing.T) {
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
				"views": 2,
				"relevanceScore": 2
			  }
			]
		  }`
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	c.StatsEndpoints = []string{ts.URL}
	defer ts.Close()
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=views&limit=3", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(apiResp.Data), 3)
	assert.Equal(t, apiResp.Count, 3)
	assert.Equal(t, apiResp.Data[0].Views, int32(11))
	assert.Equal(t, apiResp.Data[0].Url, "www.testurl3.com")

}

func TestGetStatsInvalidLimit(t *testing.T) {
	c := &config.Configuration{}
	config.SetDefaults(c)
	logger, _ := util.GetLogger(c)
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=views&limit=abc", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIErrorResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, apiResp.ErrorCode, 500)
	assert.Equal(t, apiResp.Message, ErrorInvalidLimit)

}

func TestGetStatsInvalidSortKey(t *testing.T) {
	c := &config.Configuration{}
	config.SetDefaults(c)
	logger, _ := util.GetLogger(c)
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=something&limit=2", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIErrorResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, apiResp.ErrorCode, 500)
	assert.Equal(t, apiResp.Message, ErrorInvalidSortKey)

}

func TestGetStatsRetrySuccess(t *testing.T) {
	c := &config.Configuration{}
	config.SetDefaults(c)
	logger, _ := util.GetLogger(c)
	var shouldSucceed bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldSucceed {
			w.WriteHeader(500)
			shouldSucceed = true
			return
		}
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
				"views": 2,
				"relevanceScore": 2
			  }
			]
		  }`
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	c.StatsEndpoints = []string{ts.URL}
	defer ts.Close()
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=views&limit=3", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, len(apiResp.Data), 3)
	assert.Equal(t, apiResp.Count, 3)
	assert.Equal(t, apiResp.Data[0].Views, int32(11))
	assert.Equal(t, apiResp.Data[0].Url, "www.testurl3.com")

}

func TestGetStatsRetryError(t *testing.T) {
	c := &config.Configuration{}
	config.SetDefaults(c)
	logger, _ := util.GetLogger(c)
	var shouldSucceed int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldSucceed < 4 {
			w.WriteHeader(500)
			shouldSucceed += 1
			return
		}
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
				"views": 2,
				"relevanceScore": 2
			  }
			]
		  }`
		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
	c.StatsEndpoints = []string{ts.URL}
	defer ts.Close()
	api := NewV1API(c, nil, logger)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stats?sortKey=views&limit=3", nil)
	w := httptest.NewRecorder()

	api.handleGetStats()(w, req)
	resp := w.Result()
	respData, err := ioutil.ReadAll(resp.Body)
	assert.ErrorIs(t, err, nil)
	var apiResp APIErrorResponse
	err = json.Unmarshal(respData, &apiResp)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, apiResp.ErrorCode, 500)
}
