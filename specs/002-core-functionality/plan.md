# issue2md 技术实现方案

**文档版本**: 3.0
**创建日期**: 2026-04-01
**适用规格**: `specs/002-core-functionality/spec.md`
**遵循宪法**: `constitution.md`

---

## 1. 技术上下文总结

### 1.1 技术选型

| 技术领域 | 选型 | 理由 |
|---------|------|------|
| **编程语言** | Go >= 1.21.0 | 静态类型、强并发支持、跨平台编译 |
| **HTTP 客户端** | 标准库 `net/http` | 遵循"简单性原则"，避免过度依赖 |
| **GitHub API** | `google/go-github` v56 | 现有依赖，已验证稳定 |
| **配置管理** | 命令行 flags + 环境变量 | Unix 习惯，无配置文件 |
| **日志记录** | 标准库 `log` | 简单可靠 |
| **输出格式** | 仅 GitHub Flavored Markdown | 保持简单，专注核心功能 |

### 1.2 依赖管理

**现有依赖**：
```
require (
    github.com/google/go-github/v56 v56.0.0
    github.com/google/go-querystring v1.1.0 // indirect
)
```

**原则**：保持最小依赖，不引入非必需库。

### 1.3 性能目标

| 指标 | 目标值 | 实现策略 |
|------|--------|---------|
| 单次转换时间 | < 5 秒 | 并发获取评论，复用 HTTP 连接 |
| 内存占用 | < 50 MB | 流式处理，大 Diff 截断 |
| 启动时间 | < 100ms | 最小化启动逻辑 |

---

## 2. "合宪性"审查

### 2.1 第一条：简单性原则 (Simplicity First)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **1.1 YAGNI** | 只实现 spec.md 明确需求，移除 HTML/JSON 输出 | ✅ |
| **1.2 标准库优先** | HTTP/日志均用标准库，唯一外部依赖 `go-github/v56` | ✅ |
| **1.3 反过度工程** | 移除配置文件，只用 flags + 环境变量 | ✅ |

**简化决策**：
```go
// ❌ 移除：多格式输出
type Converter interface { Convert(data interface{}) (string, error) }

// ✅ 简化：直接函数
func IssueToMarkdown(issue *github.Issue, opts Options) (string, error)
```

### 2.2 第二条：测试先行铁律 (Test-First Imperative)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **2.1 TDD 循环** | 所有新功能从 `_test.go` 开始，Red-Green-Refactor | ✅ |
| **2.2 表格驱动** | 单元测试强制使用 `tests := []struct{...}{...}` 模式 | ✅ |
| **2.3 拒绝 Mocks** | 集成测试使用真实 GitHub API | ✅ |

**测试架构**：
```
internal/
  cli/
    cli.go
    cli_test.go         # 表格驱动测试
  config/
    config.go
    config_test.go      # 表格驱动测试
  converter/
    issue.go
    issue_test.go       # 表格驱动测试
    pullrequest.go
    pullrequest_test.go # 表格驱动测试
  github/
    client.go
    client_test.go      # 集成测试（真实 API）
  parser/
    parser.go
    parser_test.go      # 表格驱动测试
```

### 2.3 第三条：明确性原则 (Clarity and Explicitness)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **3.1 错误处理** | 所有错误使用 `fmt.Errorf("context: %w", err)` 包装 | ✅ |
| **3.2 无全局变量** | 所有依赖通过函数参数显式传递 | ✅ |

**错误处理示例**：
```go
func ParseURL(rawURL string) (*ParsedURL, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse URL %q: %w", rawURL, err)
    }
    // ...
}
```

---

## 3. 项目结构细化

### 3.1 目录树

```
issue2md/
├── cmd/
│   ├── issue2md/          # CLI 工具入口
│   │   └── main.go        # 调用 cli.Run()
│   └── issue2mdweb/       # Web 服务入口（Phase 4）
│       └── main.go
│
├── internal/              # 内部包
│   ├── cli/               # 命令行协调器
│   │   └── cli.go         # Run(), Execute()
│   ├── config/            # 配置管理
│   │   └── config.go      # ParseFlags(), Options
│   ├── converter/         # Markdown 转换
│   │   ├── issue.go       # IssueToMarkdown()
│   │   ├── pullrequest.go # PullRequestToMarkdown()
│   │   ├── discussion.go  # DiscussionToMarkdown()
│   │   └── formatter.go   # 辅助函数
│   ├── github/            # GitHub API 客户端
│   │   ├── client.go      # Client, Fetch*()
│   │   └── types.go       # Issue, PR, Discussion
│   └── parser/            # URL 解析
│       └── parser.go      # ParseURL(), ParsedURL
│
├── web/                   # Web 资源
│   ├── static/
│   └── templates/
│
├── specs/                 # 功能规格
├── .claude/              # Claude 配置
├── Makefile              # 构建脚本
├── go.mod                # Go 模块
├── README.md             # 项目文档
└── constitution.md       # 项目宪法
```

### 3.2 包职责矩阵

| 包 | 职责 | 依赖 | 导出 |
|---|------|------|------|
| `internal/parser` | URL 解析、类型识别 | 标准库 | `ParseURL()`, `ParsedURL` |
| `internal/config` | flags 解析、选项生成 | 标准库 | `ParseFlags()`, `Options` |
| `internal/github` | GitHub API 调用 | `parser`, `go-github` | `Client`, `Fetch*()` |
| `internal/converter` | 数据转 Markdown | `github`, `config` | `IssueToMarkdown()` 等 |
| `internal/cli` | CLI 入口协调 | 所有内部包 | `Run()`, `Execute()` |
| `cmd/issue2md` | 程序入口 | `internal/cli` | `main()` |

### 3.3 依赖关系图

```
                    ┌─────────────┐
                    │   cli       │
                    │  (入口)     │
                    └──────┬──────┘
                           │
        ┌──────────────────┼──────────────────┐
        ▼                  ▼                  ▼
   ┌─────────┐       ┌─────────┐       ┌─────────┐
   │ config  │       │ parser  │       │ github  │
   └─────────┘       └────┬────┘       └────┬────┘
                           │                  │
                           └────────┬─────────┘
                                    ▼
                             ┌─────────┐
                             │converter│
                             └─────────┘
```

**依赖原则**：
- 单向依赖：`cmd` → `internal` → 外部库
- `cli` 协调所有包，但其他包互不依赖
- `converter` 聚合 `github` 和 `config` 的结果

---

## 4. 核心数据结构

### 4.1 parser 包数据结构

```go
// internal/parser/parser.go

// ResourceType GitHub 资源类型
type ResourceType string

const (
    TypeIssue       ResourceType = "issue"
    TypePullRequest ResourceType = "pull_request"
    TypeDiscussion  ResourceType = "discussion"
)

// ParsedURL 解析后的 GitHub URL
type ParsedURL struct {
    Owner  string
    Repo   string
    Number int
    Type   ResourceType
}
```

### 4.2 config 包数据结构

```go
// internal/config/config.go

// Config 应用配置（来自 flags）
type Config struct {
    EnableReactions  bool  // -enable-reactions
    EnableUserLinks  bool  // -enable-user-links
    OutputFile       string // [output_file] 位置参数
    versionRequested bool   // -v, -version
    helpRequested    bool   // -h, -help
}

// Options 转换器选项
type Options struct {
    EnableReactions bool
    EnableUserLinks bool
}
```

### 4.3 github 包数据结构

```go
// internal/github/types.go

// Issue GitHub Issue 数据
type Issue struct {
    Number      int
    Title       string
    Body        string
    State       string
    Author      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Labels      []Label
    Milestone   *Milestone
    Reactions   Reactions
    Comments    []Comment
}

// PullRequest GitHub PR 数据
type PullRequest struct {
    Issue
    Mergeable    bool
    Merged       bool
    MergedAt     *time.Time
    Additions    int
    Deletions    int
    ChangedFiles int
    HeadBranch   string
    BaseBranch   string
    Files        []File
}

// Discussion GitHub Discussion 数据
type Discussion struct {
    Number    int
    Title     string
    Body      string
    Author    string
    CreatedAt time.Time
    Category  string
    Upvotes   int
    Replies   []DiscussionReply
}

// Label, Milestone, Reactions, Comment, File, DiscussionReply...
```

---

## 5. 接口设计

### 5.1 parser 包接口

```go
// ParseURL 解析 GitHub URL
func ParseURL(rawURL string) (*ParsedURL, error)

// String 返回字符串表示
func (p *ParsedURL) String() string
```

### 5.2 config 包接口

```go
// ParseFlags 解析命令行参数
func ParseFlags(args []string) (*Config, error)

// Validate 验证配置
func (c *Config) Validate() error

// OutputOptions 获取转换选项
func (c *Config) OutputOptions() Options

// VersionInfo 返回版本信息
func VersionInfo() string
```

### 5.3 github 包接口

```go
// Client GitHub API 客户端
type Client struct {
    token string
}

// NewClient 创建客户端
func NewClient(token string) *Client

// FetchIssue 获取 Issue
func (c *Client) FetchIssue(ctx context.Context, parsed *ParsedURL) (*Issue, error)

// FetchPullRequest 获取 PR
func (c *Client) FetchPullRequest(ctx context.Context, parsed *ParsedURL) (*PullRequest, error)

// FetchDiscussion 获取 Discussion（Phase 3）
func (c *Client) FetchDiscussion(ctx context.Context, parsed *ParsedURL) (*Discussion, error)
```

### 5.4 converter 包接口

```go
// IssueToMarkdown Issue 转 Markdown
func IssueToMarkdown(issue *github.Issue, opts config.Options) (string, error)

// PullRequestToMarkdown PR 转 Markdown
func PullRequestToMarkdown(pr *github.PullRequest, opts config.Options) (string, error)

// DiscussionToMarkdown Discussion 转 Markdown
func DiscussionToMarkdown(disc *github.Discussion, opts config.Options) (string, error)
```

### 5.5 cli 包接口

```go
// Run 执行 CLI 主逻辑
func Run() int

// Execute 运行命令（用于测试）
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) int
```

---

## 6. CLI 执行流程

```
┌─────────────────────────────────────────────────────────────────┐
│  issue2md [flags] <github_url> [output_file]                    │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
              ┌─────────────────────┐
              │  1. ParseFlags()    │
              └─────────┬───────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
   helpRequested? versionRequested?  继续执行
        │               │
        ▼               ▼
   printHelp()     printVersion()
        │               │
        └───────┬───────┘
                ▼
            exit 0
                │
                ▼ (未请求 help/version)
      ┌─────────────────────┐
      │  2. Validate()      │
      └─────────┬───────────┘
                │
        ┌───────┴───────┐
        ▼               ▼
     无 URL?          有 URL
        │               │
        ▼               ▼
     error       ┌─────────────────┐
        │       │ 3. ParseURL()   │
        │       └────────┬────────┘
        │                │
        │        ┌───────┴───────┐
        │        ▼               ▼
        │     解析失败         解析成功
        │        │               │
        │        ▼               ▼
        │       error     ┌──────────────────┐
        │                 │ 4. NewClient()   │
        │                 └────────┬─────────┘
        │                          │
        │                 ┌────────┴────────┐
        │                 ▼                 ▼
        │            Token 空?        Token 有效
        │                 │                 │
        │                 ▼                 ▼
        │                error      ┌─────────────────┐
        │                            │ 5. Fetch*()     │
        │                            └────────┬────────┘
        │                                     │
        │                          ┌──────────┴──────────┐
        │                          ▼                     ▼
        │                     API 成功              API 失败
        │                          │                     │
        │                          ▼                     ▼
        │                   ┌─────────────┐          error
        │                   │ 6. ToMarkdown()        │
        │                   └──────┬──────┘          │
        │                          │                 │
        │                   ┌──────┴──────┐          │
        │                   ▼             ▼          │
        │              转换成功      转换失败          │
        │                   │             │          │
        │                   ▼             ▼          │
        │            ┌────────────────┐  error        │
        │            │ 7. Write output│             │
        │            └────────┬───────┘             │
        │                     │                      │
        │              ┌──────┴──────┐              │
        │              ▼             ▼              │
        │          写入成功      写入失败             │
        │              │             │              │
        ▼              ▼             ▼              ▼
    exit 1       exit 0        exit 5          exit 2/3/4/6
```

---

## 7. 错误处理与退出码

### 7.1 退出码常量

```go
const (
    ExitSuccess       = 0  // 正常 / help / version
    ExitInvalidURL    = 1  // URL 格式无效
    ExitNotFound      = 2  // 资源不存在
    ExitAPIFailed     = 3  // API 请求失败
    ExitTokenMissing  = 4  // GITHUB_TOKEN 未设置
    ExitWriteFailed   = 5  // 文件写入失败
    ExitTimeout       = 6  // 网络超时
)
```

### 7.2 错误处理规范

```go
// 所有错误必须包装
func ParseURL(rawURL string) (*ParsedURL, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse URL %q: %w", rawURL, err)
    }
    // ...
}
```

---

## 8. 实现路线图

### Phase 1: MVP（Issue 支持）

| 任务 | 优先级 | 验收标准 |
|------|--------|---------|
| `internal/parser` | P0 | 表格驱动测试覆盖所有 URL 格式 |
| `internal/config` | P0 | `-h`, `-v`, flags 解析正确 |
| `internal/github` | P0 | 集成测试通过真实 API |
| `internal/converter` | P0 | Issue 转 Markdown 符合 spec |
| `internal/cli` | P0 | 完整流程可运行 |
| `cmd/issue2md` | P0 | 可执行文件构建成功 |

### Phase 2: PR 支持

| 任务 | 优先级 | 验收标准 |
|------|--------|---------|
| `internal/github` - PR | P0 | 获取 PR + Files + Diff |
| `internal/converter` - PR | P0 | PR 转 Markdown，Diff 折叠 |

### Phase 3: Discussion 支持

| 任务 | 优先级 | 验收标准 |
|------|--------|---------|
| `internal/github` - GraphQL | P1 | GraphQL 基础设施 |
| `internal/github` - Discussion | P1 | 获取 Discussion |
| `internal/converter` - Discussion | P1 | 嵌套回复缩进正确 |

### Phase 4: Web 服务

| 任务 | 优先级 | 验收标准 |
|------|--------|---------|
| `cmd/issue2mdweb` | P2 | HTTP 服务器启动 |
| GET /health | P2 | 返回健康状态 |
| POST /convert | P2 | 接收 URL，返回 Markdown |

---

## 9. 测试策略

### 9.1 表格驱动测试示例

```go
func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    *ParsedURL
        wantErr bool
    }{
        {"valid issue", "https://github.com/owner/repo/issues/123", &ParsedURL{...}, false},
        {"invalid", "not-a-url", nil, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
            }
            // ...
        })
    }
}
```

### 9.2 集成测试

```go
func TestFetchIssue_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        t.Skip("GITHUB_TOKEN not set")
    }
    // 测试真实 API
}
```

---

## 10. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| GitHub API 限流 | 获取失败 | 指数退避重试 |
| 大 PR 内存占用 | OOM | Diff 超过 500 行截断 |
| 并发安全 | 数据竞态 | 无状态设计 |

---

## 11. 附录

### 11.1 CLI 参数

| 参数 | 说明 |
|------|------|
| `-h`, `-help` | 显示帮助 |
| `-v`, `-version` | 显示版本 |
| `-enable-reactions` | 包含反应统计 |
| `-enable-user-links` | 用户名转链接 |

### 11.2 环境变量

| 变量 | 说明 |
|------|------|
| `GITHUB_TOKEN` | GitHub API Token（必需） |
