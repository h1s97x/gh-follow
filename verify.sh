#!/bin/bash

# 验证脚本 - 检查项目结构和导入路径

echo "========================================"
echo "GH-Follow 项目结构验证"
echo "========================================"

# 检查目录结构
echo -e "\n📁 检查目录结构..."
required_dirs=(
    "cmd/gh-follow"
    "internal/models"
    "internal/storage"
    "internal/github"
    "internal/config"
    "internal/cache"
    "internal/sync"
    "internal/suggest"
    "internal/cmd"
    "internal/flags"
    "internal/errors"
)

for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "  ✅ $dir"
    else
        echo "  ❌ $dir (缺失)"
    fi
done

# 检查关键文件
echo -e "\n📄 检查关键文件..."
required_files=(
    "go.mod"
    "Makefile"
    "README.md"
    "cmd/gh-follow/main.go"
    "cmd/gh-follow/version.go"
    "internal/models/models.go"
    "internal/storage/storage.go"
    "internal/github/client.go"
    "internal/config/config.go"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "  ✅ $file"
    else
        echo "  ❌ $file (缺失)"
    fi
done

# 检查 go.mod
echo -e "\n📦 检查 go.mod..."
if grep -q "module github.com/h1s97x/gh-follow" go.mod; then
    echo "  ✅ 模块路径正确: github.com/h1s97x/gh-follow"
else
    echo "  ❌ 模块路径不正确"
fi

# 检查导入路径
echo -e "\n🔍 检查导入路径..."
import_count=$(grep -r "github.com/h1s97x/gh-follow" --include="*.go" . 2>/dev/null | wc -l)
echo "  找到 $import_count 处导入 github.com/h1s97x/gh-follow"

# 检查旧的导入路径
old_imports=$(grep -r "github.com/Link-/gh-follow" --include="*.go" . 2>/dev/null | wc -l)
if [ "$old_imports" -gt 0 ]; then
    echo "  ⚠️  发现 $old_imports 处旧导入路径需要更新"
else
    echo "  ✅ 无旧导入路径"
fi

# 统计文件数量
echo -e "\n📊 项目统计..."
go_files=$(find . -name "*.go" -not -path "./vendor/*" | wc -l)
test_files=$(find . -name "*_test.go" | wc -l)
source_files=$((go_files - test_files))

echo "  Go 源文件: $source_files"
echo "  测试文件: $test_files"
echo "  总计: $go_files 个 Go 文件"

# 检查依赖
echo -e "\n📚 检查依赖..."
if [ -f "go.sum" ]; then
    echo "  ✅ go.sum 存在"
else
    echo "  ⚠️  go.sum 不存在，需要运行 'go mod tidy'"
fi

echo -e "\n========================================"
echo "验证完成！"
echo "========================================"
echo ""
echo "下一步操作："
echo "1. 在有 Go 环境的机器上运行: go mod tidy"
echo "2. 构建项目: make build"
echo "3. 运行测试: make test"
echo "4. 提交代码: git add . && git commit -m 'feat: restructure project with internal packages'"
