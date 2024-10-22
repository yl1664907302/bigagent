package diskinfo

import (
	"context"
	"encoding/json"
	"github.com/shirou/gopsutil/v4/disk"
	"log"
	"strings"
)

// diskinfo
type DiskInfo struct {
	Usage      *disk.UsageStat
	Partition  *disk.PartitionStat
	IOCounters *disk.IOCountersStat
}

// GetDiskPath Disk 路径信息
func (d *DiskInfo) GetDiskPath() ([]string, error) {
	var paths []string
	partition, err := d.GetPartition()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, v := range partition {
		paths = append(paths, v.Device)
	}
	return paths, nil
}

// GetUsage Disk信息
func (d *DiskInfo) GetUsage() (*disk.UsageStat, error) {
	paths, err := d.GetDiskPath()
	if err != nil {
		log.Println("get disk path failed", err)
		return nil, err
	}
	for _, path := range paths {
		usageStat, err := disk.UsageWithContext(context.Background(), path)
		if err != nil {
			log.Println("get disk usage error:", err)
			return nil, err
		}
		return usageStat, nil
	}
	return nil, err
}

// GetPartition Disk分区信息
func (d *DiskInfo) GetPartition() ([]disk.PartitionStat, error) {
	partitionStats, err := disk.PartitionsWithContext(context.Background(), false)
	if err != nil {
		log.Println("get disk partitions error:", err)
		return nil, err
	}
	return partitionStats, nil
}

// GetIOCounters Disk IO信息
func (d *DiskInfo) GetIOCounters() (map[string]disk.IOCountersStat, error) {
	ioCountersStats, err := disk.IOCountersWithContext(context.Background())
	if err != nil {
		log.Println("get disk io counters error:", err)
		return nil, err
	}
	return ioCountersStats, nil
}

func (d *DiskInfo) ToString() (string, error) {
	// 定义返回string切片
	var v []string

	// Disk使用信息加入v中
	usageStat, err := d.GetUsage()
	if err != nil {
		log.Println("get disk usage error:", err)
		return "", err
	}
	v = append(v, usageStat.String())

	// Disk分区信息加入v中
	partitionStats, err := d.GetPartition()
	if err != nil {
		log.Println("get disk partitions error:", err)
		return "", err
	}
	for _, p := range partitionStats {
		v = append(v, p.String())
	}

	// Disk IO信息加入v
	ioCountersStat, err := d.GetIOCounters()
	if err != nil {
		log.Println("get disk io counters error:", err)
		return "", err
	}
	for _, ioCounters := range ioCountersStat {
		v = append(v, ioCounters.String())
	}

	disks := strings.Join(v, ",")
	marshal, err := json.Marshal(disks)

	return string(marshal), err
}

// NewDiskInfo 提供对外接口
func NewDiskInfo() *DiskInfo {
	return &DiskInfo{}
}
