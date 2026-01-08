# Exchange Backend（Go / go-zero 微服务）

一个基于 **go-zero**（`rest` + `zrpc`）构建的交易所后端练习项目，按领域拆分为用户中心、行情中心、交易撮合/订单、任务调度等多个服务，并通过 **Kafka** 串联实时数据流，`market-api` 通过 **Socket.IO(WebSocket)** 对外推送行情/盘口。

## 目录结构

- `ucenter`：用户中心 **RPC**（注册/登录/资产/会员/提币等），并消费部分业务 Kafka 主题完成异步入账/状态更新。
- `ucenter-api`：用户中心 **REST API** 网关（对外 HTTP），内部调用 `ucenter` / `market` RPC。
- `market`：行情中心 **RPC**（币种/交易对/汇率/K 线历史等），落库 MySQL + Mongo。
- `market-api`：行情 **REST API + WebSocket**（`/socket.io`），从 Kafka 消费 K 线/盘口等数据并推送到前端。
- `exchange`：交易/撮合 **RPC**（订单相关），对接 `ucenter` / `market` RPC，并通过 Kafka 进行异步流程与推送。
- `exchange-api`：交易 **REST API** 网关（对外 HTTP），内部调用 `exchange` RPC。
- `jobcenter`：任务中心（定时拉取 OKX 行情与汇率、写入 Mongo/Redis、推送 Kafka；代码里也包含 BTC 扫描逻辑入口）。
- `grpc-common`：goctl 生成的 **RPC client** 封装（用于各服务间调用）。
- `mscoin-common`：公共库（`msdb`、工具函数、代码生成器等）。

## 技术栈与依赖

- 语言/框架：Go、go-zero（rest/zrpc）、gRPC
- 存储：MySQL（GORM）、MongoDB
- 缓存：Redis（go-zero cache）
- 消息队列：Kafka（`segmentio/kafka-go`）
- 实时推送：Socket.IO（`github.com/googollee/go-socket.io`）
- 外部数据源：OKX Open API（`jobcenter` 拉取 K 线与汇率）

## 快速开始（本地开发）

### 1) 建议使用 Go Workspace 管理多模块

本仓库包含多个 `go.mod`，并且模块间使用本地导入路径（例如 `grpc-common/...`）。建议在仓库根目录创建 `go.work`：

```bash
go work init
go work use ./exchange ./exchange-api ./market ./market-api ./ucenter ./ucenter-api ./jobcenter ./grpc-common ./mscoin-common
```

### 2) 配置文件

各服务启动时会默认读取 `etc/*.yaml`（见各 `main.go` 的 `-f` 参数）。仓库当前未包含 `etc` 配置文件，你需要自行创建，例如：

- `ucenter/etc/config.yaml`
- `ucenter-api/etc/config.yaml`
- `market/etc/conf.yaml`
- `market-api/etc/conf.yaml`
- `exchange/etc/conf.yaml`
- `exchange-api/etc/conf.yaml`
- `jobcenter/etc/conf.yaml`

配置字段可对照各服务的 `internal/config/config.go`（包含 MySQL/Mongo/Redis/Kafka/OKX/JWT/RPC Client 等配置结构）。

### 3) 启动依赖服务

根据你的配置，先准备并启动：

- MySQL
- Redis
- MongoDB
- Kafka

`jobcenter` 还需要 OKX 的 `ApiKey/SecretKey/Pass` 等配置（可选代理 `Proxy`）。

### 4) 启动各服务

在对应服务目录下运行（示例）：

```bash
cd ucenter && go run . -f etc/config.yaml
cd market && go run . -f etc/conf.yaml
cd exchange && go run . -f etc/conf.yaml
cd jobcenter && go run . -f etc/conf.yaml

cd ucenter-api && go run . -f etc/config.yaml
cd market-api && go run . -f etc/conf.yaml
cd exchange-api && go run . -f etc/conf.yaml
```

仓库根目录提供了 `services.sh`（bash）用于一键启动/停止部分服务，你也可以按需补齐（脚本当前未包含 `exchange` / `exchange-api`）。

## WebSocket（market-api）

- Socket.IO 路径：`/socket.io`
- 房间：固定为 `market`（见 `market-api/internal/ws/ws.go`）
- 数据来源：`market-api` 从 Kafka 消费 K 线与盘口等主题，然后广播推送到前端

## Kafka 主题（部分）

代码中可见的主题包含（仅做导航，具体以业务为准）：

- K 线：`kline_1m`（以及其它周期相关逻辑）
- 盘口：`exchange_order_trade_plate`
- 订单/异步流程：`add-exchange-order`、`exchange_order_complete_update_success`、`withdraw`、`BTC_TRANSACTION`

## 代码生成说明（gRPC types）

`grpc-common/*client` 依赖 `grpc-common/**/types/*`（gRPC stub）包。目前仓库未包含这些 `types` 目录时会导致编译失败；你需要根据 `.proto`/goctl 生成或补齐对应代码后再编译运行。

## 免责声明

本项目更偏向学习/演示用途，并不包含生产级的安全、风控、资金安全与审计能力。请勿直接用于真实资金场景。

