package diskinfo

import (
	"github.com/shirou/gopsutil/v4/disk"
	"log"
)

// diskinfo
type DiskInfo struct {
	Usage      *disk.UsageStat
	Partition  []disk.PartitionStat
	IOCounters map[string]disk.IOCountersStat
}

func GetUsage() *disk.UsageStat {
	usage, err := disk.Usage("/")
	if err != nil {
		log.Println(err)
		return nil
	}
	return usage
}

func GetPartition() []disk.PartitionStat {
	partitions, err := disk.Partitions(false)
	if err != nil {
		log.Println(err)
		return nil
	}
	return partitions
}

func GetIOCounters() map[string]disk.IOCountersStat {
	counters, err := disk.IOCounters()
	if err != nil {
		log.Println(err)
		return nil
	}
	return counters
}
