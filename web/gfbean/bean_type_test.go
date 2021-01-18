package gfbean

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

type b1 struct {
	id string
}

func (b *b1) getId() {
	fmt.Println(b.id)
}

func TestBeanDig(t *testing.T) {
	// 1.初始化全局bean工厂
	Init()

	// 2.设置bean
	if err := ApplicationContext.setBean(func() b1 {
		return b1{
			id: uuid.New().String(),
		}
	}); err != nil {
		fmt.Println(err)
	}

	var b b1

	// 3.获取bean
	if err := ApplicationContext.getBean(func(b1 b1) {
		b = b1
	}); err != nil {
		fmt.Println(err)
	}

	b.getId()

}
