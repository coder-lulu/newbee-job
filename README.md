# 新蜂资产管理平台 — 定时任务服务

基于 Asynq 的定时任务 RPC 服务，提供任务配置、调度和执行相关能力。Redis 承载任务队列，Ent 管理数据库实体，可供平台业务按 RPC 契约接入。

仓库：[coder-lulu/newbee-job](https://github.com/coder-lulu/newbee-job) · [平台工作区](https://github.com/coder-lulu/newbee)

## 获取代码

推荐通过完整工作区开发，保留兄弟模块目录及本地 `replace` 依赖。以下命令使用 Bash；Go 工作区要求 Go 1.25.1 或更高版本。

```bash
git clone --recurse-submodules https://github.com/coder-lulu/newbee.git
cd newbee/job
```

已有工作区执行 `git submodule update --init --recursive`。单独克隆模块时，需要自行补齐 `go.mod` 中的本地依赖路径。

## 目录导航

| 路径 | 用途 |
| --- | --- |
| `job.go`、`job.proto` | 服务入口和 RPC 协议 |
| `jobclient/`、`types/` | RPC 客户端和生成类型 |
| `internal/mqs/amq/task/` | 队列、定时及动态周期任务 |
| `internal/logic/` | 任务业务逻辑 |
| `ent/` | 数据模型 |
| `etc/job.yaml.example` | 配置模板 |

## 配置与运行

首次复制 `etc/job.yaml.example` 到 `etc/job.yaml`，设置 `DatabaseConf`、`RedisConf`、`AsynqConf` 和 `TaskConf`，并按环境调整监听地址。示例 RPC 端口为 `9105`。任务开关应与所需消费者和调度任务一致。

入口使用 `conf.UseEnv()` 展开配置中的 `${变量名}`。启动前注入模板引用的环境变量；真实配置与凭据不要提交。首次部署先创建空数据库，配置独立 Redis 逻辑库，并保持 `AsynqConf.Enable`、`TaskConf.EnableDPTask` 与 `EnableScheduledTask` 为 `false`（模板和代码默认均关闭）。数据库和 Redis 就绪后运行：

```bash
go run . -f etc/job.yaml
```

RPC 启动后，在另一个终端从 `job` 目录初始化空库（需要 `grpcurl`）：

```bash
grpcurl -plaintext -max-time 300 -import-path . -proto job.proto \
  -d '{}' 127.0.0.1:9105 job.Job/initDatabase
```

初始化仅创建或补齐表结构，不插入 `hello_world` 或其他示例任务；重复调用成功返回，不重置已有任务。RPC 初始化入口仅向受控内网开放。已有数据库升级前仍需备份并审核 Schema 差异。

随后创建所需任务、核对任务状态和 Pattern 对应的处理器，先启用 `AsynqConf.Enable`，再按需启用动态周期任务 `EnableDPTask` 或预定义调度 `EnableScheduledTask` 并重启。`AsynqConf.Enable: false` 时仅启动 RPC，不启动 MQ 消费者及调度器，可用于首次建表。启用 Asynq 后，即使两个调度开关关闭，MQ 消费者仍会消费配置队列中的既有任务；空库验证应使用独立 Redis 逻辑库。

## 构建与验证

```bash
go build .
go test ./...
go vet ./...
```

生成 RPC 和 Ent 的命令见 [Makefile](Makefile)。平台其他服务是否已切换到本仓库客户端，应以其 `go.mod` 和调用代码为准；部署本模块不会自动替换对上游 Simple Admin Job 的依赖。

## 上游文档

- [Simple Admin 定时任务说明](https://doc.ryansu.tech/zh/guide/official-comp/cron.html)

## 许可证与来源

本仓库采用 [Apache-2.0](LICENSE)。源自 [simple-admin-job](https://github.com/suyuan32/simple-admin-job)，保留 Ryan SU Authors 版权。第三方依赖遵循各自许可证，保留原有版权与许可声明。
