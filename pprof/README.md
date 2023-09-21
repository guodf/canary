# pprof使用

## 启动pprof性能工具

```
go pprof.Run(6060)
```

## 使用pprof生成性能图

**需要安装：https://graphviz.org/download/**

`收集 CPU 30s 性能监控`

```
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?second=30
```