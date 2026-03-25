# GH-Follow 项目结构说明

## 项目概览

GH-Follow 是一个 GitHub CLI 扩展，提供终端管理 GitHub 关注列表的功能。本项目采用模块化包结构，遵循 Go 最佳实践。

## 目录结构

```
gh-follow/
├── cmd/                          # 命令入口
│   └── gh-follow/
│       ├── main.go               # 主程序入口
│       └── version.go            # 版本信息
│
├── internal/                     # 内部包（不对外暴露）
│   ├── cache/                    # 缓存功能
│   │   └── cache.go              # 缓存实现
│   │
│   ├── cmd/                      # 命令处理器
│   │   ├── add.go                # follow 命令
│   │   ├── remove.go             # unfollow 命令
│   │   ├── list.go               # list 命令
│   │   ├── sync.go               # sync 命令
│   │   ├── stats.go              # stats 命令
│   │   ├── batch.go              # batch 命令
│   │   ├── suggest.go            # suggest 命令
│   │   ├── config.go             # config 命令
│   │   ├── gist.go               # gist 命令
│   │   ├── cache.go              # cache 命令
│   │   └── integration_test.go   # 集成测试
│   │
│   ├── config/                   # 配置管理
│   │   ├── config.go             # 配置读写
│   │   └── cmd.go                # 配置命令处理
│   │
│   ├── errors/                   # 错误定义
│   │   └── errors.go             # 结构化错误
│   │
│   ├── flags/                    # 命令行标志定义
│   │   ├── add.go                # add 命令标志
│   │   ├── remove.go             # remove 命令标志
│   │   ├── list.go               # list 命令标志
│   │   ├── sync.go               # sync 命令标志
│   │   ├── stats.go              # stats 命令标志
│   │   ├── batch.go              # batch 命令标志
│   │   ├── suggest.go            # suggest 命令标志
│   │   ├── config.go             # config 命令标志
│   │   ├── gist.go               # gist 命令标志
│   │   ├── cache.go              # cache 命令标志
│   │   ├── export.go             # export 命令标志
│   │   ├── import.go             # import 命令标志
│   │   └── autosync.go           # autosync 命令标志
│   │
│   ├── github/                   # GitHub API 客户端
│   │   ├── client.go             # API 客户端封装
│   │   └── batch.go              # 批量操作
│   │
│   ├── models/                   # 数据模型
│   │   ├── models.go             # 核心数据结构
│   │   └── models_test.go        # 模型测试
│   │
│   ├── storage/                  # 本地存储
│   │   ├── storage.go            # 存储实现
│   │   ├── storage_test.go       # 存储测试
│   │   ├── export.go             # 导出功能
│   │   └── import.go             # 导入功能
│   │
│   ├── suggest/                  # 建议引擎
│   │   ├── engine.go             # 建议生成引擎
│   │   └── suggest.go            # suggest 命令处理
│   │
│   └── sync/                     # 同步功能
│       ├── gist_sync.go          # Gist 同步
│       ├── gist_sync_test.go     # Gist 同步测试
│       └── auto_sync.go          # 自动同步
│
├── .github/                      # GitHub 配置
│   └── workflows/
│       ├── ci.yml                # CI 工作流
│       └── release.yml           # 发布工作流
│
├── go.mod                        # Go 模块定义
├── go.sum                        # 依赖校验
├── Makefile                      # 构建脚本
├── README.md                     # 项目说明
├── CONTRIBUTING.md               # 贡献指南
├── SECURITY.md                   # 安全策略
└── LICENSE                       # MIT 许可证
```

## 包职责说明

### 1. `cmd/gh-follow`
命令入口点，包含：
- `main.go`: 定义 CLI 应用和所有子命令
- `version.go`: 版本信息

### 2. `internal/models`
核心数据模型：
- `FollowList`: 关注列表数据结构
- `FollowEntry`: 单个关注条目
- `Config`: 应用配置
- `SyncConfig`: 同步配置

### 3. `internal/storage`
本地存储管理：
- JSON 文件的读写
- 数据导入导出
- 路径管理

### 4. `internal/github`
GitHub API 封装：
- API 客户端初始化
- follow/unfollow 操作
- 批量操作支持
- Token 管理

### 5. `internal/config`
配置管理：
- 配置文件读写
- 配置项访问
- 配置更新

### 6. `internal/cache`
缓存系统：
- 用户信息缓存
- 缓存清理
- 过期管理

### 7. `internal/sync`
同步功能：
- Gist 云同步
- GitHub 同步
- 自动同步
- 冲突检测

### 8. `internal/suggest`
建议引擎：
- 基于网络分析的建议
- 趋势用户发现
- 活跃度检测

### 9. `internal/cmd`
命令处理器：
- 各子命令的业务逻辑
- 输出格式化
- 交互处理

### 10. `internal/flags`
命令行标志：
- 各命令的标志定义
- 标志验证
- 默认值设置

### 11. `internal/errors`
错误处理：
- 结构化错误类型
- 错误码定义
- 错误判断辅助函数

## 构建与测试

### 构建项目
```bash
# 安装依赖
make deps

# 构建
make build

# 多平台构建
make build-all

# 发布版本（优化体积）
make release
```

### 测试
```bash
# 运行测试
make test

# 测试覆盖率
make test-coverage
```

### 代码质量
```bash
# 格式化
make fmt

# 静态检查
make lint

# 所有检查
make check
```

## 依赖说明

### 主要依赖
- `github.com/urfave/cli/v2`: CLI 框架
- `github.com/google/go-github/v55`: GitHub API 客户端

### 间接依赖
- `golang.org/x/crypto`: 加密支持
- `golang.org/x/sys`: 系统调用

## 配置路径

### 数据存储
- 关注列表: `~/.config/gh/follow-list.json`
- 配置文件: `~/.config/gh-follow/config.json`
- 缓存目录: `~/.cache/gh-follow/`

### 环境变量
- `GH_TOKEN`: GitHub 个人访问令牌（可选，默认从 `gh` CLI 获取）

## 扩展安装

### 通过 GitHub CLI 安装
```bash
gh extension install h1s97x/gh-follow
```

### 从源码安装
```bash
git clone https://github.com/h1s97x/gh-follow.git
cd gh-follow
make install
```

## 注意事项

1. **模块路径**: 已更新为 `github.com/h1s97x/gh-follow`
2. **Go 版本**: 要求 Go 1.22 或更高
3. **权限**: 需要 GitHub Token 的 `user:follow` 权限
4. **并发**: 批量操作支持并发控制，默认限制并发数避免 API 限流

## 后续工作

1. 在有 Go 环境的机器上运行 `make test` 验证所有测试
2. 运行 `make build` 确保编译成功
3. 运行 `make lint` 检查代码规范
4. 提交代码并推送到 GitHub 仓库
