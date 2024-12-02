package diskinfo

import (
	"bigagent/scrape/machine/formatsize"
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/v4/disk"
)

type smpInfo struct {
	Path        string `json:"path"`
	Total       string `json:"total"`
	Free        string `json:"free"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedpercent"`
	Device      string `json:"device"`
	Fstype      string `json:"fstype"`
	MountPoint  string `json:"mountpoint"`
}

type SmpDisk map[string]smpInfo

func NewSmpDisk() *SmpDisk {
	// 初始化smpdisk
	smpdisk := make(SmpDisk)

	// 获取disk信息
	diskInfo, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range diskInfo {
		usage, err := disk.Usage(i.Mountpoint)
		if err != nil {
			continue
		}

		// smpdisk[i.Device] = make(map[string]smpInfo)

		smpinfo := smpInfo{
			Path:        usage.Path,
			Total:       formatsize.FormatSize(usage.Total),
			Free:        formatsize.FormatSize(usage.Free),
			Used:        formatsize.FormatSize(usage.Used),
			UsedPercent: formatsize.FormatPercent(usage.UsedPercent),
			Device:      i.Device,
			Fstype:      usage.Fstype,
			MountPoint:  i.Mountpoint,
		}
		smpdisk[i.Mountpoint] = smpinfo

	}

	return &smpdisk
}

func (d *SmpDisk) ToString() string {
	b, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
