package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/airfocusio/rate-limit-proxy/internal"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-cmp/cmp"
)

func main() {
	configPath := flag.String("config", "config.yaml", "The config file path.")
	listen := flag.String("listen", "0.0.0.0:8080", "The interface to listen on.")
	target := flag.String("target", "http://127.0.0.1:80", "The target to proxy to.")
	flag.Parse()

	config, identifiers, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("unable to load config: %v\n", err)
		return
	}
	redisClient, err := createRedisClient(*config)
	if err != nil {
		log.Fatalf("unable to create redis client: %v\n", err)
		return
	}
	proxyTarget, err := createProxyTarget(*target)
	if err != nil {
		log.Fatalf("unable to create proxy target: %v\n", err)
	}
	proxy := &internal.RateLimitProxy{
		Config:         *config,
		RedisClient:    *redisClient,
		Identifiers:    *identifiers,
		InnerServeHTTP: proxyTarget.ServeHTTP,
	}

	go reloadConfigLoop(*configPath, proxy)

	log.Printf("starting proxy server on %s to %s\n", *listen, *target)
	defer log.Printf("starting proxy server on %s\n", *listen)
	if err := http.ListenAndServe(*listen, proxy); err != nil {
		log.Fatalf("unable to start server on %s: %v\n", *listen, err)
	}
}

func loadConfig(configPath string) (*internal.RateLimitProxyConfig, *[]internal.Identifier, error) {
	configYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, nil, err
	}
	config, identifiers, err := internal.LoadRateLimitProxyConfig(configYaml)
	if err != nil {
		return nil, nil, err
	}
	return config, identifiers, nil
}

func reloadConfigLoop(configPath string, proxy *internal.RateLimitProxy) {
	for {
		time.Sleep(10 * time.Second)
		config, identifiers, err := loadConfig(configPath)
		if err != nil {
			log.Printf("unable to reload config: %v\n", err)
			continue
		}

		if config.Equal(proxy.Config) {
			continue
		}

		redisClient := &proxy.RedisClient
		if !cmp.Equal(proxy.Config.Redis, config.Redis) {
			redisClient, err = createRedisClient(*config)
			if err != nil {
				log.Printf("unable to reload config: %v\n", err)
				continue
			}
			proxy.RedisClient.Close()
		}

		proxy.Config = *config
		proxy.RedisClient = *redisClient
		proxy.Identifiers = *identifiers
		log.Printf("reloaded config\n")
	}
}

func createRedisClient(config internal.RateLimitProxyConfig) (*redis.Client, error) {
	var tlsConfig *tls.Config
	if config.Redis.TLS {
		tlsConfig = &tls.Config{}
	}
	client := redis.NewClient(&redis.Options{
		Addr:         config.Redis.Address,
		Password:     config.Redis.Password,
		TLSConfig:    tlsConfig,
		MaxRetries:   1,
		DialTimeout:  time.Duration(time.Duration.Milliseconds(1000)),
		ReadTimeout:  time.Duration(time.Duration.Milliseconds(1000)),
		WriteTimeout: time.Duration(time.Duration.Milliseconds(1000)),
	})
	return client, nil
}

func createProxyTarget(targetStr string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetStr)
	if err != nil {
		return nil, err
	}
	target := httputil.NewSingleHostReverseProxy(url)
	return target, nil
}
