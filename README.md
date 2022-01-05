# Gf-Extend 基于v1.15.1版本

**GoFrame框架扩展工具包**

![gf-plus](https://img.shields.io/badge/gf-plus-ea7b99)
![GitHub last commit](https://img.shields.io/github/last-commit/happylay-cloud/gf-extend?style=flat-square)
[![Go Doc](https://godoc.org/happylay-cloud/gf-extend?status.svg)](https://pkg.go.dev/github.com/happylay-cloud/gf-extend)
[![Go Report](https://goreportcard.com/badge/happylay-cloud/gf-extend?v=1)](https://goreportcard.com/report/happylay-cloud/gf-extend)
[![Production Ready](https://img.shields.io/badge/production-ready-blue.svg)](https://github.com/happylay-cloud/gf-extend)
[![License](https://img.shields.io/github/license/happylay-cloud/gf-extend.svg?style=flat)](https://github.com/happylay-cloud/gf-extend)

#### **安装**

> go get -u -v -d github.com/happylay-cloud/gf-extend

**或**

> require github.com/happylay-cloud/gf-extend latest

#### 实例化gf-casbin实例bean（推荐使用）

* **自动注册（无需关心数据源种类）**

```
e, err := gfadapter.NewEnforcerBean()
```

* **手动注册**

```markdown
e, err := gfadapter.NewEnforcerBean(g.DB())
e, err := gfadapter.NewEnforcerBean(g.DB("sqlite"))
e, err := gfadapter.NewEnforcerBean(g.DB("mysql"))
e, err := gfadapter.NewEnforcerBean(g.DB("pgsql"))
```

#### 实例化gf-casbin执行器

* **自动注册**

```
e, err := gfadapter.NewEnforcer()
```

* **手动注册**

```
e, err := gfadapter.NewEnforcer(g.DB())
e, err := gfadapter.NewEnforcer(g.DB("mysql"))
e, err := gfadapter.NewEnforcer(g.DB("sqlite"))
e, err := gfadapter.NewEnforcer(g.DB("pgsql"))
```

#### 解压二进制中单文件到本地

```
gfboot.SingleFileMemoryToLocal("./db", "sqlite3.db", "db/sqlite3.db")
```

#### 获取web响应对象

```
hres.Ok()
```

#### 获取page分页对象

```
hdto.NewPage()
```

#### **更多...**
