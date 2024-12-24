//go:build windows
// +build windows

package kmodule

import (
	grpc_server "bigagent/grpcs/server"
	"context"
	"log"

	"github.com/yusufpapurcu/wmi"
)

type Win32_SystemDriver struct {
	Name      string `json:"name"`
	State     string `json:"state"`
	StartMode string `json:"start_mode"`
}

func Info() ([]Win32_SystemDriver, error) {
	return InfoWithContext(context.Background())
}

func InfoWithContext(ctx context.Context) ([]Win32_SystemDriver, error) {
	var dst []Win32_SystemDriver
	query := "SELECT Name, State, StartMode FROM Win32_SystemDriver"

	err := wmi.Query(query, &dst)
	if err != nil {
		log.Fatal(err)
	}

	return dst, nil
}

type Kmodules map[string]Win32_SystemDriver

func NewKmodules() *Kmodules {
	k := make(Kmodules)

	kinfos, err := Info()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range kinfos {
		wininfo := Win32_SystemDriver{
			Name:      i.Name,
			State:     i.State,
			StartMode: i.StartMode,
		}

		k[i.Name] = wininfo
	}

	return &k
}

func NewKmodulesGrpc() *map[string]*grpc_server.Win32_SystemDriver {
	k := make(map[string]*grpc_server.Win32_SystemDriver)
	kinfos, err := Info()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range kinfos {
		wininfo := &grpc_server.Win32_SystemDriver{
			Name:      i.Name,
			State:     i.State,
			Startmode: i.StartMode,
		}

		k[i.Name] = wininfo
	}

	return &k
}
