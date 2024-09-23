package machine

import (
	"bigagent/collection/machine"
	"bigagent/util"
	"fmt"
	"testing"
)

func TestGetLocalIP(t *testing.T) {
	ipv4, err := machine.GetLocalIPv4()
	if err != nil {
		util.Log.Println(err.Error())
	}
	fmt.Println(ipv4)
}
