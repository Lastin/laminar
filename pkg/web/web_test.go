package web

import (
	"github.com/digtux/laminar/pkg/shared"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestClient_StartWeb(t *testing.T) {
	type fields struct {
		logger      *zap.SugaredLogger
		PauseChan   chan time.Time
		BuildChan   chan DockerBuildJSON
		githubToken string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "foo",
			fields: fields{
				logger:      shared.GetLogger(true),
				PauseChan:   make(chan time.Time, 1),
				BuildChan:   make(chan DockerBuildJSON, 1),
				githubToken: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				logger:      tt.fields.logger,
				PauseChan:   tt.fields.PauseChan,
				BuildChan:   tt.fields.BuildChan,
				githubToken: tt.fields.githubToken,
			}
			client.StartWeb()
		})
	}
}
