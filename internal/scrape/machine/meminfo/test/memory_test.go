package test

import (
	"bigagent/scrape/machine/memory"
	"fmt"
	"testing"
)

func TestMemory_Total(t *testing.T) {
	total := memory.NewMemory().UsedPercent()
	fmt.Println(total)
}
