package test

import (
	model "bigagent/model/machine"
	"fmt"
	"testing"
)

func TestNewData(t *testing.T) {
	data := model.NewData()
	fmt.Println(data.ToString())
}
