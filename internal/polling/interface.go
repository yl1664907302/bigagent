package polling

import "bigagent/internal/check"

type polling interface {
	Check() (check.Result, error)
}
