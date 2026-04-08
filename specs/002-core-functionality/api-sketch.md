# issue2md API 设计草图

**文档版本**: 2.0
**创建日期**: 2026-04-01
**基于规格**: `specs/002-core-functionality/spec.md`
**遵循宪法**: `constitution.md`

---

## 1. 包结构总览

```
issue2md/
├── cmd/
│   ├── issue2md/          # CLI 入口
│   └── issue2mdweb/       # Web 入口
│
├── internal/
│   ├── cli/               # 命令行参数解析与执行
│   ├── config/            # 配置管理（flags, 环境变量）
│   ├── parser/            # URL 解析与类型识别
│   ├── github/            # GitHub API 交互
│   └── converter/         # 数据转换为 Markdown
│
└── web/
    ├── templates/         # HTML 模板
    └── static/            # 静态资源
```

### 1.1 包职责与依赖

| 包 | 职责 | 依赖 | 被依赖 |
|---|------|------|--------|
| `internal/parser` | URL 解析、类型识别 | 标准库 | `cli`, `github` |
| `internal/config` | 配置加载、验证 | 标准库 | `cli`, `converter` |
| `internal/github` | GitHub API 调用 | `parser`, `google/go-github` | `cli`, `converter` |
| `internal/converter` | 数据转 Markdown | `github`, `config` | `cli` |
| `internal/cli` | CLI 入口协调 | 所有内部包 | - |

### 1.2 依赖关系图

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

**依赖原则（constitution.md 第一条）**：
- 单向依赖，禁止循环
- `cli` 协调所有包，但其他包互不依赖
- `converter` 聚合 `github` 的输出

---

## 2. internal/parser 包

**职责**: 解析 GitHub URL，识别资源类型，提取参数

### 2.1 数据结构

```go
// parser/parser.go

// ResourceType GitHub 资源类型
type ResourceType string

const (
    TypeIssue       ResourceType = "issue"
    TypePullRequest ResourceType = "pull_request"
    TypeDiscussion  ResourceType = "discussion"
)

// ParsedURL 解析后的 GitHub URL
type ParsedURL struct {
    Owner  string       // 仓库所有者
    Repo   string       // 仓库名称
    Number int          // Issue/PR/Discussion 编号
    Type   ResourceType // 资源类型
}
```

### 2.2 导出接口

```go
// ParseURL 解析 GitHub URL
// 支持的格式:
//   - https://github.com/owner/repo/issues/123
//   - https://github.com/owner/repo/pull/456
//   - https://github.com/owner/repo/discussions/789
//   - 带 .git 后缀、www 子域、http 协议等变体
//
// 返回 *ParsedURL 或错误（格式无效、不支持的类型）
func ParseURL(rawURL string) (*ParsedURL, error)

// String 返回 ParsedURL 的字符串表示（用于测试）
func (p *ParsedURL) String() string
```

### 2.3 错误定义

```go
var (
    ErrInvalidURLFormat   = errors.New("invalid URL format")
    ErrUnsupportedType    = errors.New("unsupported resource type")
    ErrMissingOwner       = errors.New("missing repository owner")
    ErrMissingRepo        = errors.New("missing repository name")
    ErrMissingNumber      = errors.New("missing issue/PR number")
)
```

---

## 3. internal/config 包

**职责**: 管理 flags、环境变量，提供转换选项

### 3.1 数据结构

```go
// config/config.go

// Config 应用配置
type Config struct {
    // Flags
    EnableReactions  bool  // -enable-reactions
    EnableUserLinks  bool  // -enable-user-links

    // 输出
    OutputFile       string // 位置参数 [output_file]，空表示 stdout

    // 内部使用
    versionRequested bool   // -v 或 -version
    helpRequested    bool   // -h 或 -help
}

// OutputOptions 转换器选项（从 Config 派生）
type OutputOptions struct {
    EnableReactions  bool
    EnableUserLinks  bool
}
```

### 3.2 导出接口

```go
// ParseFlags 解析命令行参数
// 参数: os.Args[1:]
// 返回 *Config 或错误
func ParseFlags(args []string) (*Config, error)

// Validate 验证配置有效性
func (c *Config) Validate() error

// OutputOptions 获取转换器选项
func (c *Config) OutputOptions() OutputOptions

// VersionInfo 返回版本信息字符串
func VersionInfo() string
```

### 3.3 Flag 定义

```go
const (
    FlagHelp            = "help"
    FlagHelpShort       = "h"
    FlagVersion         = "version"
    FlagVersionShort    = "v"
    FlagEnableReactions = "enable-reactions"
    FlagEnableUserLinks = "enable-user-links"
)
```

---

## 4. internal/github 包

**职责**: 与 GitHub API 交互，获取 Issue/PR/Discussion 数据

### 4.1 数据结构

```go
// github/types.go

// Issue GitHub Issue 数据
type Issue struct {
    Number      int
    Title       string
    Body        string
    State       string // "open", "closed"
    Author      string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Labels      []Label
    Milestone   *Milestone
    Reactions   Reactions
    Comments    []Comment
}

// PullRequest GitHub PR 数据（继承 Issue）
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
    UpdatedAt time.Time
    Category  string
    Upvotes   int
    Replies   []DiscussionReply
}

// Label Issue 标签
type Label struct {
    Name  string
    Color string
}

// Milestone 里程碑
type Milestone struct {
    Title string
    State string
}

// Reactions 反应统计
type Reactions struct {
    ThumbsUp   int
    ThumbsDown int
    Laugh      int
    Hooray     int
    Confused   int
    Heart      int
    Rocket     int
    Eyes       int
}

// Comment Issue/PR 评论
type Comment struct {
    ID        int64
    Author    string
    Body      string
    CreatedAt time.Time
    Reactions Reactions
}

// File PR 变更文件
type File struct {
    Path      string
    Additions int
    Deletions int
    Patch     string // uni-diff
    BlobURL   string // GitHub 链接
}

// DiscussionReply Discussion 回复
type DiscussionReply struct {
    ID        string
    Author    string
    Body      string
    CreatedAt time.Time
    Replies   []DiscussionReply // 嵌套
}
```

### 4.2 导出接口

```go
// github/client.go

// Client GitHub API 客户端
type Client struct {
    token string
    // 内部使用 *github.Client，但不暴露
}

// NewClient 创建 GitHub 客户端
// token: GITHUB_TOKEN 环境变量的值
func NewClient(token string) *Client

// FetchIssue 获取 Issue 完整数据
// 包含: Issue 信息 + Labels + Milestone + Reactions + 所有评论
func (c *Client) FetchIssue(ctx context.Context, parsed *parser.ParsedURL) (*Issue, error)

// FetchPullRequest 获取 PR 完整数据
// 包含: Issue 信息 + PR 特有字段 + Files + Patch（可能截断）+ 评论
func (c *Client) FetchPullRequest(ctx context.Context, parsed *parser.ParsedURL) (*PullRequest, error)

// FetchDiscussion 获取 Discussion 完整数据
// 注意: 需要 GraphQL API，Phase 3 实现
func (c *Client) FetchDiscussion(ctx context.Context, parsed *parser.ParsedURL) (*Discussion, error)
```

### 4.3 错误定义

```go
var (
    ErrTokenRequired     = errors.New("GITHUB_TOKEN required")
    ErrResourceNotFound  = errors.New("resource not found")
    ErrAPIRequestFailed  = errors.New("API request failed")
    ErrNetworkTimeout    = errors.New("network timeout")
)
```

---

## 5. internal/converter 包

**职责**: 将 GitHub 数据转换为 GitHub Flavored Markdown

### 5.1 数据结构

```go
// converter/options.go

// Options 转换选项（从 config.OutputOptions 映射）
type Options struct {
    EnableReactions bool
    EnableUserLinks bool
}
```

### 5.2 导出接口

```go
// converter/issue.go

// IssueToMarkdown 将 Issue 转换为 Markdown
// 输出格式符合 spec.md 第 2.2 节要求:
//   - YAML Front Matter
//   - 标题 (H1)
//   - 描述内容
//   - 标签/里程碑元数据
//   - 评论区（时间顺序，嵌套用引用块）
func IssueToMarkdown(issue *github.Issue, opts Options) (string, error)

// converter/pullrequest.go

// PullRequestToMarkdown 将 PR 转换为 Markdown
// 额外包含 spec.md 第 6 节要求:
//   - 变更文件列表
//   - 折叠的 Diff (<details>)
//   - 大 PR 截断处理（>500 行）
func PullRequestToMarkdown(pr *github.PullRequest, opts Options) (string, error)

// converter/discussion.go

// DiscussionToMarkdown 将 Discussion 转换为 Markdown
// 处理嵌套回复的缩进和引用块
func DiscussionToMarkdown(discussion *github.Discussion, opts Options) (string, error)
```

### 5.3 辅助函数（内部使用）

```go
// converter/formatter.go

// writeYAMLFrontMatter 生成 YAML 头
func writeYAMLFrontMatter(w io.Writer, data interface{}) error

// writeComments 格式化评论列表
func writeComments(w io.Writer, comments []github.Comment, opts Options) error

// writeReactions 格式化反应统计
func writeReactions(w io.Writer, reactions github.Reactions) error

// writeUserLink 格式化用户链接（条件性）
func writeUserLink(w io.Writer, username string, enabled bool) string

// sanitizeFilename 清理文件名（spec.md 第 5.2 节）
func sanitizeFilename(title string) string
```

### 5.4 输出格式示例

```go
// IssueToMarkdown 输出示例:
/*
---
title: "Bug: 修复登录失败"
type: issue
number: 123
repository: owner/repo
author: username
created_at: 2024-01-01T00:00:00Z
state: closed
labels:
  - bug
reactions:
  thumbs_up: 5
  heart: 3
---

# Bug: 修复登录失败

## 描述

登录时出现 panic...

## 评论

### @user1 (2024-01-01 10:00)

这是一个已知问题...

> ### @user2 (2024-01-01 11:00)
>
> 我也遇到了这个问题...

👍 5  ❤️ 3
*/
```

---

## 6. internal/cli 包

**职责**: CLI 入口，协调所有包

### 6.1 导出接口

```go
// cli/cli.go

// Run 执行 CLI 主逻辑
// 返回退出码（0-6，见 spec.md 第 7.1 节）
func Run() int

// Execute 运行命令（用于测试）
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) int
```

### 6.2 执行流程

```
1. ParseFlags(args)
   │
   ├─► helpRequested?    ──► printHelp()    ──► exit 0
   ├─► versionRequested? ──► printVersion() ──► exit 0
   │
2. Validate()
   │
   ├─► no URL?          ──► error          ──► exit 1
   │
3. ParseURL(githubURL)
   │
   ├─► error?           ──► error          ──► exit 1
   │
4. NewClient(GITHUB_TOKEN)
   │
   ├─► token empty?     ──► error          ──► exit 4
   │
5. FetchIssue/PR/Discussion()
   │
   ├─► error?           ──► error          ──► exit 2/3/5/6
   │
6. ToMarkdown(data, opts)
   │
7. Write to stdout or file
   │
   ├─► write error?     ──► error          ──► exit 5
   │
8. exit 0
```

---

## 7. 错误处理规范

**遵循 constitution.md 第三条**

### 7.1 错误包装

```go
// 所有错误必须使用 %w 包装
func ParseURL(rawURL string) (*ParsedURL, error) {
    u, err := url.Parse(rawURL)
    if err != nil {
        return nil, fmt.Errorf("failed to parse URL %q: %w", rawURL, err)
    }
    // ...
}
```

### 7.2 退出码映射

```go
const (
    ExitSuccess          = 0  // 正常退出 / help / version
    ExitInvalidURL       = 1  // URL 格式无效
    ExitNotFound         = 2  // 资源不存在
    ExitAPIFailed        = 3  // API 请求失败
    ExitTokenMissing     = 4  // GITHUB_TOKEN 未设置
    ExitWriteFailed      = 5  // 文件写入失败
    ExitTimeout          = 6  // 网络超时
)
```

---

## 8. 测试策略

**遵循 constitution.md 第二条**

### 8.1 表格驱动测试示例

```go
// parser/parser_test.go

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
            want:    &ParsedURL{Owner: "owner", Repo: "repo", Number: 123, Type: TypeIssue},
            wantErr: false,
        },
        {
            name:    "invalid URL",
            url:     "not-a-url",
            want:    nil,
            wantErr: true,
        },
        // 更多测试用例...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 8.2 集成测试

```go
// github/github_test.go

func TestFetchIssue_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        t.Skip("GITHUB_TOKEN not set")
    }

    client := NewClient(token)
    parsed, _ := parser.ParseURL("https://github.com/golang/go/issues/12345")

    issue, err := client.FetchIssue(context.Background(), parsed)
    if err != nil {
        t.Fatalf("FetchIssue() failed: %v", err)
    }

    if issue.Number != 12345 {
        t.Errorf("got number %d, want 12345", issue.Number)
    }
}
```

---

## 9. 实现顺序

### Phase 1: Issue 支持

1. `internal/parser` - URL 解析（可独立测试）
2. `internal/config` - 参数解析（可独立测试）
3. `internal/github` - Issue 获取（集成测试）
4. `internal/converter` - Issue 转 Markdown
5. `internal/cli` - 组装

### Phase 2: PR 支持

1. `internal/github` - PR 获取
2. `internal/converter` - PR 转 Markdown（含 Diff）

### Phase 3: Discussion 支持

1. `internal/github` - GraphQL 基础设施
2. `internal/github` - Discussion 获取
3. `internal/converter` - Discussion 转 Markdown
