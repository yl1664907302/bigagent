package polling

import (
	"bigagent/internal/result"
)

type PollFunc func(args ...interface{}) result.Result
