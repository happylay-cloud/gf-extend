# gf-extend
GoFrame框架扩展工具包

### **安装**

> go get -u -v github.com/happylay-cloud/gf-extend

或

> require github.com/happylay-cloud/gf-extend latest

### 自动实例化gf-casbin执行器
```
e, err := gfadapter.NewEnforcer()
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


