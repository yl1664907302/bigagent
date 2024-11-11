package diskinfo

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v4/disk"
	"log"
)

type DISK struct {
	U  *disk.UsageStat                `json:"disk_usage"`
	P  []disk.PartitionStat           `json:"disk_partitions"`
	IO map[string]disk.IOCountersStat `json:"disk_io"`
}
type Options func(d *DISK)

func (d *DISK) ToString() string {

	indent, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return err.Error()
	}
	return string(indent)
}

func WithIO(io map[string]disk.IOCountersStat, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(d *DISK) {
		d.IO = io
	}
}

func WithPartition(par []disk.PartitionStat, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(d *DISK) {
		d.P = par
	}
}

func WithUsage(use *disk.UsageStat, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(d *DISK) {
		d.U = use
	}
}

func NewDISK(options ...Options) *DISK {

	par, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}
	usage, err := disk.Usage("/")
	if err != nil {
		log.Fatal(err)
	}

	counters, err := disk.IOCounters()
	if err != nil {
		log.Fatal(err)
	}
	disk := &DISK{
		U:  usage,
		P:  par,
		IO: counters,
	}
	for _, option := range options {
		option(disk)
	}
	return disk
}
