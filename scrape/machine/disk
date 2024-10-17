package disk

import (
	"github.com/shirou/gopsutil/v4/disk"
	"log"
)

// 定义磁盘类
type DiskInfo struct {
	Usages     []disk.UsageStat
	Partitions []disk.PartitionStat
}

func NewDiskInfo() *DiskInfo {
	return &DiskInfo{
		Partitions: make([]disk.PartitionStat, 0),
		Usages:     make([]disk.UsageStat, 0),
	}
}

// 获取内存总量
func (diskInfo *DiskInfo) GetToatl() uint64 {
	//toatl := diskInfo.Usages[0]
	return diskInfo.Usages[0].Total / 1024 / 1024 / 1024
}

// 获取使用的内存
func (diskInfo *DiskInfo) GetUsed() uint64 {
	//used := diskInfo.Usages[0]
	return diskInfo.Usages[0].Used / 1024 / 1024 / 1024
}

// 获取空闲内存
func (diskInfo *DiskInfo) GetFree() uint64 {
	//free := diskInfo.Usages[0]
	return diskInfo.Usages[0].Free / 1024 / 1024 / 1024
}

// 获取内存使用率
func (diskInfo *DiskInfo) GetUsedPercent() float64 {
	return diskInfo.Usages[0].UsedPercent
}

// 获取文件系统使用情况
func (diskInfo *DiskInfo) GetUsage(path string) []disk.UsageStat {
	usage, _ := disk.Usage(path)
	diskInfo.Usages = append(diskInfo.Usages, *usage)
	return diskInfo.Usages
}

// 返回所有磁盘信息
func (diskInfo *DiskInfo) GetPartitions(path string) []disk.PartitionStat {

	stats, err := disk.Partitions(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, stat := range stats {
		diskInfo.Partitions = append(diskInfo.Partitions, stat)
	}
	return diskInfo.Partitions
}
