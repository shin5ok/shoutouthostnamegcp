package shoutouthostnamegcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/rs/zerolog"

	// "github.com/rs/zerolog/log"
	log "github.com/rs/zerolog/log"
	// "google.golang.org/grpc/metadata"
	// "github.com/gin-gonic/gin"
)

func init() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func Get() string {
	if metadata.OnGCE() {
		hostname, _ := metadata.Hostname()
		return hostname
	}
	hostname, _ := os.Hostname()
	return hostname
}

func SetSigHandler(slackAPI, slackChannel string) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	type slackStruct struct {
		Payload string `json:"text"`
		Channel string `json:"channel"`
	}

	go func() {
		s := <-sigs
		message := fmt.Sprintf("%s on %s", s.String(), Get())
		log.Info().Msg(message)
		postData := slackStruct{}
		postData.Payload = message
		postData.Channel = "#" + slackChannel
		postDataJson, _ := json.Marshal(&postData)
		resp, err := http.PostForm(slackAPI, url.Values{"payload": {string(postDataJson)}})
		if err != nil {
			log.Error().Err(err)
		}
		log.Info().Str("slackapi", slackAPI).Str("slackchannel", slackChannel).Send()
		log.Info().Msgf("%+v", resp)
		os.Exit(0)
	}()
}
