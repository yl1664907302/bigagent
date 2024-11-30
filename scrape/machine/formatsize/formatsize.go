package formatsize

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func FormatSize[T constraints.Integer](size T) string {

	sizeInBytes := float64(size)

	if sizeInBytes < 1024 {
		return fmt.Sprintf("%.2f B", sizeInBytes)
	}
	if sizeInBytes < 1024*1024 {
		return fmt.Sprintf("%.2f KB", sizeInBytes/1024)
	}
	if sizeInBytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", sizeInBytes/(1024*1024))
	} else {
		return fmt.Sprintf("%.2f GB", sizeInBytes/(1024*1024*1024))
	}
}

func FormatPercent(size float64) string {
	return fmt.Sprintf("%.2f%%", size)
}
