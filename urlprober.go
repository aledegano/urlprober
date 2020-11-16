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
	query          string
	requiredStatus map[int]bool
	tick           time.Duration
	url            string
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
		envInterval       = "interval"
		envPrefix         = "urlprober"
		envQuery          = "query"
		envRequiredStatus = "required_status"
		envURL            = "url"
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

	// Set the optional required status parameter from environment
	viper.BindEnv(envRequiredStatus)
	requiredStatus := viper.Get(envRequiredStatus)
	c.requiredStatus = make(map[int]bool)
	if requiredStatus == nil {
		c.requiredStatus[200] = true //Set the default to 200: OK
	} else {
		tokens := strings.Split(requiredStatus.(string), ",")
		for _, t := range tokens {
			status, err := strconv.Atoi(t)
			if err != nil {
				return fmt.Errorf("Cannot convert status code: %s to integer", t)
			}
			if http.StatusText(status) == "" {
				return fmt.Errorf("The provided status code: %d does not correspond to a valid HTTP status code", status)
			}
			c.requiredStatus[status] = true
		}
		log.Info().Msgf("The status codes: %s will be considered successful responses", tokens)
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
			if !c.requiredStatus[resp.StatusCode] {
				log.Warn().Str("url", c.url).Msgf("The returned status code: %d is not among those requested", resp.StatusCode)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Warn().Err(err).Str("url", c.url).Msg("Failed to decode response of query")
			}
			log.Info().
				Str("url", c.url).
				Str("status", strconv.Itoa(resp.StatusCode)).
				Msgf("Successful query response: %s", body)
		}
	}
}
