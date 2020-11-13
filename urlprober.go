package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type config struct {
	url  string
	tick time.Duration
}

func (c *config) init() error {
	viper.SetEnvPrefix("dyndns")

	// Set the URL to probe from environment
	viper.BindEnv("url")
	url := viper.Get("url")
	if url == nil {
		log.Fatal().
			Msg("The enviornment variable DYNDNS_URL must be set.")
	}
	c.url = url.(string)

	// Set the probing interval, in seconds, from the environment
	viper.BindEnv("interval")
	interval := viper.Get("interval")
	if interval == nil {
		log.Fatal().Msg("The environment variable DYNDNS_INTERVAL must be set.")
	}
	tick, err := strconv.Atoi(interval.(string))
	if err != nil {
		log.Fatal().Msgf("Cannot convert interval: %s to integer.", interval.(string))
	}
	c.tick = time.Duration(tick) * time.Second

	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := &config{}

	defer func() {
		cancel()
	}()

	if err := run(ctx, c); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, c *config) error {
	c.init()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.Tick(c.tick):
			log.Info().Str("url", c.url).Msg("Querying.")
			resp, err := http.Get(c.url)
			if err != nil {
				log.Warn().Err(err).Str("url", c.url).Msg("The query failed.")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Warn().Err(err).Str("url", c.url).Msg("Failed to decode response of query")
			}
			log.Info().Str("url", c.url).Msgf("Successful query response: %s", body)
		}
	}
}
