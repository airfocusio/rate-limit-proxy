package internal

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

type testIdentifier struct{}

func (t testIdentifier) IdentifyRequest(req *http.Request) (*string, error) {
	id := req.URL.Query().Get("id")
	if id == "" {
		return nil, fmt.Errorf("anonymous")
	}
	return &id, nil
}

func TestRateLimitProxy(t *testing.T) {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
	proxy := RateLimitProxy{
		Config: RateLimitProxyConfig{
			Redis: RateLimitProxyConfigRedis{},
			Limits: RateLimitProxyConfigLimits{
				Interval:   60,
				Anonymous:  1,
				Identified: 2,
				Other: map[string]int64{
					"system": 0,
				},
			},
			Paths: RateLimitProxyConfigPaths{
				Includes: []string{"/api/"},
				Excludes: []string{"/api/unlimited/"},
			},
		},
		RedisClient: *redisClient,
		Identifiers: []Identifier{testIdentifier{}},
		InnerServeHTTP: func(wr http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/notfound/" || req.URL.Path == "/api/notfound/" {
				wr.WriteHeader(http.StatusNotFound)
			} else {
				wr.WriteHeader(http.StatusOK)
			}
		},
	}
	resetCounters := func() {
		proxy.Config.Redis.KeyPrefix = strconv.FormatInt(seededRand.Int63(), 10) + ":"
	}

	test := func(path string, token string, expectedCode int, expectedLimit int64, expectedRemaining int64) {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if token != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
		if w.Code != expectedCode {
			t.Errorf("expected code %d, got %d", expectedCode, w.Code)
		}
		limit := w.Header().Get("X-Ratelimit-Limit")
		if expectedLimit == -1 && limit != "" {
			t.Errorf("expected no limit")
		}
		if expectedLimit != -1 && limit != strconv.FormatInt(expectedLimit, 10) {
			t.Errorf("expected limit %d, got %s", expectedLimit, limit)
		}
		remaining := w.Header().Get("X-Ratelimit-Remaining")
		if expectedRemaining == -1 && remaining != "" {
			t.Errorf("expected no remaining")
		}
		if expectedRemaining != -1 && remaining != strconv.FormatInt(expectedRemaining, 10) {
			t.Errorf("expected remaining %d, got %s", expectedRemaining, remaining)
		}
	}

	t.Run("anonymous", func(t *testing.T) {
		resetCounters()
		test("/api/", "", 200, 1, 0)
		test("/api/", "", 429, 1, 0)
	})

	t.Run("identified", func(t *testing.T) {
		resetCounters()
		test("/api/?id=1", jwtHs256User1, 200, 2, 1)
		test("/api/?id=2", jwtHs256User2, 200, 2, 1)
		test("/api/?id=1", jwtHs256User1, 200, 2, 0)
		test("/api/?id=2", jwtHs256User2, 200, 2, 0)
		test("/api/?id=1", jwtHs256User1, 429, 2, 0)
		test("/api/?id=2", jwtHs256User2, 429, 2, 0)
	})

	t.Run("included and excluded routes", func(t *testing.T) {
		resetCounters()
		test("/api/unlimited/", "", 200, -1, -1)
		test("/other/", "", 200, -1, -1)
	})

	t.Run("underyling server returns not found", func(t *testing.T) {
		resetCounters()
		test("/api/notfound/", "", 404, 1, 0)
		test("/api/notfound/", "", 429, 1, 0)
		test("/notfound/", "", 404, -1, -1)
		test("/notfound/", "", 404, -1, -1)
	})

	t.Run("identity with special limits", func(t *testing.T) {
		resetCounters()
		test("/api/", "", 200, 1, 0)
		test("/api/?id=1", "", 200, 2, 1)
		test("/api/?id=system:1", "", 200, -1, -1)
		test("/api/?id=system:2", "", 200, -1, -1)
	})
}
