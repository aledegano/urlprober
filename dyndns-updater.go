package main

import (
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("dyndns")
	viper.BindEnv("url")
	url := viper.Get("url")
	if url == nil {
		log.Fatal().
			Msg("The enviornment variable DYNDNS_URL must be set.")
	}

	log.Info().Msg("Querying Dyn-dns to update IP.")
	resp, err := http.Get(url.(string))
	if err != nil {
		log.Fatal().Err(err).Msg("The query failed.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to decode response of query.")
	}
	log.Info().Msgf("Successful query, response: %s", body)
}
