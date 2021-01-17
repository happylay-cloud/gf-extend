# gf-extend
GoFrame框架扩展工具包

### **安装**

> go get -u -v -d github.com/happylay-cloud/gf-extend

或

> require github.com/happylay-cloud/gf-extend latest

### 自动实例化gf-casbin执行器
```
// 自动注册
e, err := gfadapter.NewEnforcer()

// 手动注册
e, err := gfadapter.NewEnforcer(g.DB())

// 手动注册-自定义数据源
e, err := gfadapter.NewEnforcer(g.DB("mysql"))

// 手动注册-自定义数据源
e, err := gfadapter.NewEnforcer(g.DB("casbin"))
```

### 解压二进制中单文件到本地
```
gfboot.SingleFileMemoryToLocal("./db", "sqlite3.db", "db/sqlite3.db")
```

### 获取web响应对象
```
gfres.Ok()
```

### 获取page分页对象
```
dto.NewPage()
```

### 更多...


