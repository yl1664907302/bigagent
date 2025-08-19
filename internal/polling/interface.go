package polling

import (
	"bigagent/internal/check/result"
)

type PollFunc func(args ...interface{}) result.Result
