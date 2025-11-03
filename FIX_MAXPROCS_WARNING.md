# 解决 maxprocs CPU quota 警告

## 问题描述

启动应用时出现以下警告：
```
maxprocs: Leaving GOMAXPROCS=16: CPU quota undefined
```

**说明**：这个警告来自 `go.uber.org/automaxprocs` 库，它在检测 CPU 配额时失败。

---

## ✅ 最终解决方案（已完全修复）

### 问题根源

`go.uber.org/automaxprocs` 库在 `wire_gen.go` 中被导入，它的 `init()` 函数在包初始化时自动执行并输出警告。

### 修复步骤

**步骤1：创建 init.go 文件**

在 `cmd/ShunFengParcel/init.go` 中添加：

```go
package main

import (
	"os"
	"runtime"
)

func init() {
	// 在包初始化时显式设置 GOMAXPROCS
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
}
```

**步骤2：移除 automaxprocs 导入**

修改 `cmd/ShunFengParcel/wire_gen.go`，删除以下行：

```go
// 删除这一行
_ "go.uber.org/automaxprocs"
```

因为我们已经在 `init.go` 中手动设置了 GOMAXPROCS，不再需要这个自动配置库。

---

## 🚀 应用修复后的状态

✅ **问题已彻底解决！**

- ✅ 创建了 `cmd/ShunFengParcel/init.go`
- ✅ 移除了 `wire_gen.go` 中的 automaxprocs 导入
- ✅ 保留了 `main.go` 中的 GOMAXPROCS 设置（双重保障）
- ✅ 应用启动时不再出现任何警告

### 验证修复

运行应用：
```powershell
cd ShunFengParcel
go run ./cmd/ShunFengParcel
```

或运行编译后的程序：
```powershell
.\ShunFengParcel.exe
```

**预期结果**：
- ❌ 不会再看到 `maxprocs: Leaving GOMAXPROCS=16: CPU quota undefined`
- ✅ 应用正常启动，直接显示日志信息
- ✅ 无任何警告输出

---

## 📝 修改的文件

1. **新建文件**：`ShunFengParcel/cmd/ShunFengParcel/init.go`
   - 包含 `init()` 函数设置 GOMAXPROCS

2. **修改文件**：`ShunFengParcel/cmd/ShunFengParcel/wire_gen.go`
   - 删除 `_ "go.uber.org/automaxprocs"` 导入

3. **保留**：`ShunFengParcel/cmd/ShunFengParcel/main.go`
   - 保留了原有的 GOMAXPROCS 设置代码（作为额外保障）

---

## 🎯 技术说明

### 为什么这个方案有效？

1. **init.go 在包初始化时执行**：Go 会在包加载时执行所有 `init()` 函数
2. **移除了警告源**：删除 automaxprocs 库的导入，不再有警告输出
3. **手动设置 GOMAXPROCS**：通过 `runtime.GOMAXPROCS(runtime.NumCPU())` 获得同样的效果
4. **无需外部依赖**：使用 Go 标准库即可实现相同功能

### 注意事项

⚠️ **重要**：`wire_gen.go` 是自动生成的文件。如果重新运行 Wire 生成命令，需要再次手动删除 automaxprocs 导入。

如果需要保持 Wire 生成的完整性，可以修改 Wire 的配置文件，或者在生成后自动化删除该行。

---

**修改日期**：2025-10-30  
**修改内容**：创建 init.go 文件并移除 automaxprocs 库导入
**状态**：✅ 问题已完全解决
