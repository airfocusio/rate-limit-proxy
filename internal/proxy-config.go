package internal

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
)

// RateLimitProxyConfigLimits ...
type RateLimitProxyConfigLimits struct {
	Interval   int64            `yaml:"interval"`
	Anonymous  int64            `yaml:"anonymous"`
	Identified int64            `yaml:"identified"`
	Other      map[string]int64 `yaml:"other"`
}

// RateLimitProxyConfigRedis ...
type RateLimitProxyConfigRedis struct {
	Address   string `yaml:"address"`
	Password  string `yaml:"password"`
	TLS       bool   `yaml:"tls"`
	KeyPrefix string `yaml:"keyPrefix"`
}

// RateLimitProxyConfigPaths ...
type RateLimitProxyConfigPaths struct {
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
}

// RateLimitProxyConfig ...
type RateLimitProxyConfig struct {
	Redis             RateLimitProxyConfigRedis  `yaml:"redis"`
	Paths             RateLimitProxyConfigPaths  `yaml:"paths"`
	Limits            RateLimitProxyConfigLimits `yaml:"limits"`
	IdentifiersConfig []IdentifierConfig         `yaml:"identifiers"`
}

// LoadRateLimitProxyConfig ...
func LoadRateLimitProxyConfig(bytes []byte) (*RateLimitProxyConfig, *[]Identifier, error) {
	config := &RateLimitProxyConfig{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
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
