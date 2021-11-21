package internal

import (
	utiljson "encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// RateLimitProxyConfigLimits ...
type RateLimitProxyConfigLimits struct {
	Interval   int64            `json:"interval"`
	Anonymous  int64            `json:"anonymous"`
	Identified int64            `json:"identified"`
	Other      map[string]int64 `json:"other"`
}

// RateLimitProxyConfigRedis ...
type RateLimitProxyConfigRedis struct {
	Address   string `json:"address"`
	Password  string `json:"password"`
	TLS       bool   `json:"tls"`
	KeyPrefix string `json:"keyPrefix"`
}

// RateLimitProxyConfigPaths ...
type RateLimitProxyConfigPaths struct {
	Includes []string `json:"includes"`
	Excludes []string `json:"excludes"`
}

// RateLimitProxyConfig ...
type RateLimitProxyConfig struct {
	Redis             RateLimitProxyConfigRedis  `json:"redis"`
	Paths             RateLimitProxyConfigPaths  `json:"paths"`
	Limits            RateLimitProxyConfigLimits `json:"limits"`
	IdentifiersConfig []IdentifierConfig         `json:"identifiers"`
}

// LoadRateLimitProxyConfig ...
func LoadRateLimitProxyConfig(yaml []byte) (*RateLimitProxyConfig, *[]Identifier, error) {
	config := &RateLimitProxyConfig{}

	json, err := utilyaml.ToJSON(yaml)
	if err != nil {
		return nil, nil, err
	}

	if err = utiljson.Unmarshal(json, config); err != nil {
		return nil, nil, err
	}

	identifiers := []Identifier{}
	for i, c := range config.IdentifiersConfig {
		if c.JwtBearerHeader != nil {
			identifiers = append(identifiers, JwtIdentifier{
				Algorithm:         c.JwtBearerHeader.Algorithm,
				KeyID:             c.JwtBearerHeader.KeyID,
				Verifier:          c.JwtBearerHeader.Verifier,
				TokenExtractor:    ExtractBearerToken,
				IdentityExtractor: ExtractClaim(c.JwtBearerHeader.Claim),
			})
		} else if c.JwtQueryToken != nil {
			identifiers = append(identifiers, JwtIdentifier{
				Algorithm:         c.JwtQueryToken.Algorithm,
				KeyID:             c.JwtQueryToken.KeyID,
				Verifier:          c.JwtQueryToken.Verifier,
				TokenExtractor:    ExtractQueryParameter(c.JwtQueryToken.Name),
				IdentityExtractor: ExtractClaim(c.JwtQueryToken.Claim),
			})
		} else {
			return nil, nil, fmt.Errorf("identifier #%d is invalid", i)
		}
	}

	return config, &identifiers, nil
}

// Equal ...
func (c1 RateLimitProxyConfig) Equal(c2 RateLimitProxyConfig) bool {
	return true &&
		cmp.Equal(c1.Redis, c2.Redis) &&
		cmp.Equal(c1.Paths, c2.Paths) &&
		cmp.Equal(c1.Limits, c2.Limits) &&
		cmp.Equal(c1.IdentifiersConfig, c2.IdentifiersConfig)
}
