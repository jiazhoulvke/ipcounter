package ipcounter

import (
	"sync"
	"time"
)

//IPCounter ip访问数统计器
type IPCounter struct {
	countMap map[int64]map[string]int
	//timeRange 时间区域，指最大统计的秒数。
	timeRange int
	mux       sync.RWMutex
}

//New 新的IPCounter
func New(maxLength int) *IPCounter {
	return &IPCounter{
		countMap:  make(map[int64]map[string]int),
		timeRange: maxLength,
	}
}

//Add 添加IP
func (l *IPCounter) Add(ip string) {
	now := time.Now().Unix()
	l.mux.Lock()
	if _, ok := l.countMap[now]; !ok {
		l.countMap[now] = make(map[string]int)
		go func() {
			time.Sleep(time.Duration(l.timeRange+1) * time.Second) //定时删除过时的数据
			l.mux.Lock()
			delete(l.countMap, now)
			l.mux.Unlock()
		}()
	}
	l.countMap[now][ip]++
	l.mux.Unlock()
}

//Count 统计IP总数
func (l *IPCounter) Count(ip string, timeRange int) int {
	var count int
	now := time.Now().Unix()
	l.mux.Lock()
	if timeRange <= 0 {
		timeRange = l.timeRange
	}
	for t := now; t > now-int64(timeRange); t-- {
		if m, exist := l.countMap[t]; exist {
			if c, ok := m[ip]; ok {
				count += c
			}
		}
	}
	l.mux.Unlock()
	return count
}
