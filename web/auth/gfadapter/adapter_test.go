package gdbadapter

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"log"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
)

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	myRes := e.GetPolicy()
	log.Print("策略：", myRes)

	if !util.Array2DEquals(res, myRes) {
		t.Error("策略：", myRes, "，应该是", res)
	}
}

func initPolicy(t *testing.T, a *Adapter) {

	// 因为数据库一开始是空的，
	// 因此，我们需要首先从文件适配器（.csv）加载策略。
	e, err := casbin.NewEnforcer("./config/rbac_model.conf", "./config/rbac_policy.csv")
	if err != nil {
		panic(err)
	}

	// 这是一个保存当前策略到数据库的技巧。
	// 我们不能调用e.SavePolicy（），因为强制执行器中的适配器仍然是文件适配器。
	// 当前策略是指Casbin执行器中的策略（也称为内存中的策略）。
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}

	// 清除当前策略。
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// 从数据库加载策略。
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"}})
}

func testSaveLoad(t *testing.T, a *Adapter) {
	// 初始化数据库中的某些策略。
	initPolicy(t, a)
	// 注意：如果您已经有一个正在工作的数据库，其中包含策略，您不需要查看上面的代码

	// 现在数据库有了策略，所以我们可以提供一个正常的用例。
	// 创建适配器和执行器。
	// NewEnforcer()将自动加载策略。
	e, _ := casbin.NewEnforcer("./config/rbac_model.conf", a)
	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"}})
}

func initAdapter(t *testing.T, driverName string, dataSourceName string) *Adapter {

	// 创建适配器
	a, err := NewAdapter(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}

	// 初始化数据库中的某些策略。
	initPolicy(t, a)
	// 现在数据库有了策略，所以我们可以提供一个正常的用例。
	// 注意：如果您已经有一个包含策略的工作数据库，您不需要查看上面的代码。

	return a
}

func initAdapterFormOptions(t *testing.T, adapter *Adapter) *Adapter {
	// 创建一个适配器
	a, _ := NewAdapterFromOptions(adapter)
	// 在数据库中初始化一些策略。
	initPolicy(t, a)
	// 现在数据库有了策略，所以我们可以提供一个正常的用例。
	// 注意：如果您已经有一个包含策略的工作数据库，您不需要查看上面的代码。

	return a
}

func testAutoSave(t *testing.T, a *Adapter) {

	// NewEnforcer()会自动加载策略。
	e, _ := casbin.NewEnforcer("./config/rbac_model.conf", a)
	// 默认情况下启用自动保存。
	// 现在我们禁用它。
	e.EnableAutoSave(false)

	// 由于禁用了自动保存，策略更改只影响Casbin执行器中的策略，它不影响存储中的策略。
	_, _ = e.AddPolicy("alice", "data1", "write")
	// 从存储重新加载策略以查看效果。

	_ = e.LoadPolicy()
	// 这仍然是原来的政策。
	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"}})

	// 现在我们启用自动保存。
	e.EnableAutoSave(true)

	// 由于启用了自动保存，策略更改不仅影响Casbin执行器中的策略，但也会影响存储中的策略。
	_, _ = e.AddPolicy("alice", "data1", "write")

	// 从存储中重新加载策略以查看效果。
	_ = e.LoadPolicy()

	// 策略有了新规则：{"alice", "data1", "write"}.
	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"},
		{"alice", "data1", "write"}})

	// 删除添加的规则。
	_, _ = e.RemovePolicy("alice", "data1", "write")
	_ = e.LoadPolicy()

	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"},
		{"data2_admin", "data2", "read"},
		{"data2_admin", "data2", "write"}})

	// 通过过滤器删除与"data2_admin"相关的策略规则。
	// 删除两条规则：{"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}
	_, _ = e.RemoveFilteredPolicy(0, "data2_admin")
	_ = e.LoadPolicy()
	testGetPolicy(t, e, [][]string{
		{"alice", "data1", "read"},
		{"bob", "data2", "write"}})
}

func TestMysqlAdapters(t *testing.T) {

	// 测试mysql
	a := initAdapter(t, "mysql", "root:root@tcp(127.0.0.1:3306)/casbin")
	testAutoSave(t, a)
	testSaveLoad(t, a)

	a = initAdapterFormOptions(t, &Adapter{
		driverName:     "mysql",
		dataSourceName: "root:root@tcp(127.0.0.1:3306)/casbin",
	})

	testAutoSave(t, a)
	testSaveLoad(t, a)

}

func TestPgsqlAdapters(t *testing.T) {

	// 测试pgsql
	a := initAdapter(t, "pgsql", "user=postgres host=127.0.0.1 port=5432 sslmode=disable dbname=casbin")
	testAutoSave(t, a)
	testSaveLoad(t, a)

	a = initAdapterFormOptions(t, &Adapter{
		driverName:     "pgsql",
		dataSourceName: "user=postgres host=127.0.0.1 port=5432 sslmode=disable dbname=casbin",
	})

	testAutoSave(t, a)
	testSaveLoad(t, a)
}
