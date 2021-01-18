package gfbean

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/google/uuid"
	"go.uber.org/dig"
)

// bean管理容器
var beanContainer *dig.Container

// bean变量
var beanVar cusBean

// 自定义bean
type cusBean struct {
	Id   string
	Time time.Time
}

// 自定义bean方法
func (b *cusBean) getBeanInfo() {
	g.Dump(b)
}

// 注入bean
func autowired(b cusBean) {
	beanVar = b
}

// 测试bean注入
func TestDiBean(t *testing.T) {
	// 实例化bean容器
	beanContainer = dig.New()
	// 注入自定义bean
	err := beanContainer.Provide(func() cusBean {
		return cusBean{
			Id:   uuid.New().String(),
			Time: time.Now(),
		}
	})
	if err != nil {
		fmt.Println(err)
	}
	// 使用bean
	err = beanContainer.Invoke(autowired)
	if err != nil {
		fmt.Println(err)
	}

	beanVar.getBeanInfo()
}
