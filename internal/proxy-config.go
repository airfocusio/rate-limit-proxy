package internal

import (
	"fmt"
	"strings"

	"github.com/airfocusio/go-expandenv"
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
func LoadRateLimitProxyConfig(bytesRaw []byte) (*RateLimitProxyConfig, *[]Identifier, error) {
	var expansionTemp interface{}
	err := yaml.Unmarshal(bytesRaw, &expansionTemp)
	if err != nil {
		return nil, nil, err
	}
	expansionTemp, err = expandenv.ExpandEnv(expansionTemp)
	if err != nil {
		return nil, nil, err
	}
	bytes, err := yaml.Marshal(expansionTemp)
	if err != nil {
		return nil, nil, err
	}

	config := &RateLimitProxyConfig{}
	if err := yaml.Unmarshal(bytes, config); err != nil {
		return nil, nil, err
	}

	identifiers := []Identifier{}
	for i, c := range config.IdentifiersConfig {
		if c.JwtBearerHeader != nil {
			claimNames := []string{}
			if c.JwtBearerHeader.Claim != "" {
				claimNames = strings.Split(c.JwtBearerHeader.Claim, "|")
			}
			identifiers = append(identifiers, JwtIdentifier{
				Algorithm:         c.JwtBearerHeader.Algorithm,
				KeyID:             c.JwtBearerHeader.KeyID,
				Verifier:          c.JwtBearerHeader.Verifier,
				TokenExtractor:    ExtractBearerToken,
				IdentityExtractor: ExtractClaim(claimNames),
			})
		} else if c.JwtQueryToken != nil {
			claimNames := []string{}
			if c.JwtQueryToken.Claim != "" {
				claimNames = strings.Split(c.JwtQueryToken.Claim, "|")
			}
			identifiers = append(identifiers, JwtIdentifier{
				Algorithm:         c.JwtQueryToken.Algorithm,
				KeyID:             c.JwtQueryToken.KeyID,
				Verifier:          c.JwtQueryToken.Verifier,
				TokenExtractor:    ExtractQueryParameter(c.JwtQueryToken.Name),
				IdentityExtractor: ExtractClaim(claimNames),
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
