//go:build linux
// +build linux

package kmodule

import (
	"bigagent/scrape/machine/formatsize"
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
)

// one module base info
type InfoStat struct {
	Name string   `json:"name"`
	Size string   `json:"size"`
	Used int      `json:"used"`
	Stat string   `json:"stat"`
	By   []string `json:"by"`
}

// get module info
func Info() ([]InfoStat, error) {
	return InfoWithContext(context.Background())
}

// get module info whit context
func InfoWithContext(ctx context.Context) ([]InfoStat, error) {
	// info is []InfoStat
	var ret []InfoStat

	// open file /proc/modules
	modInfo, err := os.Open("/proc/modules")
	if err != nil {
		return nil, err
	}

	// close
	defer modInfo.Close()

	// get buf read file
	scanner := bufio.NewScanner(modInfo)

	// scan
	for scanner.Scan() {
		line := scanner.Text()
		info := strings.Fields(line)
		size, err := strconv.ParseUint(info[1], 10, 64)
		if err != nil {
			return nil, err
		}

		used, err := strconv.Atoi(info[2])
		if err != nil {
			return nil, err
		}

		if info[3] == "-" {
			info[3] = ""
		}

		infostat := InfoStat{
			Name: info[0],
			Size: formatsize.FormatSize(size),
			Used: used,
			Stat: info[4],
			By:   strings.Fields(info[3]),
		}

		ret = append(ret, infostat)

	}

	return ret, nil
}

func (m InfoStat) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

type Kmodules map[string]InfoStat

func NewKmodules() *Kmodules {
	k := make(Kmodules)

	kinfos, err := Info()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range kinfos {
		infostat := InfoStat{
			Name: i.Name,
			Size: i.Size,
			Used: i.Used,
			Stat: i.Stat,
			By:   i.By,
		}

		k[i.Name] = infostat
	}

	return &k
}

func (k *Kmodules) ToString() string {
	b, err := json.MarshalIndent(k, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
