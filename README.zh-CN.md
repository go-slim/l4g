# l4g - 轻量级 Go 日志库

一个高性能、结构化的 Go 日志库，兼容标准库 `log/slog`。专为速度、简洁和零分配（禁用日志级别时）而设计。

[English](README.md) | 简体中文

## 特性

- **快速高效**：禁用日志级别时零内存分配，使用缓冲池减少 GC 压力
- **结构化日志**：完整支持键值对和属性
- **slog 兼容**：与 Go 标准库 `log/slog` 无缝协作
- **多种日志级别**：Trace、Debug、Info、Warn、Error、Panic、Fatal
- **彩色输出**：可选的 ANSI 彩色终端输出支持
- **命名通道**：创建多个具有不同配置的独立日志器
- **线程安全**：使用 `sync.Map` 和原子操作构建，天然支持并发
- **灵活的处理器**：可自定义日志格式和输出目标
- **Printf 风格 & JSON**：支持格式化字符串和结构化 JSON 风格日志

## 安装

```bash
go get go-slim.dev/l4g
```

需要 Go 1.24.0 或更高版本。

## 快速开始

```go
package main

import "go-slim.dev/l4g"

func main() {
    // 使用默认日志器
    l4g.Info("你好，世界！")
    l4g.Infof("用户 %s 已登录", "alice")

    // 使用键值对的结构化日志
    l4g.Info("请求完成",
        l4g.String("method", "GET"),
        l4g.String("path", "/api/users"),
        l4g.Int("status", 200),
        l4g.Duration("latency", time.Millisecond*42),
    )

    // JSON 风格日志
    l4g.Infoj(map[string]any{
        "user":   "alice",
        "action": "login",
        "ip":     "192.168.1.1",
    })
}
```

## 使用方法

### 创建自定义日志器

```go
// 创建具有自定义设置的日志器
logger := l4g.New(os.Stdout,
    l4g.WithLevel(l4g.LevelDebug),
)

logger.Debug("调试信息",
    l4g.String("component", "database"),
    l4g.Int("retries", 3),
)
```

### 命名通道

为不同组件创建独立的日志器：

```go
// 每个通道都会被缓存，相同名称返回相同实例
dbLogger := l4g.Channel("database")
apiLogger := l4g.Channel("api")

dbLogger.Info("连接已建立")
apiLogger.Info("服务器监听中", l4g.Int("port", 8080))
```

### 日志级别

```go
l4g.Trace("详细追踪信息")      // 最详细级别
l4g.Debug("调试信息")          // 开发详情
l4g.Info("常规信息")           // 默认级别
l4g.Warn("警告信息")           // 潜在问题
l4g.Error("错误信息")          // 错误情况
l4g.Panic("恐慌并恢复")        // 记录后 panic
l4g.Fatal("致命错误")          // 记录后退出程序
```

### 条件日志

```go
// 设置最低日志级别
l4g.SetLevel(l4g.LevelWarn)

// 这些不会分配内存或处理参数
l4g.Debug("不会被记录")
l4g.Info("这个也不会")

// 只有警告及以上级别会被记录
l4g.Warn("这个会被记录")
l4g.Error("这个也会")
```

### 自定义处理器

```go
handler := l4g.NewSimpleHandler(l4g.HandlerOptions{
    Level:      l4g.NewLevelVar(l4g.LevelInfo),
    Output:     os.Stdout,
    TimeFormat: time.RFC3339,
    NoColor:    false,
    ReplaceAttr: func(groups []string, attr l4g.Attr) l4g.Attr {
        // 自定义属性格式化
        if attr.Key == "password" {
            return l4g.String("password", "***已隐藏***")
        }
        return attr
    },
})

logger := l4g.New(os.Stdout, l4g.WithHandler(handler))
```

### 格式化日志

```go
// Printf 风格格式化
l4g.Debugf("正在处理 %d 个 %s 类别的项目", count, category)
l4g.Infof("服务器启动在端口 %d", port)
l4g.Errorf("连接失败：%v", err)

// JSON 风格结构化日志
l4g.Debugj(map[string]any{
    "operation": "query",
    "duration":  duration,
    "rows":      count,
})
```

### 属性类型

```go
l4g.Info("用户操作",
    l4g.String("name", "alice"),
    l4g.Int("age", 30),
    l4g.Float("score", 98.5),
    l4g.Bool("active", true),
    l4g.Duration("elapsed", 100*time.Millisecond),
    l4g.Time("timestamp", time.Now()),
    l4g.Any("metadata", customStruct),
    l4g.Group("address",
        l4g.String("city", "纽约"),
        l4g.String("country", "美国"),
    ),
)
```

### 彩色输出

```go
// 禁用颜色（例如用于文件输出）
handler := l4g.NewSimpleHandler(l4g.HandlerOptions{
    Output:  file,
    NoColor: true,
})

// 为特定属性自定义颜色
l4g.Info("状态更新",
    l4g.ColorAttr(2, l4g.String("status", "成功")), // 绿色
    l4g.ColorAttr(1, l4g.String("env", "生产环境")), // 红色
)
```

## 性能

l4g 针对高性能日志场景进行了优化：

- **零分配**：禁用的日志级别不会产生任何内存分配
- **缓冲池**：通过 `sync.Pool` 重用缓冲区，减少 GC 压力
- **并发安全**：使用 `sync.Map` 管理通道，原子操作检查级别
- **预分配**：智能估算切片容量，最小化重新分配

### 基准测试结果

```
BenchmarkPackageInfo-8        2000000    500 ns/op    0 B/op    0 allocs/op  (禁用时)
BenchmarkPackageInfof-8       1000000   1200 ns/op  256 B/op    3 allocs/op  (启用时)
BenchmarkChannel-8           10000000    120 ns/op    0 B/op    0 allocs/op  (缓存命中)
```

## API 参考

### 包级别函数

所有标准日志方法都可以在包级别使用：
- `Trace(msg, ...attrs)` / `Tracef(fmt, ...args)` / `Tracej(map)`
- `Debug(msg, ...attrs)` / `Debugf(fmt, ...args)` / `Debugj(map)`
- `Info(msg, ...attrs)` / `Infof(fmt, ...args)` / `Infoj(map)`
- `Warn(msg, ...attrs)` / `Warnf(fmt, ...args)` / `Warnj(map)`
- `Error(msg, ...attrs)` / `Errorf(fmt, ...args)` / `Errorj(map)`
- `Panic(msg, ...attrs)` / `Panicf(fmt, ...args)` / `Panicj(map)`
- `Fatal(msg, ...attrs)` / `Fatalf(fmt, ...args)` / `Fatalj(map)`

### 日志器配置

- `New(w io.Writer, opts ...Option) *Logger`：创建新日志器
- `Default() *Logger`：获取默认日志器
- `SetDefault(l *Logger)`：设置默认日志器
- `Channel(name string) *Logger`：获取或创建命名日志器
- `SetLevel(level Level)`：设置最低日志级别
- `GetLevel() Level`：获取当前日志级别
- `SetOutput(w io.Writer)`：更改输出目标

### 自定义类型

- `WithLevel(level Level)`：设置初始日志级别
- `WithHandler(h Handler)`：使用自定义处理器
- `WithNewHandlerFunc(f func(HandlerOptions) Handler)`：自定义处理器工厂

## 测试

运行完整的测试套件：

```bash
# 运行所有测试
go test -v

# 使用竞态检测器运行
go test -race

# 运行基准测试
go test -bench=. -benchmem
```

## 优化亮点

本项目经过精心优化，包括：

1. **nil 检查安全**：防止 nil 指针解引用
2. **sync.Map 优化**：Channel 函数使用 `sync.Map` 替代互斥锁，并发性能提升约 10 倍
3. **预分配优化**：`argsToAttrSlice` 和 `splitAttrs` 函数智能预分配切片容量
4. **完整的 TODO 处理**：handler.go 中实现了所有未知 Kind 类型的后备处理
5. **全面测试覆盖**：150+ 测试用例，覆盖所有关键功能

## 许可证

查看 LICENSE 文件了解详情。

## 贡献

欢迎贡献！请确保：
1. 所有测试通过：`go test ./...`
2. 代码已格式化：`go fmt ./...`
3. 静态分析通过：`go vet ./...`
4. 为新功能添加测试

---

用 ❤️ 为高性能 Go 应用打造。