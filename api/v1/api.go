package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/renjugeo/go-server/config"
	"github.com/renjugeo/go-server/util"
	"go.uber.org/zap"
)

const (
	pathPrefix = "/api/v1"

	SortKey      = "sortKey"
	Limit        = "limit"
	DefaultLimit = "200"
	SortKeyScore = "relevanceScore"
	SortKeyViews = "views"

	ErrorInvalidLimit   = "invalid value for limit"
	ErrorInvalidSortKey = "sortKey should be either relevanceScore or views"
)

type API struct {
	logger     *zap.Logger
	cfg        *config.Configuration
	httpclient *retryablehttp.Client
}

func NewV1API(cfg *config.Configuration, logger *zap.Logger) *API {
	if logger == nil {
		logger, _ = util.GetLogger(cfg)
	}
	cli := retryablehttp.NewClient()
	cli.RetryMax = 3
	return &API{
		logger:     logger,
		cfg:        cfg,
		httpclient: cli,
	}
}

func (api *API) SetHttpClient(client *http.Client) {
	api.httpclient.HTTPClient = client
}

func (api *API) RegisterPaths(r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("%s/stats", pathPrefix), api.handleGetStats()).Methods(http.MethodGet)
}

func (api *API) handleGetStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get and validate relevanceScore
		key := r.URL.Query().Get(SortKey)
		if key == "" {
			key = SortKeyScore
		}
		if key != SortKeyScore && key != SortKeyViews {
			api.jsonResponse(w, APIErrorResponse{Message: ErrorInvalidSortKey, ErrorCode: 500})
			return
		}
		// get and validate views
		limit := r.URL.Query().Get(Limit)
		if limit == "" {
			limit = DefaultLimit
		}
		l, err := strconv.Atoi(limit)
		if err != nil {
			api.jsonResponse(w, APIErrorResponse{Message: ErrorInvalidLimit, ErrorCode: 500})
			return
		}
		resp, err := api.getStats(key, l)
		if err != nil {
			r := APIErrorResponse{
				Message:   err.Error(),
				ErrorCode: 500,
			}
			api.jsonResponse(w, r)
			return
		}
		api.jsonResponse(w, resp)
	}
}

func (api *API) jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		w.WriteHeader(500)
	}
}

func (api *API) getStats(key string, limit int) (*APIResponse, error) {
	urls := api.cfg.StatsEnpoints
	st, err := api.fetchUrls(urls)
	if err != nil {
		return nil, err
	}

	aggregated := []EndpointStat{}
	for _, v := range st {
		aggregated = append(aggregated, v.Data...)
	}

	if key == SortKeyScore {
		sort.Slice(aggregated, func(i, j int) bool {
			return aggregated[i].Score > aggregated[j].Score
		})
	} else if key == SortKeyViews {
		sort.Slice(aggregated, func(i, j int) bool {
			return float32(aggregated[i].Views) > float32(aggregated[j].Views)
		})
	}
	if limit < len(aggregated) {
		aggregated = aggregated[:limit]
	}
	r := &APIResponse{
		Stats: Stats{
			Data: aggregated,
		},
		Count: len(aggregated),
	}

	return r, nil
}

func (api *API) fetchUrls(urls []string) ([]Stats, error) {
	var wg sync.WaitGroup
	var result []Stats
	resultMap := make(map[string]Stats)
	rw := &sync.RWMutex{}

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			stat, err := api.fetchUrl(url)
			if err != nil {
				api.logger.Debug("error curling", zap.Error(err))
				return
			}
			rw.Lock()
			defer rw.Unlock()
			resultMap[url] = *stat

		}(urls[i])
	}
	wg.Wait()
	for _, url := range urls {
		data, ok := resultMap[url]
		if !ok {
			return nil, errors.New(fmt.Sprintf("could not get data from %s", url))
		}
		result = append(result, data)
	}

	return result, nil
}

func (api *API) fetchUrl(url string) (*Stats, error) {
	req, err := retryablehttp.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating a new request")
	}
	req.Header.Add("Accept", "application/json")
	resp, err := api.httpclient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error occurred when making http request")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read response body")
	}

	var stats Stats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, errors.Wrap(err, "unable to decode response body")
	}

	return &stats, nil
}
