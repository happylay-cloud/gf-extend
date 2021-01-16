package gfboot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gres"
	"io/ioutil"
)

// SingleFileMemoryToLocal 将打包到二进制中的单文件，解压到本地
//
// 示例：
// start.SingleFileMemoryToLocal("./db", "sqlite3.db", "db/sqlite3.db")
//
// @Param folderPath：文件所在文件夹（可以是相对路径以./开头）
//
// @Param fileName：  文件名称
//
// @Param memoryPath：内存文件路径（不能以./开头）
func SingleFileMemoryToLocal(folderPath, fileName, memoryPath string) {

	// 判断文件夹是否为空
	if empty := gfile.IsEmpty(folderPath); empty {
		// 文件夹为空则创建文件夹
		if err := gfile.Mkdir(folderPath); err == nil {
			// 从内存中获取资源文件
			file := gres.Get(memoryPath)
			if file == nil {
				g.Log().Error("\n" +
					"无法获取指定文件：" + fileName + "，解压失败。\n" +
					"1.请检查可执行文件是否已打包资源文件\n" +
					"2.内存文件路径memoryPath：" + memoryPath + "+是否正确。")
				gres.Dump()
				return
			}
			// 将资源文件写入本地
			if err := ioutil.WriteFile(folderPath+"/"+fileName, file.Content(), 0666); err != nil {
				// 记录错误日志
				g.Log().Error(err)
				return
			}
			// 打印文件信息
			g.Dump(file)
		} else {
			// 记录错误日志
			g.Log().Error(err)
		}

	}

}
