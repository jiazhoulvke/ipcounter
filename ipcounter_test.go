package ipcounter

import (
	"math/rand"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIPCounter(t *testing.T) {
	Convey("测试IP计数器", t, func() {
		var (
			maxLength = 20
			wg        sync.WaitGroup
		)
		ipCounter := New(5)
		ipList := []string{"127.0.0.1", "192.168.1.1", "8.8.8.8"}
		for i := 0; i < maxLength; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 10; j++ {
					ipCounter.Add(ipList[rand.Intn(len(ipList))])
				}
				wg.Done()
			}()
		}
		wg.Wait()
		var count int
		for _, ip := range ipList {
			c := ipCounter.Count(ip, 0)
			count += c
			t.Logf("ip: %s count: %d\n", ip, c)
		}
		So(count, ShouldEqual, 10*maxLength)

	})
}
