package gfbean

import (
	"go.uber.org/dig"
)

// 全局bean上下文
var ApplicationContext BeanFactory

// BeanFactory 定义bean容器工厂
type BeanFactory struct {
	Bean *dig.Container
}

// Init 实例化bean容器
func Init() {
	container := dig.New()
	ApplicationContext = BeanFactory{
		Bean: container,
	}
}

// NewBeanFactory 实例化bean容器
func NewBeanFactory() BeanFactory {
	return BeanFactory{
		Bean: dig.New(),
	}
}

// 设置bean
func (a *BeanFactory) setBean(function interface{}) error {
	err := a.Bean.Provide(function)
	return err
}

// 获取bean
func (a *BeanFactory) getBean(function interface{}) error {
	err := a.Bean.Invoke(function)
	return err
}
