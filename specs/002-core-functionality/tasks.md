# issue2md 任务列表

**生成日期**: 2026-04-01
**基于规格**: `specs/002-core-functionality/spec.md`
**技术方案**: `specs/002-core-functionality/plan.md`
**遵循宪法**: `constitution.md`

---

## 任务说明

- **[P]** 标记表示该任务与其他任务无依赖关系，可并行执行
- **TDD 强制**：测试任务必须先于实现任务执行
- **依赖关系**：任务按顺序排列，后面的任务依赖前面的任务完成

---

## Phase 1: Foundation (数据结构定义)

> 目标：定义所有核心数据结构，建立项目骨架

### 1.1 [P] 创建 parser 包基础结构

**文件**: `internal/parser/parser.go`

**任务**:
- 定义 `ResourceType` 类型及常量
- 定义 `ParsedURL` 结构体
- 添加文档注释

**验收**:
- [ ] `ResourceType` 类型定义，包含 `TypeIssue`, `TypePullRequest`, `TypeDiscussion`
- [ ] `ParsedURL` 结构体包含 `Owner`, `Repo`, `Number`, `Type` 字段
- [ ] 所有公开类型有 godoc 注释

---

### 1.2 [P] 创建 config 包基础结构

**文件**: `internal/config/config.go`

**任务**:
- 定义 `Config` 结构体
- 定义 `Options` 结构体
- 添加文档注释

**验收**:
- [ ] `Config` 包含 `EnableReactions`, `EnableUserLinks`, `OutputFile`, `versionRequested`, `helpRequested`
- [ ] `Options` 包含 `EnableReactions`, `EnableUserLinks`
- [ ] 所有公开类型有 godoc 注释

---

### 1.3 [P] 创建 github 包数据结构

**文件**: `internal/github/types.go`

**任务**:
- 定义 `Issue` 结构体
- 定义 `PullRequest` 结构体
- 定义 `Discussion` 结构体
- 定义辅助类型：`Label`, `Milestone`, `Reactions`, `Comment`, `File`, `DiscussionReply`

**验收**:
- [ ] 所有结构体字段完整
- [ ] `PullRequest` 嵌入 `Issue`
- [ ] 所有公开类型有 godoc 注释

---

### 1.4 定义退出码常量

**文件**: `internal/cli/exit_codes.go`

**任务**:
- 定义退出码常量（0-6）
- 添加文档注释

**验收**:
- [ ] `ExitSuccess = 0`, `ExitInvalidURL = 1`, ..., `ExitTimeout = 6`
- [ ] 每个常量有注释说明使用场景

---

## Phase 2: GitHub Fetcher (API 交互逻辑)

> 目标：实现 URL 解析和 GitHub API 交互

### 2.1 [P] 创建 parser 包测试文件

**文件**: `internal/parser/parser_test.go`

**任务**:
- 编写 `TestParseURL` 表格驱动测试
- 覆盖所有 URL 格式变体（spec.md 2.1.1 节）
- 包含错误用例

**测试用例**:
- [ ] 标准 Issue URL
- [ ] 标准 PR URL
- [ ] 标准 Discussion URL
- [ ] 带 `.git` 后缀
- [ ] 带 www 子域
- [ ] HTTP 协议
- [ ] 无效 URL 格式
- [ ] 缺少 owner/repo
- [ ] 缺少编号

**验收**: 所有测试失败（Red 阶段）

---

### 2.2 [P] 实现 ParseURL 函数

**文件**: `internal/parser/parser.go`

**任务**:
- 实现 `ParseURL(rawURL string) (*ParsedURL, error)` 函数
- 实现 `(*ParsedURL) String() string` 方法
- 错误处理使用 `fmt.Errorf` 包装

**验收**:
- [ ] `go test ./internal/parser -v` 全部通过
- [ ] URL 解析逻辑符合 spec.md 2.1.1 节
- [ ] 错误消息清晰

---

### 2.3 创建 github 包集成测试框架

**文件**: `internal/github/client_test.go`

**任务**:
- 创建测试文件结构
- 创建 `TestFetchIssue_Integration` 测试框架
- 添加 `testing.Short()` 跳过逻辑
- 添加 `GITHUB_TOKEN` 环境变量检查

**验收**:
- [ ] 测试文件编译通过
- [ ] 跳过逻辑工作正常

---

### 2.4 创建 github 包客户端基础

**文件**: `internal/github/client.go`

**任务**:
- 定义 `Client` 结构体
- 实现 `NewClient(token string) *Client` 函数
- 添加 `FetchIssue()` 和 `FetchPullRequest()` 空函数签名
- 添加文档注释

**验收**:
- [ ] 代码编译通过
- [ ] 函数签名正确
- [ ] 所有测试失败（Red 阶段）

---

### 2.5 实现 FetchIssue 函数（测试）

**文件**: `internal/github/client_test.go`

**任务**:
- 完善 `TestFetchIssue_Integration` 测试用例
- 覆盖基本 Issue 获取
- 覆盖 Issue 带评论
- 覆盖 Issue 带 Labels
- 添加 `TestFetchPullRequest_Integration` 测试用例

**验收**:
- [ ] 所有测试用例编写完成
- [ ] 运行测试全部失败（函数未实现或返回 nil）

---

### 2.6 实现 FetchIssue 函数

**文件**: `internal/github/client.go`

**任务**:
- 实现 `FetchIssue(ctx, owner, repo, number)` 函数
- 使用 `google/go-github` 库获取 Issue
- 转换数据结构（Issue, Labels, Milestone, Reactions, Comments）
- 实现评论分页获取
- 所有错误使用 `fmt.Errorf("context: %w")` 包装

**验收**:
- [ ] `go test ./internal/github -run TestFetchIssue_Integration` 通过
- [ ] Issue 数据完整正确
- [ ] 评论获取正确（支持分页）

---

### 2.7 实现 FetchPullRequest 函数

**文件**: `internal/github/client.go`

**任务**:
- 实现 `FetchPullRequest(ctx, owner, repo, number)` 函数
- 获取 PR 信息 + Files + Patch
- 复用 Issue 转换逻辑
- 注意：go-github v56 的 PullRequest 没有 Reactions 字段

**验收**:
- [ ] `go test ./internal/github -run TestFetchPullRequest_Integration` 通过
- [ ] PR 数据完整正确
- [ ] Files 列表正确
- [ ] 验证 Number, Title, State, Author 等字段
- [ ] 测试失败时能正确跳过

---

### 2.6 实现 FetchIssue 函数

**文件**: `internal/github/client.go`

**任务**:
- 实现 `FetchIssue(ctx context.Context, parsed *ParsedURL) (*Issue, error)` 函数
- 使用 `google/go-github` 库
- 获取 Issue 信息 + Labels + Milestone + Reactions + Comments
- 所有错误使用 `%w` 包装

**验收**:
- [ ] `go test ./internal/github -v -run TestFetchIssue_Integration` 通过
- [ ] 能获取完整 Issue 数据
- [ ] 错误处理符合规范

---

### 2.7 实现 FetchPullRequest 函数（测试）

**文件**: `internal/github/client_test.go`

**任务**:
- 创建 `TestFetchPullRequest_Integration` 测试
- 使用真实 GitHub PR 测试

**验收**:
- [ ] 测试框架正确
- [ ] 能获取真实 PR 数据

---

### 2.8 实现 FetchPullRequest 函数

**文件**: `internal/github/client.go`

**任务**:
- 实现 `FetchPullRequest(ctx context.Context, parsed *ParsedURL) (*PullRequest, error)` 函数
- 获取 PR 信息 + Files + Patch
- Patch 超过 500 行时截断（spec.md 6.2 节）

**验收**:
- [ ] `go test ./internal/github -v -run TestFetchPullRequest_Integration` 通过
- [ ] 包含 Files 和 Patch 数据
- [ ] 大 PR 截断逻辑正确

---

## Phase 3: Markdown Converter (转换逻辑)

> 目标：实现 GitHub 数据到 Markdown 的转换

### 3.1 [P] 创建 converter 包 Issue 测试

**文件**: `internal/converter/issue_test.go`

**任务**:
- 编写 `TestIssueToMarkdown` 表格驱动测试
- 测试 YAML Front Matter 生成
- 测试评论格式化
- 测试 Reactions 开关
- 测试用户链接开关

**测试用例**:
- [ ] 基础 Issue 转换
- [ ] 包含 Reactions
- [ ] 包含用户链接
- [ ] 空评论列表
- [ ] 多条评论

**验收**: 所有测试失败（Red 阶段）

---

### 3.2 [P] 创建 converter 包辅助函数

**文件**: `internal/converter/formatter.go`

**任务**:
- 实现 `writeYAMLFrontMatter` 函数
- 实现 `writeComments` 函数
- 实现 `writeReactions` 函数
- 实现 `writeUserLink` 函数

**验收**:
- [ ] YAML Front Matter 格式符合 spec.md 2.2.1 节
- [ ] 评论按时间顺序排列
- [ ] Reactions 格式正确

---

### 3.3 实现 IssueToMarkdown 函数

**文件**: `internal/converter/issue.go`

**任务**:
- 实现 `IssueToMarkdown(issue *github.Issue, opts config.Options) (string, error)` 函数
- 组织 YAML Front Matter + 标题 + 描述 + 评论

**验收**:
- [ ] `go test ./internal/converter -v -run TestIssueToMarkdown` 通过
- [ ] 输出符合 spec.md 示例格式
- [ ] `EnableReactions` 开关生效
- [ ] `EnableUserLinks` 开关生效

---

### 3.4 [P] 创建 converter 包 PR 测试

**文件**: `internal/converter/pullrequest_test.go`

**任务**:
- 编写 `TestPullRequestToMarkdown` 表格驱动测试
- 测试变更文件列表
- 测试 Diff 折叠
- 测试大 PR 截断

**验收**: 所有测试失败（Red 阶段）

---

### 3.5 实现 PullRequestToMarkdown 函数

**文件**: `internal/converter/pullrequest.go`

**任务**:
- 实现 `PullRequestToMarkdown(pr *github.PullRequest, opts config.Options) (string, error)` 函数
- 复用 IssueToMarkdown 逻辑
- 添加变更文件列表
- 添加折叠的 Diff

**验收**:
- [ ] `go test ./internal/converter -v -run TestPullRequestToMarkdown` 通过
- [ ] Diff 使用 `<details>` 折叠
- [ ] 大 PR 截断逻辑正确

---

### 3.6 [P] 创建 converter 包 Discussion 测试

**文件**: `internal/converter/discussion_test.go`

**任务**:
- 编写 `TestDiscussionToMarkdown` 表格驱动测试
- 测试嵌套回复缩进

**验收**: 所有测试失败（Red 阶段）

---

### 3.7 实现 DiscussionToMarkdown 函数

**文件**: `internal/converter/discussion.go`

**任务**:
- 实现 `DiscussionToMarkdown(disc *github.Discussion, opts config.Options) (string, error)` 函数
- 处理嵌套回复的引用块缩进

**验收**:
- [ ] `go test ./internal/converter -v -run TestDiscussionToMarkdown` 通过
- [ ] 嵌套回复缩进正确

---

## Phase 4: CLI Assembly (命令行入口集成)

> 目标：组装 CLI 工具，实现完整功能

### 4.1 [P] 创建 config 包测试

**文件**: `internal/config/config_test.go`

**任务**:
- 编写 `TestParseFlags` 表格驱动测试
- 测试所有 flags 组合
- 测试位置参数解析

**测试用例**:
- [ ] 空参数
- [ ] `-h` / `-help`
- [ ] `-v` / `-version`
- [ ] `-enable-reactions`
- [ ] `-enable-user-links`
- [ ] 组合 flags
- [ ] 带 URL
- [ ] 带 URL + 输出文件

**验收**: 所有测试失败（Red 阶段）

---

### 4.2 [P] 实现 ParseFlags 函数

**文件**: `internal/config/config.go`

**任务**:
- 实现 `ParseFlags(args []string) (*Config, error)` 函数
- 手动解析命令行参数（不使用 flag 库，保持简单）
- 实现 `Validate() error` 方法
- 实现 `OutputOptions() Options` 方法
- 实现 `VersionInfo() string` 函数

**验收**:
- [ ] `go test ./internal/config -v -run TestParseFlags` 通过
- [ ] 所有 flags 正确解析
- [ ] 位置参数正确处理
- [ ] Validate 正确检测无 URL 情况

---

### 4.3 [P] 创建 cli 包测试框架

**文件**: `internal/cli/cli_test.go`

**任务**:
- 创建 `TestExecute` 表格驱动测试框架
- 测试各种输入场景

**测试用例**:
- [ ] `-h` 输出帮助并返回 0
- [ ] `-v` 输出版本并返回 0
- [ ] 无参数返回错误码 1
- [ ] 无 GITHUB_TOKEN 返回错误码 4
- [ ] 无效 URL 返回错误码 1
- [ ] 成功转换返回 0

**验收**: 测试框架创建完成

---

### 4.4 [P] 创建 cli 包基础结构

**文件**: `internal/cli/cli.go`

**任务**:
- 定义 `Execute(stdin, stdout, stderr io.Writer, args []string) int` 函数签名
- 定义 `Run() int` 函数
- 添加帮助信息常量
- 添加版本信息常量

**验收**:
- [ ] 函数签名正确
- [ ] 帮助信息符合 spec.md 3.6 节格式

---

### 4.5 实现 Execute 函数（基础流程）

**文件**: `internal/cli/cli.go`

**任务**:
- 实现 Execute 函数的基础流程
- 步骤 1: ParseFlags
- 步骤 2: 处理 help/version
- 步骤 3: Validate 检查 URL

**验收**:
- [ ] `-h` 正常输出帮助
- [ ] `-v` 正常输出版本
- [ ] 无 URL 时返回错误码 1

---

### 4.6 实现 Execute 函数（完整流程）

**文件**: `internal/cli/cli.go`

**任务**:
- 完整实现 Execute 函数的 8 步流程（plan.md 第 6 节）
- 步骤 4: ParseURL
- 步骤 5: NewClient + 检查 Token
- 步骤 6: FetchIssue/PR/Discussion
- 步骤 7: ToMarkdown
- 步骤 8: Write output

**验收**:
- [ ] `go test ./internal/cli -v -run TestExecute` 通过
- [ ] 完整流程可运行
- [ ] 错误码正确返回
- [ ] stdout 输出正确
- [ ] 文件输出正确

---

### 4.7 实现 Run 函数

**文件**: `internal/cli/cli.go`

**任务**:
- 实现 `Run() int` 函数
- 调用 Execute(os.Stdin, os.Stdout, os.Stderr, os.Args[1:])

**验收**:
- [ ] Run 函数正确调用 Execute

---

### 4.8 创建 cmd/issue2md 入口

**文件**: `cmd/issue2md/main.go`

**任务**:
- 创建 `main()` 函数
- 调用 `cli.Run()`
- 将返回值作为 `os.Exit()` 参数

**验收**:
- [ ] `go build ./cmd/issue2md` 成功
- [ ] 可执行文件可运行
- [ ] 所有功能正常工作

---

## Phase 5: 验收与完善

### 5.1 全量测试

**任务**:
- 运行所有测试
- 生成覆盖率报告

**验收**:
- [ ] `go test ./... -v` 全部通过
- [ ] 覆盖率 `internal/github` ≥ 80%
- [ ] 覆盖率 `internal/converter` ≥ 80%
- [ ] 覆盖率 `internal/config` ≥ 80%
- [ ] 覆盖率 `internal/parser` ≥ 80%

---

### 5.2 端到端测试

**任务**:
- 测试真实 Issue 转换
- 测试真实 PR 转换
- 测试 stdout 输出
- 测试文件输出

**验收**:
- [ ] 真实 Issue 转换正确
- [ ] 真实 PR 转换正确
- [ ] 输出符合 spec.md 示例

---

### 5.3 代码质量检查

**任务**:
- `go fmt` 检查
- `go vet` 检查
- 检查所有错误是否使用 `%w` 包装

**验收**:
- [ ] `go fmt ./...` 无需修改
- [ ] `go vet ./...` 无警告
- [ ] 所有错误正确包装

---

## 任务依赖关系图

```
Phase 1: Foundation
├── 1.1 [P] parser 包基础结构
├── 1.2 [P] config 包基础结构
├── 1.3 [P] github 包数据结构
└── 1.4    退出码常量

Phase 2: GitHub Fetcher
├── 2.1 [P] parser 测试 ─────────────┐
├── 2.2 [P] ParseURL 实现 <───────────┤
├── 2.3     github 客户端基础 <───────┤── 1.3
├── 2.4     集成测试框架 <────────────┤
├── 2.5     FetchIssue 测试 <─────────┤
├── 2.6     FetchIssue 实现 <─────────┤
├── 2.7     FetchPullRequest 测试 ────┤
└── 2.8     FetchPullRequest 实现 ────┘

Phase 3: Markdown Converter
├── 3.1 [P] Issue 测试 <─────────────────────┐
├── 3.2 [P] 辅助函数 <────────────────────────┤── 2.6
├── 3.3     IssueToMarkdown <─────────────────┤
├── 3.4 [P] PR 测试 <────────────────────────┐── 2.8
├── 3.5     PullRequestToMarkdown <──────────┤
├── 3.6 [P] Discussion 测试 <────────────────┤
└── 3.7     DiscussionToMarkdown <───────────┘

Phase 4: CLI Assembly
├── 4.1 [P] config 测试 <───────────────────┐── 1.2
├── 4.2 [P] ParseFlags <────────────────────┤
├── 4.3 [P] cli 测试框架 <──────────────────┤
├── 4.4 [P] cli 基础结构 <──────────────────┤── 1.4
├── 4.5     Execute 基础 <──────────────────┤── 4.2
├── 4.6     Execute 完整 <──────────────────┼── 2.2, 2.6, 2.8, 3.3, 3.5, 3.7
├── 4.7     Run <───────────────────────────┤
└── 4.8     cmd/issue2md/main <─────────────┘── 4.7

Phase 5: 验收
├── 5.1 全量测试 <───────────────────────────┤── 所有 Phase 1-4
├── 5.2 端到端测试 <─────────────────────────┤
└── 5.3 代码质量检查 <───────────────────────┘
```

---

## 任务统计

| Phase | 任务数 | 测试任务 | 实现任务 |
|-------|--------|----------|----------|
| Phase 1 | 4 | 0 | 4 |
| Phase 2 | 8 | 4 | 4 |
| Phase 3 | 7 | 3 | 4 |
| Phase 4 | 8 | 3 | 5 |
| Phase 5 | 3 | 0 | 3 |
| **总计** | **30** | **10** | **20** |

---

## 执行建议

1. **严格按照顺序执行**：任务按依赖关系排列，不要跳过
2. **TDD 循环**：每个测试任务后紧跟实现任务，确保 Red-Green-Refactor
3. **小步提交**：每个任务完成后提交代码
4. **并行机会**：标记 `[P]` 的任务可以与其他 `[P]` 任务同时进行
5. **集成测试注意**：Phase 2 的集成测试需要真实 GitHub Token，提前准备好

---

## 附录：快速检查清单

### 开始每个任务前

- [ ] 理解任务目标和验收标准
- [ ] 阅读相关 spec.md 章节
- [ ] 查看依赖任务的代码

### 提交代码前

- [ ] `go test` 通过
- [ ] `go fmt` 运行
- [ ] `go vet` 运行
- [ ] 错误使用 `%w` 包装
- [ ] 添加必要的注释

### 完成每个 Phase 后

- [ ] 运行 Phase 内所有测试
- [ ] 手动验证功能
- [ ] 检查代码一致性
