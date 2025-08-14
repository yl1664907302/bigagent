package polling

import (
	"bigagent/internal/check/result"
)

type polling interface {
	Check() (result.Result, error)
}
