package disk

import (
	"github.com/shirou/gopsutil/v4/disk"
)

// DiskInfo 定义磁盘类
type DiskInfo struct {
	//Usages     *disk.UsageStat                `json:"usages"`
	Usages     []*disk.UsageStat              `json:"usages"`
	Partitions []disk.PartitionStat           `json:"partitions"`
	IOCounters map[string]disk.IOCountersStat `json:"iocounters"`
}

// 容纳系统磁盘路径
var fsPath []string

func init() {
	// 构造系统filesystem 挂载路径的字符串切片
	path, _ := disk.Partitions(true)
	// 循环加入fsPath字符串切片中
	for _, part := range path {
		fsPath = append(fsPath, part.Mountpoint)
		//fmt.Printf("type: %T\n", part.Mountpoint)
	}
	//fmt.Printf("fssystem path: %v len: %v\n", fsPath, len(fsPath))
}

// NewDiskInfo 类型
func NewDiskInfo() *DiskInfo {
	//初始化一个disk.UsageStat对象的切片，长度是上文中磁盘路径的个数
	usage := make([]*disk.UsageStat, len(fsPath))
	// 遍历fsPath 将切片中的路径传给disk.Usage进行获取数据
	for i, path := range fsPath {
		// 返回对应路径的disk.UsageStat对象指针
		stat, err := disk.Usage(path)
		if err != nil {
			continue
		}
		// 赋值给索引i
		usage[i] = stat
	}
	partitions, _ := disk.Partitions(true)
	ioCounters, _ := disk.IOCounters()
	return &DiskInfo{
		Usages:     usage,
		Partitions: partitions,
		IOCounters: ioCounters,
	}
}

// GetUsage 获取文件系统使用情况
func (diskInfo *DiskInfo) GetUsage(path string) *disk.UsageStat {
	// 遍历fsPath 内的路径字符串，如果传入的path有则返回对于的disk.UsageStat对象
	for _, usage := range diskInfo.Usages {
		if usage.Path == path {
			return usage
		}
	}
	return nil
}

// GetPartitions 返回所有磁盘信息
func (diskInfo *DiskInfo) GetPartitions() []disk.PartitionStat {
	return diskInfo.Partitions
}

// GetIOCounters 返回io计数器
func (diskInfo *DiskInfo) GetIOCounters() map[string]disk.IOCountersStat {
	return diskInfo.IOCounters
}
