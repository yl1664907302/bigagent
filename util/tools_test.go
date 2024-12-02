package utils

import (
	"testing"
)

func Test_replaceKeyValue(t *testing.T) {
	err := ReplaceKeyValue("../config.yml", "grpc_cmdb1_stand1", "0.0.0.0:7777")
	if err != nil {
		DefaultLogger.Error(err)
	}
}
