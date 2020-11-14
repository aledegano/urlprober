package main

import (
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestMissingUrl(t *testing.T) {
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should have failed.")
	}
}

func TestMissingInterval(t *testing.T) {
	os.Setenv("URLPROBER_URL", "https://example.com")
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should have failed.")
	}
}

func TestWrongInterval(t *testing.T) {
	os.Setenv("URLPROBER_INTERVAL", "not_a_number")
	c := &config{}
	err := c.init()
	if err == nil {
		log.Fatal().Msg("Should have failed.")
	}
}

func TestCorrectInit(t *testing.T) {
	os.Setenv("URLPROBER_URL", "https://example.com")
	os.Setenv("URLPROBER_INTERVAL", "42")
	c := &config{}
	err := c.init()
	if err != nil {
		log.Fatal().Msg("Should have not failed.")
	}
}
