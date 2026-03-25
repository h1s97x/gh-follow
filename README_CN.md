# gh-follow 使用指南

[English](README.md) | 中文

一个强大的 GitHub CLI 扩展，让你在终端管理 GitHub 关注列表。

## 安装

### 方法 1：作为 GitHub CLI 扩展安装（推荐）

```bash
gh extension install h1s97x/gh-follow
```

### 方法 2：从源码构建

```bash
git clone https://github.com/h1s97x/gh-follow.git
cd gh-follow
go build -o bin/gh-follow ./cmd/gh-follow

# 安装到 gh 扩展目录
mkdir -p ~/.local/share/gh/extensions/gh-follow
cp bin/gh-follow ~/.local/share/gh/extensions/gh-follow/
```

### 方法 3：下载预编译二进制

从 [GitHub Releases](https://github.com/h1s97x/gh-follow/releases) 下载最新版本。

## 前置要求

- [GitHub CLI](https://cli.github.com/) 已安装并登录
- Go 1.22+（从源码构建需要）

```bash
# 检查登录状态
gh auth status

# 如果未登录
gh auth login
```

## 快速开始

```bash
# 查看帮助
gh follow --help

# 查看版本
gh follow --version

# 列出关注的用户
gh follow list

# 查看统计信息
gh follow stats
```

## 命令详解

### 📋 列表管理

#### 列出关注的用户

```bash
# 列出所有关注的用户
gh follow list

# 别名
gh follow ls

# JSON 格式输出
gh follow list --format json

# 简单列表格式（只显示用户名）
gh follow list --format simple

# 限制数量
gh follow list --limit 50

# 排序
gh follow list --sort name --order asc
gh follow list --sort date --order desc

# 按标签筛选
gh follow list --tag developer

# 按日期筛选
gh follow list --date-from 2024-01-01 --date-to 2024-03-31
```

#### 查看统计数据

```bash
gh follow stats

# 输出示例：
# Total Following: 128
# Mutual Followers: 45
# Non-followers: 83
# Organizations: 12
```

---

### ➕ 关注用户

#### 关注单个用户

```bash
gh follow add octocat

# 别名
gh follow follow octocat
```

#### 关注多个用户

```bash
gh follow add octocat torvalds github
```

#### 带备注和标签关注

```bash
gh follow add octocat --notes "GitHub mascot" --tags developer,go
```

#### 从文件批量关注

```bash
# users.txt 每行一个用户名
gh follow add --file users.txt
```

---

### ➖ 取消关注

#### 取消关注单个用户

```bash
gh follow remove octocat

# 别名
gh follow unfollow octocat
gh follow rm octocat
gh follow delete octocat
```

#### 取消关注多个用户

```bash
gh follow remove octocat torvalds github
```

#### 跳过确认

```bash
gh follow remove octocat --force
```

#### 从文件批量取消关注

```bash
gh follow remove --file users.txt
```

---

### 🔄 批量操作

```bash
# 查看批量操作帮助
gh follow batch --help

# 批量关注
gh follow batch follow user1 user2 user3

# 从文件批量关注
gh follow batch follow --file users.txt

# 预览模式（不执行）
gh follow batch follow user1 user2 --dry-run

# 批量取消关注
gh follow batch unfollow user1 user2 user3

# 批量检查是否关注你
gh follow batch check user1 user2 user3

# 从文件导入用户名
gh follow batch import users.txt

# 导入并关注所有用户
gh follow batch import users.txt --follow
```

---

### 💡 关注建议

```bash
# 获取个性化关注建议
gh follow suggest

# 别名
gh follow recommend

# 限制建议数量
gh follow suggest --limit 10

# 查看热门用户
gh follow suggest trending

# 按编程语言筛选热门用户
gh follow suggest trending --language go

# 检查互相关注
gh follow suggest mutual

# 查找不活跃的关注用户
gh follow suggest inactive --days 365
```

---

### 💾 数据导入导出

#### 导出关注列表

```bash
# 导出到 JSON 文件
gh follow export --output following.json

# 导出为 CSV 格式
gh follow export --format csv --output following.csv
```

#### 导入关注列表

```bash
# 从文件导入
gh follow import --input following.json

# 从 CSV 导入
gh follow import --format csv --input following.csv

# 导入并合并
gh follow import --input following.json --merge
```

---

### ☁️ 云端同步

#### Gist 同步

```bash
# 创建同步用的 Gist
gh follow gist create

# 查看 Gist 同步状态
gh follow gist status

# 从 Gist 拉取
gh follow gist pull

# 推送到 Gist
gh follow gist push

# 同步时指定 Gist
gh follow sync --gist --gist-id <gist-id>
```

#### 自动同步

```bash
# 查看自动同步状态
gh follow autosync status

# 手动触发同步
gh follow autosync trigger

# 启用自动同步
gh follow config set sync.auto_sync true

# 禁用自动同步
gh follow config set sync.auto_sync false

# 设置同步间隔（秒）
gh follow config set sync.sync_interval 3600
```

#### GitHub 同步

```bash
# 双向同步（默认）
gh follow sync

# 从 GitHub 拉取
gh follow sync --direction pull

# 推送到 GitHub
gh follow sync --direction push

# 预览模式
gh follow sync --dry-run
```

---

### ⚡ 缓存管理

```bash
# 查看缓存状态
gh follow cache status

# 列出缓存的用户
gh follow cache list

# 清除缓存
gh follow cache clear --force

# 清理过期缓存
gh follow cache cleanup

# 刷新所有关注用户的缓存
gh follow cache refresh

# 查看缓存用户详情
gh follow cache show octocat
```

---

### ⚙️ 配置管理

```bash
# 查看当前配置
gh follow config

# 查看当前配置（详细）
gh follow config show

# 获取特定配置项
gh follow config get display.default_format

# 设置配置项
gh follow config set display.default_format json

# 重置为默认配置
gh follow config reset --force
```

## 配置项说明

| 配置项 | 说明 | 默认值 |
|-------|------|--------|
| `storage.local_path` | 本地存储路径 | `~/.config/gh/follow-list.json` |
| `storage.use_gist` | 启用 Gist 同步 | `false` |
| `storage.gist_id` | Gist ID | `""` |
| `sync.auto_sync` | 启用自动同步 | `true` |
| `sync.sync_interval` | 同步间隔（秒） | `3600` |
| `display.default_format` | 默认输出格式 | `table` |
| `display.default_sort` | 默认排序字段 | `date` |
| `display.default_order` | 默认排序顺序 | `desc` |

## 冲突解决

同步时可能发生冲突，支持三种解决策略：

1. **newest-wins**（默认）：使用最新的条目
2. **local-wins**：本地优先
3. **remote-wins**：远程优先

使用 `--force` 标志可跳过冲突检测。

## 数据存储

### 本地存储

存储位置：`~/.config/gh/follow-list.json`

```json
{
  "version": "1.0.0",
  "updated_at": "2024-03-25T10:30:00Z",
  "follows": [
    {
      "username": "octocat",
      "followed_at": "2024-03-20T15:30:00Z",
      "notes": "GitHub mascot",
      "tags": ["developer", "go"]
    }
  ],
  "metadata": {
    "total_count": 1,
    "last_sync": "2024-03-25T10:30:00Z",
    "sync_status": "success"
  }
}
```

### Gist 存储

启用 Gist 同步后，会创建私有 Gist 存储：
- 跨设备同步
- 版本历史
- 备份恢复

## 常见使用场景

### 场景 1：备份关注列表

```bash
# 导出到本地
gh follow export --output backup.json

# 同时同步到 Gist
gh follow sync --gist
```

### 场景 2：清理不互关的用户

```bash
# 1. 导出当前列表
gh follow export --output following.json

# 2. 查看统计（找出不互关的）
gh follow stats

# 3. 检查互相关注状态
gh follow suggest mutual

# 4. 批量取消关注
gh follow batch unfollow --file non_mutual.txt
```

### 场景 3：迁移到新账号

```bash
# 旧账号导出
gh follow export --output following.json

# 新账号导入
gh follow import --input following.json
```

### 场景 4：发现有趣的用户

```bash
# 获取推荐
gh follow suggest --limit 20

# 查看热门用户
gh follow suggest trending --language go

# 选择几个关注
gh follow add user1 user2 user3
```

## 开发

### 构建

```bash
make build
# 或
go build -o bin/gh-follow ./cmd/gh-follow
```

### 测试

```bash
make test
# 或
go test ./...
```

### 本地安装

```bash
make install
```

## 许可证

[MIT License](LICENSE)

## 相关项目

- [GitHub CLI](https://github.com/cli/cli) - 官方 GitHub CLI
- [gh-token](https://github.com/Link-/gh-token) - GitHub App token 管理
