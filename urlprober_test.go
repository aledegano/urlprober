package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

func clearEnv() {
	for _, env := range []string{"URL", "INTERVAL", "QUERY"} {
		os.Unsetenv("URLPROBER_" + env)
	}
}

func TestMissingUrl(t *testing.T) {
	clearEnv()
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should fail when URLPROBER_URL is not set.")
	}
}

func TestMissingInterval(t *testing.T) {
	clearEnv()
	os.Setenv("URLPROBER_URL", "https://example.com")
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should fail when URLPROBER_INTERVAL is not set.")
	}
}

func TestWrongInterval(t *testing.T) {
	clearEnv()
	os.Setenv("URLPROBER_INTERVAL", "not_a_number")
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should fail when URLPROBER_INTERVAL cannot be interpreted as an integer.")
	}
}

func TestInvalidQuery(t *testing.T) {
	clearEnv()
	os.Setenv("URLPROBER_URL", "https://example.com")
	os.Setenv("URLPROBER_INTERVAL", "42")
	os.Setenv("URLPROBER_QUERY", "foo\bar")
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should fail when URLPROBER_QUERY forms an invalid URL.")
	}
}

func TestCorrectInit(t *testing.T) {
	clearEnv()
	os.Setenv("URLPROBER_URL", "https://example.com")
	os.Setenv("URLPROBER_INTERVAL", "42")
	c := &config{}
	err := c.init()
	if err != nil {
		log.Fatal().Msg("Should have not failed.")
	}
}
