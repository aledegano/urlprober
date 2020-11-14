package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

type config struct {
	url   string
	query string
	tick  time.Duration
}

func errorEnvNotSet(prefix string, env string) error {
	return fmt.Errorf(
		"The enviornment variable %s_%s must be set",
		strings.ToUpper(prefix),
		strings.ToUpper(env),
	)
}

func (c *config) init() error {
	var (
		envPrefix   = "urlprober"
		envURL      = "url"
		envInterval = "interval"
		envQuery    = "query"
	)
	viper.SetEnvPrefix(envPrefix)

	// Set the URL to probe from environment
	viper.BindEnv(envURL)
	targetURL := viper.Get(envURL)
	if targetURL == nil {
		return errorEnvNotSet(envPrefix, envURL)
	}
	c.url = targetURL.(string)

	// Set the probing interval, in seconds, from the environment
	viper.BindEnv(envInterval)
	interval := viper.Get(envInterval)
	if interval == nil {
		return errorEnvNotSet(envPrefix, envInterval)
	}
	tick, err := strconv.Atoi(interval.(string))
	if err != nil {
		return fmt.Errorf("Cannot convert interval: %s to integer", interval.(string))
	}
	c.tick = time.Duration(tick) * time.Second

	// Set the optional query paramter from the environment
	viper.BindEnv(envQuery)
	query := viper.Get(envQuery)
	if query == nil {
		log.Info().Msg("Optional query not set.")
	} else {
		c.query = query.(string)
		_, err := url.Parse(c.url + c.query)
		if err != nil {
			return errors.New("The provided query forms an invalid URL when appended to the provided target URL")
		}
	}

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
			resp, err := http.Get(c.url + c.query)
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
