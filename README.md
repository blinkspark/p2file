# p2file

p2file 是一个基于 libp2p 的 P2P 文件共享工具，允许用户在去中心化网络中共享和下载文件。

## 功能

- **P2P 文件共享**：使用 libp2p 协议实现去中心化文件共享
- **目录共享**：服务器模式下共享本地目录
- **文件列表**：客户端可以请求服务器的文件列表
- **文件下载**：客户端可以从服务器下载文件并指定保存路径
- **DHT 发现**：使用 Kademlia DHT 进行节点发现和路由
- **完整客户端**：支持连接、列出文件、下载文件等完整功能

## 架构

项目使用以下组件：

- **libp2p**：P2P 网络层，提供主机、流和协议支持
- **Kademlia DHT**：分布式哈希表，用于节点发现
- **自定义协议**：`/p2file/<host_id>` 用于文件共享通信
- **JSON Payload**：基于 JSON 的消息格式用于请求和响应

### 核心组件

- `App`：主应用结构体，管理 libp2p 主机和 DHT
- `Payload`：通信消息结构体，定义请求类型和数据
- `config`：命令行参数配置
- `cmd/client`：客户端入口点

### 协议消息

- `PL_LS`：请求文件列表
- `PL_LS_RES`：文件列表响应
- `PL_GET`：请求下载文件
- `PL_GET_RES`：文件数据响应
- `PL_GET_RES_DONE`：下载完成

## 安装

1. 确保安装了 Go 1.25.3 或更高版本
2. 克隆仓库：
   ```bash
   git clone https://github.com/blinkspark/p2file.git
   cd p2file
   ```
3. 安装依赖：
   ```bash
   go mod download
   ```
4. 构建：
   ```bash
   go build -o p2file cmd/client/main.go
   ```

## 使用

### 服务器模式

启动服务器共享目录：

```bash
./p2file -d /path/to/share
```

### 客户端模式

连接到服务器并列出文件：

```bash
./p2file -c <peer_id> -l
```

下载文件：

```bash
./p2file -c <peer_id> -g <filename>
```

指定下载保存路径：

```bash
./p2file -c <peer_id> -g <filename> -o /path/to/save
```

### 命令行参数

- `-c`：目标通道（对等节点 ID）
- `-d`：要共享的目标目录（默认：当前目录）
- `-g`：要获取的目标文件
- `-l`：列出所有文件
- `-o`：下载文件保存路径（可选，默认使用原文件名）

## 开发

### 项目结构

```
.
├── app.go          # 核心应用逻辑
├── payload.go      # 通信 payload 定义
├── config/
│   └── config.go   # 命令行配置
├── cmd/
│   └── client/
│       └── main.go # 客户端入口
├── go.mod
├── go.sum
└── README.md
```

### 依赖

主要依赖：
- `github.com/libp2p/go-libp2p`：libp2p 核心库
- `github.com/libp2p/go-libp2p-kad-dht`：Kademlia DHT 实现

## 待办事项

### 核心功能（已完成）

- [x] 实现客户端连接逻辑：解析 peer ID 并建立连接
- [x] 实现文件列表请求：发送 PL_LS payload 并处理 PL_LS_RES 响应
- [x] 实现文件下载：发送 PL_GET payload，接收 PL_GET_RES 数据流并保存文件
- [x] 添加 peer 发现机制：使用 DHT 查找目标 peer
- [x] 支持自定义输出路径

### 功能增强

- [ ] 实现多文件下载支持
- [ ] 添加下载进度显示
- [ ] 添加错误处理和重试机制

### 服务器端改进

- [ ] 支持目录递归共享
- [ ] 实现文件分块传输以提高大文件传输效率
- [ ] 添加并发连接处理
- [ ] 优化错误处理和日志记录

### 安全和认证

- [ ] 添加认证和授权机制
- [ ] 实现传输加密
- [ ] 添加文件完整性验证（哈希检查）

### 用户界面和体验

- [ ] 添加下载进度显示
- [ ] 添加命令行界面改进（彩色输出、进度条）
- [ ] 实现配置文件支持
- [ ] 添加 GUI 界面

### 网络和性能

- [ ] 优化 DHT 引导过程
- [ ] 添加连接池管理
- [ ] 实现断点续传功能
- [ ] 实现多文件批量下载

## 许可证

请查看 LICENSE 文件。

## 贡献

欢迎提交 issue 和 pull request。