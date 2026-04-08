# issue2md 工具规格说明文档

## 1. 概述

`issue2md` 是一个命令行工具，用于将 GitHub Issues、Pull Requests 和 Discussions 转换为格式化的 Markdown 文档。

### 1.1 目标用户

- 需要归档 GitHub 讨论的开发者
- 想要离线阅读 Issue/PR 的用户
- 需要生成技术文档的项目维护者

### 1.2 核心价值

- 一键转换，操作简单
- 保留完整上下文（评论、时间轴、元数据）
- 支持离线归档

---

## 2. 功能需求

### 2.1 支持的资源类型

| 资源类型 | URL 示例 | 优先级 |
|---------|---------|--------|
| Issue | `https://github.com/owner/repo/issues/123` | P0 |
| Pull Request | `https://github.com/owner/repo/pull/456` | P0 |
| Discussion | `https://github.com/owner/repo/discussions/789` | P1 |

#### 2.1.1 URL 解析规则

工具通过 URL 路径中的关键词识别资源类型：

| 关键词 | 资源类型 | 示例路径 |
|--------|---------|---------|
| `/issues/` | Issue | `github.com/owner/repo/issues/123` |
| `/pull/` | Pull Request | `github.com/owner/repo/pull/456` |
| `/discussions/` | Discussion | `github.com/owner/repo/discussions/789` |

**解析逻辑**：
1. 提取 URL 路径部分
2. 匹配关键词（`issues`, `pull`, `discussions`）
3. 根据匹配结果确定资源类型
4. 提取 `owner`、`repo`、`number` 等参数

**支持的 URL 格式变体**：
- 标准格式：`https://github.com/owner/repo/issues/123`
- 带 `.git` 后缀：`https://github.com/owner/repo.git/issues/123`
- 带 www 子域：`https://www.github.com/owner/repo/pull/456`
- HTTP 协议：`http://github.com/owner/repo/issues/123`

### 2.2 输出内容

#### 2.2.1 元数据（YAML Front Matter）

```yaml
---
title: "Issue 标题"
type: issue | pull_request | discussion
number: 123
repository: owner/repo
author: username
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-01-02T00:00:00Z
state: open | closed | merged
labels:
  - bug
  - enhancement
milestone: "v1.0.0"
reactions:
  thumbs_up: 5
  thumbs_down: 0
  laugh: 2
  hooray: 1
  confused: 0
  heart: 3
  rocket: 0
  eyes: 0
---
```

#### 2.2.2 正文内容

- Issue/PR 标题（H1）
- 描述内容（兼容当前生成的文档格式，确保自洽性）
- 标签、里程碑等信息（作为元数据块显示）
- PR 特有：变更文件列表 + 折叠的 diff（`<details>` 标签）

#### 2.2.3 评论区

- 所有评论按时间顺序排列
- 嵌套回复使用引用块（`>`）缩进
- 每条评论包含：作者、时间戳、内容、反应
- 保留代码块、图片、链接的原始格式

---

## 3. 命令行接口

### 3.1 语法

```bash
issue2md [flags] <github_url> [output_file]
```

### 3.2 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| `github_url` | 必需 | GitHub Issue/PR/Discussion 的完整 URL |
| `output_file` | 可选 | 输出文件路径。如果不提供，输出到 stdout |

### 3.3 Flags

| Flag | 说明 |
|------|------|
| `-h`, `-help` | 显示帮助信息并退出 |
| `-v`, `-version` | 显示版本信息并退出 |
| `-enable-reactions` | 包含 GitHub 反应统计（👍👎😄🎉😕❤️🚀👀） |
| `-enable-user-links` | 将 @username 转换为可点击的 GitHub 链接 |

### 3.4 环境变量

| 变量名 | 必需 | 说明 |
|--------|------|------|
| `GITHUB_TOKEN` | 是 | GitHub Personal Access Token |

**注意**：Token **只**通过环境变量获取，不提供 `--token` 参数，以防在 Shell 历史中泄露密钥。

### 3.5 使用示例

```bash
# 显示帮助信息
issue2md -h

# 显示版本信息
issue2md -v

# 基本用法：输出到 stdout（可重定向）
export GITHUB_TOKEN=ghp_xxx
issue2md https://github.com/owner/repo/issues/123 > output.md

# 直接指定输出文件
issue2md https://github.com/owner/repo/issues/123 output.md

# 启用反应统计
issue2md -enable-reactions https://github.com/owner/repo/issues/123 > output.md

# 启用用户链接
issue2md -enable-user-links https://github.com/owner/repo/issues/123 > output.md

# 组合使用
issue2md -enable-reactions -enable-user-links https://github.com/owner/repo/pull/456 pr.md

# 管道到其他工具
issue2md https://github.com/owner/repo/issues/123 | pandoc -o output.pdf
```

### 3.6 帮助信息格式

```
issue2md - Convert GitHub Issues/PRs/Discussions to Markdown

USAGE:
    issue2md [flags] <github_url> [output_file]

ARGUMENTS:
    github_url        GitHub Issue/PR/Discussion URL
    output_file       Output file path (default: stdout)

FLAGS:
    -h, -help              Show this help message
    -v, -version           Show version information
    -enable-reactions      Include GitHub reactions
    -enable-user-links     Convert @username to links

ENVIRONMENT:
    GITHUB_TOKEN           GitHub Personal Access Token (required)

EXAMPLES:
    issue2md https://github.com/owner/repo/issues/123
    issue2md https://github.com/owner/repo/issues/123 > output.md
    issue2md -enable-reactions https://github.com/owner/repo/pull/456 pr.md
```

---

## 4. 输出格式

### 4.1 格式说明

- **唯一格式**：GitHub Flavored Markdown (GFM)
- **不提供模板系统**：保持简单，固定的输出格式
- **不提供 HTML/JSON 输出**：专注 Markdown，用户可通过管道转换

### 4.2 输出目标

| 场景 | 输出方式 |
|------|---------|
| 默认 | 标准输出 (stdout) |
| 指定 output_file | 写入指定文件 |
| 管道 | 可与其他工具组合 (如 `pandoc`) |

---

## 5. 文件命名规则（自动命名）

### 5.1 命名格式

```
{number}-{type}-{sanitized-title}.md
```

其中 `{type}` 为资源类型关键词：
- Issue: `issue`
- Pull Request: `pr`
- Discussion: `discussion`

### 5.2 标题清理规则

1. **移除特殊字符**：`/`, `\`, `:`, `*`, `?`, `"`, `<`, `>`, `|`, `.`
2. **空格替换**：空格替换为连字符 `-`
3. **长度限制**：最多保留 50 个字符
4. **连续连字符**：多个连续 `-` 合并为一个
5. **首尾处理**：移除首尾的连字符
6. **Unicode 保留**：保留中文等非 ASCII 字符

### 5.3 示例

| 资源类型 | 原标题 | 清理后标题 | 完整文件名 |
|---------|--------|-----------|-----------|
| Issue | "Bug: 修复登录失败" | "bug-修复登录失败" | `123-issue-bug-修复登录失败.md` |
| PR | "Feature: Add ??? support" | "feature-add-support" | `456-pr-feature-add-support.md` |
| Discussion | "A very long title..." | "a-very-long-title" | `789-discussion-a-very-long-title.md` |

---

## 6. PR Diff 处理

### 6.1 Diff 输出格式

```markdown
## 变更文件

- [`path/to/file1.ts`](https://github.com/owner/repo/blob/branch/path/to/file1.ts) (+10, -5)
- [`path/to/file2.ts`](https://github.com/owner/repo/blob/branch/path/to/file2.ts) (+100, -0)

<details>
<summary>查看完整 Diff</summary>

```diff
--- a/path/to/file1.ts
+++ b/path/to/file1.ts
@@ -1,3 +1,5 @@
 function hello() {
+  console.log("hello");
   return "world";
 }
```

</details>
```

### 6.2 大型 PR 处理

- Diff 超过 500 行时，显示提示信息并截断
- 在截断处添加链接指向 GitHub 的完整 diff

---

## 7. 错误处理

### 7.1 错误类型

| 错误类型 | 处理方式 | 退出码 |
|---------|---------|--------|
| 正常退出 | - | 0 |
| 用户请求帮助 | `-h` 或 `-help` | 0 |
| 用户请求版本 | `-v` 或 `-version` | 0 |
| URL 格式无效 | 提示正确格式，退出 | 1 |
| 资源不存在 | 提示检查 URL 和权限，退出 | 2 |
| API 请求失败 | 显示错误信息，退出 | 3 |
| GITHUB_TOKEN 未设置 | 提示设置环境变量，退出 | 4 |
| 文件写入失败 | 显示错误信息，退出 | 5 |
| 网络超时 | 显示错误信息，退出 | 6 |

### 7.2 版本信息格式

```
issue2md version 1.0.0
```

### 7.3 错误消息格式

```
error: invalid URL format
usage: issue2md [flags] <github_url> [output_file]
```

---

## 8. 非功能需求

### 8.1 性能

- 单个 Issue/PR 转换时间 < 5 秒
- 支持 100+ 条评论的 Issue 转换

### 8.2 兼容性

- Go 版本 >= 1.21
- 支持 macOS、Linux、Windows

### 8.3 依赖

- GitHub REST API v3
- 外部依赖：仅 Go 标准库

---

## 9. 实现里程碑

### Phase 1: MVP（最小可行产品）

- [ ] 命令行参数解析（`-h`, `-v`, flags）
- [ ] 帮助信息显示
- [ ] 版本信息显示
- [ ] 支持 Issue 转换
- [ ] 基本元数据输出
- [ ] 评论获取和格式化
- [ ] stdout 输出和文件输出

### Phase 2: PR 支持

- [ ] PR 信息获取
- [ ] 变更文件列表
- [ ] 折叠 Diff 输出

### Phase 3: Discussion 支持

- [ ] Discussion API 集成
- [ ] Discussion 评论获取

### Phase 4: 增强功能

- [ ] 进度条显示
- [ ] 缓存支持
- [ ] 批量转换

---

## 11. 测试策略 (Test Strategy)

**本章节遵循 `constitution.md` 第二条：测试先行铁律，内容不可协商。**

### 11.1 核心原则

| 原则 | 说明 |
|------|------|
| **TDD 强制** | 所有新功能或 Bug 修复，必须从编写失败的测试开始 |
| **表格驱动** | 单元测试必须采用表格驱动测试（Table-Driven Tests）风格 |
| **真实依赖** | 优先编写集成测试，使用真实的 GitHub API，禁止使用 Mock |
| **覆盖率门槛** | 核心业务逻辑覆盖率 ≥ 80% |

### 11.2 TDD 开发流程（Red-Green-Refactor）

```go
// 步骤 1: Red - 编写失败的测试
func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    *ParsedURL
        wantErr bool
    }{
        {
            name:    "valid issue URL",
            url:     "https://github.com/owner/repo/issues/123",
            want:    &ParsedURL{Owner: "owner", Repo: "repo", Number: 123, Type: "issue"},
            wantErr: false,
        },
        {
            name:    "invalid URL format",
            url:     "not-a-url",
            want:    nil,
            wantErr: true,
        },
    }
    // 此时测试会失败，因为 ParseURL 函数还未实现
}

// 步骤 2: Green - 编写最少代码使测试通过
func ParseURL(rawURL string) (*ParsedURL, error) {
    // 实现逻辑...

    // 步骤 3: Refactor - 重构代码，保持测试通过
}

// 步骤 4: 提交代码
```

### 11.3 表格驱动测试规范

**格式要求**：

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string           // 测试用例名称
        input   InputType        // 输入参数
        want    ExpectedType     // 期望结果
        wantErr bool             // 是否期望错误
    }{
        {
            name:    "描述性的测试名称",
            input:   "test input",
            want:    "expected output",
            wantErr: false,
        },
        // 更多测试用例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := FunctionName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("FunctionName() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("FunctionName() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**命名规范**：
- 测试文件：`{filename}_test.go`
- 测试函数：`Test{FunctionName}`
- 子测试名称：使用 `t.Run()` 并提供描述性名称

### 11.4 集成测试 vs 单元测试

| 测试类型 | 范围 | 示例 | 优先级 |
|---------|------|------|--------|
| **集成测试** | 跨包/外部依赖 | GitHub API 调用、端到端转换 | **高** |
| **单元测试** | 单个函数/方法 | URL 解析、Markdown 格式化 | 中 |

**集成测试要求**：

```go
// internal/github/github_test.go
func TestFetchIssue_Integration(t *testing.T) {
    // 跳过单元测试模式
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // 使用真实环境变量
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        t.Skip("GITHUB_TOKEN not set")
    }

    // 测试真实 GitHub API
    issue, err := FetchIssue("golang", "go", 12345)
    if err != nil {
        t.Fatalf("FetchIssue() failed: %v", err)
    }

    if issue.Number != 12345 {
        t.Errorf("got number %d, want 12345", issue.Number)
    }
}
```

**禁止使用 Mock**：

```go
// ❌ 错误示范：使用 Mock
type MockGitHubClient struct{}
func (m *MockGitHubClient) FetchIssue() {...}

// ✅ 正确示范：使用真实 API
func TestFetchIssue(t *testing.T) {
    // 使用真实的 GitHub API 进行测试
}
```

### 11.5 测试覆盖率要求

| 模块 | 最低覆盖率 | 说明 |
|------|-----------|------|
| `internal/github` | 80% | 核心 API 交互逻辑 |
| `internal/converter` | 80% | Markdown 转换逻辑 |
| `internal/utils` | 90% | 工具函数（通常简单但关键） |
| `cmd/`, `web/` | 60% | 入口点和处理器 |

**运行覆盖率**：

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率
go tool cover -func=coverage.out

# HTML 报告
go tool cover -html=coverage.out
```

### 11.6 测试目录结构

```
internal/
  github/
    github.go          # 核心实现
    github_test.go     # 集成测试
    queries_test.go    # 单元测试
  converter/
    converter.go
    converter_test.go  # 表格驱动测试
```

### 11.7 运行测试

```bash
# 运行所有测试
make test

# 只运行单元测试（跳过集成测试）
go test -short ./...

# 运行特定包的测试
go test -v ./internal/github

# 运行特定测试
go test -v -run TestParseURL
```

---

## 12. 错误处理规范 (Error Handling)

**本章节遵循 `constitution.md` 第三条第 3.1 款：错误处理不可协商。**

### 12.1 错误包装规则

**强制要求**：所有错误必须使用 `fmt.Errorf("...: %w", err)` 进行包装。

```go
// ❌ 错误示范：直接返回
func FetchIssue(owner, repo string, number int) (*Issue, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err  // 失去上下文
    }
    // ...
}

// ✅ 正确示范：添加上下文后包装
func FetchIssue(owner, repo string, number int) (*Issue, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch issue %s/%s#%d: %w", owner, repo, number, err)
    }
    // ...
}
```

### 12.2 错误消息规范

错误消息必须包含足够的上下文信息：

| 场景 | 必需上下文 | 示例 |
|------|-----------|------|
| API 调用失败 | 操作名称、资源标识 | `failed to fetch issue golang/go#123` |
| 文件操作 | 文件路径、操作类型 | `failed to write output.md: permission denied` |
| 解析失败 | 输入内容、期望格式 | `invalid URL format: "not-a-url"` |

### 12.3 错误链传递

在多层调用中，每层都应添加其特定的上下文：

```go
func ConvertIssue(url string) error {
    parsed, err := ParseURL(url)
    if err != nil {
        return fmt.Errorf("parse URL: %w", err)  // 第一层包装
    }

    issue, err := FetchIssue(parsed.Owner, parsed.Repo, parsed.Number)
    if err != nil {
        return fmt.Errorf("fetch issue: %w", err)  // 第二层包装
    }

    return nil
}
```

---

## 13. 附录

### 13.1 GitHub API 相关

- Issue API: https://docs.github.com/en/rest/issues/issues
- PR API: https://docs.github.com/en/rest/pulls/pulls
- Discussion API: https://docs.github.com/en/graphql/guides/using-the-graphql-api-for-discussions

### 13.2 参考示例

输入：
```bash
export GITHUB_TOKEN=ghp_xxx
issue2md https://github.com/golang/go/issues/12345
```

输出文件：`12345-issue-runtime-panic-in-gc.md`

内容预览：
```markdown
---
title: "runtime: panic in GC"
type: issue
number: 12345
repository: golang/go
author: user123
created_at: 2024-01-01T00:00:00Z
state: closed
labels:
  - compiler
  - runtime
---

# runtime: panic in GC

## 描述

运行垃圾回收时发生 panic...

## 评论

### @user1 (2024-01-01 10:00)

这是一个已知问题...

> ### @user2 (2024-01-01 11:00)
>
> 我也遇到了这个问题...
```
