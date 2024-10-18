package test

import (
	"bigagent/scrape/machine/info"
	"fmt"
	"testing"
)

func TestInfo_Platform(t *testing.T) {
	platform := info.NewInfo().PlatformVersion()
	fmt.Println(platform)
}
