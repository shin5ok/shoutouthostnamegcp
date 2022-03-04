package shoutouthostnamegcp

import (
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	hostname, _ := os.Hostname()
	tests := []struct {
		name string
		want string
	}{
		{"foo", hostname},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Get(); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetSigHandler(t *testing.T) {
	type args struct {
		slackAPI     string
		slackChannel string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "test", args: args{slackAPI: os.Getenv("SLACK_URL"), slackChannel: "#" + os.Getenv("SLACK_CHANNEL")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSigHandler(tt.args.slackAPI, tt.args.slackChannel)
		})
	}
}
