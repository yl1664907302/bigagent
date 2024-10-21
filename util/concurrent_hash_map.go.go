package util

import (
	farmhash "github.com/leemcloughlin/gofarmhash"
	"sync"
)

type ConcurrentHashMap struct {
	//n个map
	mps []map[string]any
	//n个锁
	locks []sync.RWMutex
	//map数量
	mnum int
	//每次执行farmhash传统一的seed，使用hash算法根据此值对同一个key计算出不同的唯一值，读取不同的小map
	seed uint32
}

// 构造器
func CreateChashMap(num, cap int) *ConcurrentHashMap {
	mps := make([]map[string]any, num)
	locks := make([]sync.RWMutex, num)
	for i := 0; i < num; i++ {
		//为每个map平分容量
		mps[i] = make(map[string]any, cap/num)
	}
	return &ConcurrentHashMap{
		mps:   mps,
		locks: locks,
		mnum:  num,
		seed:  0,
	}
}

// 判断key对应的是哪个map，如果mnum是5，会在取模操作中选择五分之一作为选举的地址返回，这样做是为了并发读取地址去重
func (hm *ConcurrentHashMap) getSegIndex(key string) int {
	hash := int(farmhash.Hash32WithSeed([]byte(key), hm.seed))
	return hash % hm.mnum
}

// 读取操作
func (hm *ConcurrentHashMap) Get(key string) (any, bool) {
	//获取到本地使用的地址
	index := hm.getSegIndex(key)
	// 使用map锁,对改地址的map并发安全
	hm.locks[index].RLock()
	//最后解锁
	defer hm.locks[index].RUnlock()
	//主要读取操作
	value, exists := hm.mps[index][key]
	return value, exists
}

// 写入操作
func (hm *ConcurrentHashMap) Set(key string, value any) {
	index := hm.getSegIndex(key)
	hm.locks[index].Lock()
	defer hm.locks[index].Unlock()
	hm.mps[index][key] = value
}

// 迭代器返回值k-v载体
type MapEntry struct {
	Key   string
	Value any
}

// 迭代器模式
type MapIterator interface {
	Next() *MapEntry
}

// 迭代器主体
type ConcurrentHashMapIterator struct {
	cm    *ConcurrentHashMap
	keys  []string
	index int
}

// 创建迭代器
func (hm *ConcurrentHashMap) CreateIterator() *ConcurrentHashMapIterator {
	keys := make([]string, 0, len(hm.mps))
	//获取小map的row于col
	for _, mp := range hm.mps {
		//获取小map的所有key值,存入
		for key, _ := range mp {
			keys = append(keys, key)
		}
	}
	//返回迭代并且能够，配置好所有的key值
	return &ConcurrentHashMapIterator{
		cm:    hm,
		keys:  keys,
		index: 0,
	}
}

// 编写迭代方法:创建一个迭代器,将大map中的所有的小map的key值统一存入创建的二维字符串数组中,二维字符串数组的一维数被key的数量确定了,再使用大map的get方法遍历的将每个
// key的value,递归执行
func (iter *ConcurrentHashMapIterator) Next() *MapEntry {
	//递归结束判断
	if iter.index >= len(iter.keys) {
		return nil
	}

	//获取小map的字符串数组
	key := iter.keys[iter.index]

	//如果本行不存在就递归下一行继续寻找
	if len(key) == 0 {
		iter.index += 1
		return iter.Next()
	}

	//获取值
	value, _ := iter.cm.Get(key)
	iter.index++

	return &MapEntry{key, value}
}
