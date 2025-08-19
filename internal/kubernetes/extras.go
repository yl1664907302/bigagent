package kubernetes

import (
	"context"
	"time"
)

// TimeoutContext 返回一个带超时的 context，便于 kube 操作统一控制时限
func TimeoutContext(parent context.Context, seconds int) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithTimeout(parent, time.Duration(seconds)*time.Second)
}
