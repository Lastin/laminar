package cmd

import (
	"time"
)

// ChangeRequests is an object recording what changed and when
type ChangeRequest struct {
	Old          string    `json:"old"`
	New          string    `json:"new"`
	Time         time.Time `json:"time"`
	PatternValue string    `json:"patternValue"`
	PatternType  string    `json:"patternType"`
	Image        string    `json:"image"`
	File         string    `json:"file"`
}
