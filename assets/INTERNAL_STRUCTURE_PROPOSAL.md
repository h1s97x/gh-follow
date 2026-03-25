# Internal 目录细分方案

## 📊 当前 Internal 结构问题

```
internal/
├── models.go
├── storage.go
├── github_api.go
├── config.go
├── errors.go
├── cache.go
├── gist_sync.go
├── suggest.go
├── batch.go
├── auto_sync.go
├── add.go
├── add_flags.go
├── remove.go
├── remove_flags.go
├── list.go
├── list_flags.go
├── sync.go
├── sync_flags.go
├── export.go
├── export_flags.go
├── import.go
├── import_flags.go
├── stats.go
├── stats_flags.go
├── config_cmd.go
├── config_flags.go
├── gist_cmd.go
├── gist_flags.go
├── cache_cmd.go
├── cache_flags.go
├── suggest_cmd.go
├── suggest_flags.go
├── batch_cmd.go
├── batch_flags.go
├── auto_sync_cmd.go
├── auto_sync_flags.go
├── *_test.go          ❌ 测试混在一起
├── fixtures/
└── test-follows.json
```

**问题：**
- ❌ 40+ 个文件混在一个目录
- ❌ 测试文件和源文件混在一起
- ❌ 命令处理器和标志定义混在一起
- ❌ 难以快速定位相关代码
- ❌ 不符合 Go 项目最佳实践

---

## ✅ 推荐结构

```
internal/
├── models/                    # 数据模型
│   ├── follow.go
│   ├── config.go
│   ├── stats.go
│   └── models_test.go
│
├── storage/                   # 存储层
│   ├── storage.go
│   ├── export.go
│   ├── import.go
│   ├── storage_test.go
│   └── fixtures/
│       └── test-follows.json
│
├── github/                    # GitHub API 集成
│   ├── client.go              (github_api.go)
│   ├── batch.go
│   ├── github_test.go
│   └── batch_test.go
│
├── config/                    # 配置管理
│   ├── config.go
│   ├── config_test.go
│   └── defaults.go
│
├── cache/                     # 缓存系统
│   ├── cache.go
│   ├── cache_test.go
│   └── eviction.go
│
├── sync/                      # 同步功能
│   ├── gist.go                (gist_sync.go)
│   ├── auto.go                (auto_sync.go)
│   ├── gist_test.go
│   └── auto_test.go
│
├── suggest/                   # 推荐引擎
│   ├── engine.go              (suggest.go)
│   ├── strategies.go
│   ├── engine_test.go
│   └── strategies_test.go
│
├── cmd/                       # 命令处理器
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── sync.go
│   ├── export.go
│   ├── import.go
│   ├── stats.go
│   ├── config.go
│   ├── gist.go
│   ├── cache.go
│   ├── suggest.go
│   ├── batch.go
│   ├── autosync.go
│   ├── cmd_test.go
│   └── integration_test.go
│
├── flags/                     # 命令行标志定义
│   ├── add.go
│   ├── remove.go
│   ├── list.go
│   ├── sync.go
│   ├── export.go
│   ├── import.go
│   ├── stats.go
│   ├── config.go
│   ├── gist.go
│   ├── cache.go
│   ├── suggest.go
│   ├── batch.go
│   ├── autosync.go
│   └── flags_test.go
│
├── errors/                    # 错误处理
│   ├── errors.go
│   └── errors_test.go
│
└── testutil/                  # 测试工具（可选）
    ├── helpers.go
    └── mocks.go
```

---

## 🎯 分类说明

### 1️⃣ **models/** - 数据模型
**包含：** 所有数据结构定义

```go
// models/follow.go
type Follow struct { ... }
type FollowList struct { ... }

// models/config.go
type Config struct { ... }
type StorageConf struct { ... }

// models/stats.go
type Stats struct { ... }
```

**测试：** `models_test.go`

---

### 2️⃣ **storage/** - 存储层
**包含：** 本地文件存储操作

```go
// storage/storage.go
type Storage struct { ... }
func (s *Storage) Load() { ... }
func (s *Storage) Save() { ... }

// storage/export.go
func (s *Storage) Export() { ... }

// storage/import.go
func (s *Storage) Import() { ... }
```

**测试：** `storage_test.go`
**Fixtures：** `fixtures/test-follows.json`

---

### 3️⃣ **github/** - GitHub API 集成
**包含：** GitHub API 客户端和批量操作

```go
// github/client.go
type GitHubClient struct { ... }
func (gc *GitHubClient) Follow() { ... }
func (gc *GitHubClient) GetFollowing() { ... }

// github/batch.go
type BatchProcessor struct { ... }
func (bp *BatchProcessor) BatchFollow() { ... }
```

**测试：** `github_test.go`, `batch_test.go`

---

### 4️⃣ **config/** - 配置管理
**包含：** 配置文件读写和管理

```go
// config/config.go
type ConfigManager struct { ... }
func (cm *ConfigManager) Load() { ... }
func (cm *ConfigManager) Save() { ... }

// config/defaults.go
func DefaultConfigPath() { ... }
func NewDefaultConfig() { ... }
```

**测试：** `config_test.go`

---

### 5️⃣ **cache/** - 缓存系统
**包含：** 用户信息缓存

```go
// cache/cache.go
type Cache struct { ... }
func (c *Cache) Get() { ... }
func (c *Cache) Set() { ... }

// cache/eviction.go
func (c *Cache) evictOldest() { ... }
```

**测试：** `cache_test.go`

---

### 6️⃣ **sync/** - 同步功能
**包含：** Gist 同步和自动同步

```go
// sync/gist.go
type GistSync struct { ... }
func (gs *GistSync) Upload() { ... }
func (gs *GistSync) Download() { ... }

// sync/auto.go
type AutoSync struct { ... }
func (as *AutoSync) Trigger() { ... }
```

**测试：** `gist_test.go`, `auto_test.go`

---

### 7️⃣ **suggest/** - 推荐引擎
**包含：** 关注建议生成

```go
// suggest/engine.go
type SuggestionEngine struct { ... }
func (se *SuggestionEngine) GenerateSuggestions() { ... }

// suggest/strategies.go
func (se *SuggestionEngine) suggestFromFollowOfFollow() { ... }
func (se *SuggestionEngine) suggestFromOrgs() { ... }
```

**测试：** `engine_test.go`, `strategies_test.go`

---

### 8️⃣ **cmd/** - 命令处理器
**包含：** 所有命令的处理逻辑

```go
// cmd/add.go
func Add(c *cli.Context) error { ... }

// cmd/remove.go
func Remove(c *cli.Context) error { ... }

// cmd/list.go
func List(c *cli.Context) error { ... }

// ... 其他命令
```

**测试：** `cmd_test.go`, `integration_test.go`

---

### 9️⃣ **flags/** - 命令行标志定义
**包含：** 所有命令的标志定义

```go
// flags/add.go
func AddFlags() []cli.Flag { ... }

// flags/remove.go
func RemoveFlags() []cli.Flag { ... }

// flags/list.go
func ListFlags() []cli.Flag { ... }

// ... 其他标志
```

**测试：** `flags_test.go`

---

### 🔟 **errors/** - 错误处理
**包含：** 自定义错误类型

```go
// errors/errors.go
type FollowError struct { ... }
var (
    ErrEmptyUsername = errors.New("...")
    ErrUserNotFound = errors.New("...")
)
```

**测试：** `errors_test.go`

---

### 1️⃣1️⃣ **testutil/** - 测试工具（可选）
**包含：** 测试辅助函数和 Mock

```go
// testutil/helpers.go
func CreateTestFollowList() *FollowList { ... }
func CreateTestStorage() *Storage { ... }

// testutil/mocks.go
type MockGitHubClient struct { ... }
```

---

## 📊 结构对比

### 旧结构
```
internal/
├── 40+ 个文件混在一起
├── 测试文件和源文件混在一起
└── 难以快速定位代码
```

### 新结构
```
internal/
├── models/          (数据模型)
├── storage/         (存储层)
├── github/          (GitHub API)
├── config/          (配置管理)
├── cache/           (缓存系统)
├── sync/            (同步功能)
├── suggest/         (推荐引擎)
├── cmd/             (命令处理)
├── flags/           (标志定义)
├── errors/          (错误处理)
└── testutil/        (测试工具)
```

**改进：** 从 40+ 个文件混在一起 → 11 个清晰的子包

---

## 🔄 迁移步骤

### 第1步：创建子包目录
```bash
mkdir -p internal/models
mkdir -p internal/storage/fixtures
mkdir -p internal/github
mkdir -p internal/config
mkdir -p internal/cache
mkdir -p internal/sync
mkdir -p internal/suggest
mkdir -p internal/cmd
mkdir -p internal/flags
mkdir -p internal/errors
mkdir -p internal/testutil
```

### 第2步：移动文件

**models/**
```bash
mv internal/models.go internal/models/
mv internal/models_test.go internal/models/
```

**storage/**
```bash
mv internal/storage.go internal/storage/
mv internal/export.go internal/storage/
mv internal/import.go internal/storage/
mv internal/storage_test.go internal/storage/
mv internal/fixtures/ internal/storage/
```

**github/**
```bash
mv internal/github_api.go internal/github/client.go
mv internal/batch.go internal/github/
mv internal/github_test.go internal/github/
mv internal/batch_test.go internal/github/
```

**config/**
```bash
mv internal/config.go internal/config/
mv internal/config_test.go internal/config/
```

**cache/**
```bash
mv internal/cache.go internal/cache/
mv internal/cache_test.go internal/cache/
```

**sync/**
```bash
mv internal/gist_sync.go internal/sync/gist.go
mv internal/auto_sync.go internal/sync/auto.go
mv internal/gist_sync_test.go internal/sync/gist_test.go
mv internal/auto_sync_test.go internal/sync/auto_test.go
```

**suggest/**
```bash
mv internal/suggest.go internal/suggest/engine.go
mv internal/suggest_test.go internal/suggest/engine_test.go
```

**cmd/**
```bash
mv internal/add.go internal/cmd/
mv internal/remove.go internal/cmd/
mv internal/list.go internal/cmd/
mv internal/sync.go internal/cmd/
mv internal/export.go internal/cmd/
mv internal/import.go internal/cmd/
mv internal/stats.go internal/cmd/
mv internal/config_cmd.go internal/cmd/config.go
mv internal/gist_cmd.go internal/cmd/gist.go
mv internal/cache_cmd.go internal/cmd/cache.go
mv internal/suggest_cmd.go internal/cmd/suggest.go
mv internal/batch_cmd.go internal/cmd/batch.go
mv internal/auto_sync_cmd.go internal/cmd/autosync.go
mv internal/cmd_test.go internal/cmd/
mv internal/integration_test.go internal/cmd/
```

**flags/**
```bash
mv internal/add_flags.go internal/flags/add.go
mv internal/remove_flags.go internal/flags/remove.go
mv internal/list_flags.go internal/flags/list.go
mv internal/sync_flags.go internal/flags/sync.go
mv internal/export_flags.go internal/flags/export.go
mv internal/import_flags.go internal/flags/import.go
mv internal/stats_flags.go internal/flags/stats.go
mv internal/config_flags.go internal/flags/config.go
mv internal/gist_flags.go internal/flags/gist.go
mv internal/cache_flags.go internal/flags/cache.go
mv internal/suggest_flags.go internal/flags/suggest.go
mv internal/batch_flags.go internal/flags/batch.go
mv internal/auto_sync_flags.go internal/flags/autosync.go
mv internal/flags_test.go internal/flags/
```

**errors/**
```bash
mv internal/errors.go internal/errors/
mv internal/errors_test.go internal/errors/
```

### 第3步：更新导入路径

**在 cmd/gh-follow/main.go 中：**

```go
// 旧的
import "github.com/h1s97x/gh-follow/internal"

// 新的
import (
    "github.com/h1s97x/gh-follow/internal/cmd"
    "github.com/h1s97x/gh-follow/internal/flags"
)

// 使用
Commands: []*cli.Command{
    {
        Name: "add",
        Flags: flags.AddFlags(),
        Action: cmd.Add,
    },
    // ...
}
```

**在各个包中：**

```go
// internal/cmd/add.go
package cmd

import (
    "github.com/h1s97x/gh-follow/internal/models"
    "github.com/h1s97x/gh-follow/internal/storage"
    "github.com/h1s97x/gh-follow/internal/github"
)

// internal/github/client.go
package github

import (
    "github.com/h1s97x/gh-follow/internal/models"
    "github.com/h1s97x/gh-follow/internal/errors"
)
```

### 第4步：验证

```bash
# 验证构建
go build ./cmd/gh-follow

# 验证测试
go test ./...

# 验证帮助
./gh-follow --help
```

---

## 📋 包依赖关系

```
models/          (基础数据结构)
  ↑
  ├─ storage/    (使用 models)
  ├─ github/     (使用 models)
  ├─ config/     (使用 models)
  ├─ cache/      (使用 models)
  ├─ sync/       (使用 models)
  └─ suggest/    (使用 models)
       ↑
       ├─ cmd/   (使用所有包)
       └─ flags/ (定义命令行标志)
```

---

## ✅ 优势

✅ **清晰的包结构**
- 每个包有明确的职责
- 易于理解和维护

✅ **测试组织清晰**
- 测试文件和源文件在同一包
- 便于运行特定包的测试

✅ **代码复用性高**
- 各个包可以独立使用
- 便于单元测试

✅ **易于扩展**
- 添加新功能时，创建新包即可
- 不会污染现有包

✅ **符合 Go 最佳实践**
- 遵循 Go 包设计原则
- 便于其他开发者理解

---

## 🔍 文件数量对比

| 指标 | 旧结构 | 新结构 | 改进 |
|------|--------|--------|------|
| internal 目录文件 | 40+ | 11 个子包 | **-73%** |
| 单个目录最大文件数 | 40+ | 5-8 | **-80%** |
| 查找相关代码时间 | 长 | 短 | ✅ |

---

## 📝 提交计划更新

如果采用新的 internal 结构，提交计划需要微调：

```
1. chore: initialize project structure and configuration
   - 创建目录结构（包括 internal 子包）

2. feat: implement core data models
   - 添加 internal/models/

3. feat: implement storage layer
   - 添加 internal/storage/

4. feat: implement GitHub API integration
   - 添加 internal/github/

5. feat: implement configuration management
   - 添加 internal/config/

6. feat: implement caching system
   - 添加 internal/cache/

7. feat: implement sync functionality
   - 添加 internal/sync/

8. feat: implement suggestion engine
   - 添加 internal/suggest/

9. feat: implement error handling
   - 添加 internal/errors/

10. feat: implement CLI commands
    - 添加 internal/cmd/

11. feat: implement command flags
    - 添加 internal/flags/

12. test: add comprehensive unit tests
    - 添加所有 *_test.go 文件

13. docs: add documentation and CI/CD configuration
    - 添加文档和工作流
```

---

## 🎯 总结

**建议：** 采用细分的 internal 结构

**优势：**
- ✅ 清晰的包结构（11 个子包）
- ✅ 测试组织清晰
- ✅ 代码复用性高
- ✅ 易于扩展
- ✅ 符合 Go 最佳实践

**工作量：** 约 20 分钟（移动文件 + 更新导入）

**风险：** 低（只是文件位置改变）

---

**你同意这个方案吗？** 👍
