package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
		return errors.New("The enviornment variable DYNDNS_URL must be set")
	}
	c.url = url.(string)

	// Set the probing interval, in seconds, from the environment
	viper.BindEnv("interval")
	interval := viper.Get("interval")
	if interval == nil {
		return errors.New("The environment variable DYNDNS_INTERVAL must be set")
	}
	tick, err := strconv.Atoi(interval.(string))
	if err != nil {
		return fmt.Errorf("Cannot convert interval: %s to integer", interval.(string))
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
		log.Fatal().Str("component", "init").Err(err).Send()
	}
}

func run(ctx context.Context, c *config) error {
	log.Info().Str("component", "run").Msg("Starting prober.")
	err := c.init()
	if err != nil {
		return err
	}

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
