package shoutouthostnamegcp

import (
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
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		message := Get()
		http.PostForm(slackAPI, url.Values{"message": {message}, "slack_channel": {slackChannel}})
		log.Info().Str("slackapi", slackAPI).Str("slackchannel", slackChannel).Send()
		log.Info().Msg("catch the signal")
	}()
}
