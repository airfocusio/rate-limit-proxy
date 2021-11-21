package internal

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// RateLimitProxy ...
type RateLimitProxy struct {
	Config         RateLimitProxyConfig
	RedisClient    redis.Client
	Identifiers    []Identifier
	InnerServeHTTP func(http.ResponseWriter, *http.Request)
}

func (p *RateLimitProxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	enforce := len(p.Config.Paths.Includes) == 0
	for _, prefix := range p.Config.Paths.Includes {
		if strings.HasPrefix(req.URL.Path, prefix) {
			enforce = true
		}
	}
	for _, prefix := range p.Config.Paths.Excludes {
		if strings.HasPrefix(req.URL.Path, prefix) {
			enforce = false
		}
	}

	key := ""
	limit := p.Config.Limits.Anonymous
	count := int64(0)
	ttl := time.Duration(0)
	if enforce {
		for _, identifier := range p.Identifiers {
			if key == "" {
				identity, err := identifier.IdentifyRequest(req)
				if err == nil {
					key = "identified:" + *identity
					limit = p.Config.Limits.Identified
					for k, v := range p.Config.Limits.Other {
						if strings.HasPrefix(*identity, k) {
							limit = v
						}
					}
				}
			}
		}

		if key == "" {
			key = "anonymous:" + ExtractRequestClientIP(req)
		}

		key = p.Config.Redis.KeyPrefix + key
	} else {
		limit = 0
	}

	if limit > 0 {
		countRaw, err := p.RedisClient.Eval(ctx, `
			local current
			current = redis.call("incr",KEYS[1])
			if tonumber(current) == 1 then
				redis.call("expire",KEYS[1],ARGV[1])
			end
			return current
		`, []string{key}, (time.Duration(p.Config.Limits.Interval) * time.Second).Seconds()).Result()
		if err != nil {
			log.Printf("redis error: %v\n", err)
		} else {
			count = countRaw.(int64)
			ttl, err = p.RedisClient.PTTL(ctx, key).Result()
			if err != nil {
				log.Printf("redis error: %v\n", err)
			}
		}
	}
	remaining := limit - count

	if limit == 0 {
		p.InnerServeHTTP(wr, req)
	} else if remaining >= 0 {
		AddResponseRateLimitHeaders(wr, limit, remaining, ttl)
		p.InnerServeHTTP(wr, req)
	} else {
		AddResponseRateLimitHeaders(wr, limit, remaining, ttl)
		wr.WriteHeader(http.StatusTooManyRequests)
	}
}
