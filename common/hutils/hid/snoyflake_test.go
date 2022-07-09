package hid

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
)

func TestSonyFlake(t *testing.T) {
	parse, _ := time.Parse(time.RFC3339, "2021-01-01")

	settings := Settings{
		StartTime: parse,
	}

	sf := NewSonyFlake(settings)

	// 3395,6479,3034,2449,89
	// 339564902606177149
	// 339564939314791293
	// 339564980184089469
	// 339565231053800317
	// 339565309868049277
	// 339565349785502589
	go func() {
		for i := 0; i < 100000; i++ {
			id, err := sf.NextID()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(id)
			decompose := Decompose(id)
			g.Dump(decompose)
		}
	}()

	select {}
}

func TestIdBean(t *testing.T) {
	// 设置同步等待组
	wg := new(sync.WaitGroup)
	wg.Add(2)

	ids := gmap.New(true)

	go func() {
		// 等待组计数器减一
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("唯一id1：", id)
			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间1：", format)

			if ok := ids.Contains(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id1")

			fmt.Println("唯一id1：", id)

		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间2：", format)

			if ok := ids.Contains(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id2")

			fmt.Println("唯一id2：", id)

		}

	}()
	// 让主协程处于等待状态
	wg.Wait()
	g.Dump(ids)
}

type SafeMap struct {
	Data   map[uint64]string
	Lock   *sync.Mutex   // 通用锁
	LockRw *sync.RWMutex // 读写锁

}

func (d *SafeMap) Get(k uint64) (value string, ok bool) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	if value, ok := d.Data[k]; ok {
		return value, true
	} else {
		return value, false
	}
}

func (d *SafeMap) Set(k uint64, v string) {
	d.Lock.Lock()
	defer d.Lock.Unlock()
	d.Data[k] = v

}

func (d *SafeMap) GetV(k uint64) (value string, ok bool) {
	d.LockRw.RLock()
	defer d.LockRw.RUnlock()
	if value, ok := d.Data[k]; ok {
		return value, true
	} else {
		return value, false
	}
}

func (d *SafeMap) SetKV(k uint64, v string) {
	d.LockRw.Lock()
	defer d.LockRw.Unlock()
	d.Data[k] = v
}

func TestIdBeanSafeMap(t *testing.T) {
	// 设置同步等待组
	wg := new(sync.WaitGroup)
	wg.Add(5)

	ids := SafeMap{
		Data:   map[uint64]string{},
		Lock:   new(sync.Mutex),
		LockRw: new(sync.RWMutex),
	}

	go func() {
		defer wg.Done()
		for i := 0; i < 30000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("唯一id1：", id)
			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间1：", format)

			if _, ok := ids.Get(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id1："+format)

			fmt.Println("唯一id1：", id)

		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 30000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("唯一id2：", id)
			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间2：", format)

			if _, ok := ids.Get(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id2："+format)

			fmt.Println("唯一id2：", id)

		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 30000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("唯一id3：", id)
			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间3：", format)

			if _, ok := ids.Get(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id3："+format)

			fmt.Println("唯一id3：", id)

		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 30000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("唯一id4：", id)
			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间4：", format)

			if _, ok := ids.Get(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id4："+format)

			fmt.Println("唯一id4：", id)

		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 30000; i++ {
			id, err := IdBean.NextID()
			if err != nil {
				fmt.Println(err)
			}

			decompose := Decompose(id)
			g.Dump(decompose)

			format := time.Unix(int64(decompose["time"]*10/1000), 0).Format("2006-01-02 15:04:05")
			fmt.Println("生成时间5：", format)

			if _, ok := ids.Get(id); ok {
				g.Log().Line(false).Error("id重复xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")
			}
			ids.Set(id, "唯一id5："+format)

			fmt.Println("唯一id5：", id)

		}

	}()

	wg.Wait()
	g.Dump(ids)
}

func TestIdBeanTime(t *testing.T) {
	id1, _ := IdBean.NextID()
	fmt.Println(id1, GetIdGenTime(id1))

	id2, _ := IdBean.NextID()
	fmt.Println(id2, GetIdGenTime(id2))
}
