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
	log "github.com/rs/zerolog/log"
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

	go func() {

		type slackStruct struct {
			Payload string `json:"text"`
			Channel string `json:"channel"`
		}

		s := <-sigs
		message := fmt.Sprintf("%s:%s", s.String(), Get())

		log.Info().Msg(message)
		postData := slackStruct{}
		postData.Payload = message
		postData.Channel = "#" + slackChannel
		postDataJson, err := json.Marshal(&postData)
		if err != nil {
			log.Error().Err(err)
			postDataJson = []byte(fmt.Sprintf("{'text':'%s', 'channel':'#%s'}", err, slackChannel))
		}
		resp, err := http.PostForm(slackAPI, url.Values{"payload": {string(postDataJson)}})
		if err != nil {
			log.Error().Err(err)
		}
		log.Info().Str("slackapi", slackAPI).Str("slackchannel", slackChannel).Send()
		log.Info().Msgf("%+v", resp)
		os.Exit(0)
	}()
}
