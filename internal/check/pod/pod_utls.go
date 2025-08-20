package check

import (
	"context"
	"time"
)

func GetK8sTwContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
}

func isBadWaitingReason(reason string) bool {
	switch reason {
	case "CrashLoopBackOff", "ImagePullBackOff", "ErrImagePull", "CreateContainerError":
		return true
	default:
		return false
	}
}
