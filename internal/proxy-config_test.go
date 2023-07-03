package internal

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadRateLimitProxyConfig(t *testing.T) {
	yaml := `redis:
  address: 127.0.0.1:50002
identifiers:
- jwtBearerHeader:
    keyId: '1'
    algorithm: HS256
    verifier: v1
- jwtQueryParameter:
    algorithm: ES256
    verifier: v2
    name: token
    claim: foo+bar
limits:
  interval: 60
  anonymous: 1
  identified: 10
  other:
    special: 100
paths:
  includes:
  - /api/
  excludes:
  - /api/unlimited/
`

	c1, i1, err := LoadRateLimitProxyConfig([]byte(yaml))
	if !assert.NoError(t, err) {
		return
	}
	c2 := RateLimitProxyConfig{
		Redis: RateLimitProxyConfigRedis{
			Address:  "127.0.0.1:50002",
			Password: "",
		},
		IdentifiersConfig: []IdentifierConfig{
			{JwtBearerHeader: &IdentifierConfigJwtBearerHeader{
				KeyID:     "1",
				Algorithm: "HS256",
				Verifier:  "v1",
			}},
			{JwtQueryToken: &IdentifierConfigJwtQueryParameter{
				KeyID:     "",
				Algorithm: "ES256",
				Verifier:  "v2",
				Name:      "token",
				Claim:     "foo+bar",
			}},
		},
		Limits: RateLimitProxyConfigLimits{
			Interval:   60,
			Anonymous:  1,
			Identified: 10,
			Other: map[string]int64{
				"special": 100,
			},
		},
		Paths: RateLimitProxyConfigPaths{
			Includes: []string{"/api/"},
			Excludes: []string{"/api/unlimited/"},
		},
	}
	if !assert.Equal(t, c2, *c1) {
		return
	}

	i11, ok := (*i1)[0].(JwtIdentifier)
	if !ok {
		assert.Fail(t, "expected JwtIdentifier, got %v", (*i1)[0])
		return
	}
	r1 := httptest.NewRequest(http.MethodGet, "/", nil)
	r1.Header.Add("Authorization", "Bearer token1")
	t1, err := i11.TokenExtractor(*r1)
	if err != nil || *t1 != "token1" {
		assert.Fail(t, "expected JwtIdentifier for bearer token header, got %v", (*i1)[0])
		return
	}
	i12, ok := (*i1)[1].(JwtIdentifier)
	if !ok {
		assert.Fail(t, "expected JwtIdentifier, got %v", (*i1)[1])
		return
	}
	r2 := httptest.NewRequest(http.MethodGet, "/?token=token2", nil)
	t2, err := i12.TokenExtractor(*r2)
	if err != nil || *t2 != "token2" {
		assert.Fail(t, "expected JwtIdentifier for bearer token header, got %v", (*i1)[1])
		return
	}
}
