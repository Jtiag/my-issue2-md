# Code Review Report — `internal/`

## 静态检查结果

`go vet ./internal/...` — **无问题报告**

---

## 总体评价

代码整体结构清晰，包职责划分合理，错误处理基本遵循了宪法的 `%w` 包装要求，但在代码复用、硬编码占位符和死代码方面存在明显改进空间。

---

## 优点

1. **表格驱动测试做得好** — `converter/issue_test.go`、`converter/pullrequest_test.go`、`converter/discussion_test.go`、`parser/parser_test.go` 等测试文件均严格遵循表格驱动风格，覆盖了多种边界场景，符合宪法第二条 2.2。
2. **错误包装一致性高** — `client.go`、`parser/parser.go` 中的错误返回均使用了 `fmt.Errorf("...: %w", err)` 格式，符合宪法第三条 3.1。

---

## 待改进项

### [高优先级]

#### 1. 硬编码占位符字符串（影响正确性）

- **文件:** `converter/issue.go:39-40`
- **问题:** Issue URL 链接使用了硬编码的 `"placeholder"` 和 `"owner"/"repo"` 字符串，导致生成的 Markdown 中链接始终指向错误的地址。
  ```go
  b.WriteString(fmt.Sprintf("**Issue:** [%s/#%d](https://github.com/%s/%s/issues/%d)\n\n",
      "placeholder", issue.Number, "owner", "repo", issue.Number))
  ```
- **文件:** `converter/pullrequest.go:34-35` — 同样的问题：
  ```go
  b.WriteString(fmt.Sprintf("**Pull Request:** [#%d](https://github.com/owner/repo/pull/%d)\n\n",
      pr.Number, pr.Number))
  ```
- **建议:** 在 `Issue`/`PullRequest` 结构体中增加 `Owner` 和 `Repo` 字段，由上游在 fetch 时填充，converter 使用这些字段生成正确的链接。或者将 owner/repo 作为参数传入转换函数。

#### 2. Reactions 转换逻辑大量重复（违反宪法第一条 1.3 反过度工程）

- **文件:** `github/client.go:68-79`、`github/client.go:101-112`、`github/client.go:223-234`
- **问题:** 将 `ghIssue.Reactions`/`ghComment.Reactions` 转换为内部 `Reactions` 结构体的逻辑被复制粘贴了 **3 次**，每次都是完全相同的 8 个字段赋值。
- **建议:** 提取一个私有函数，例如：
  ```go
  func convertReactions(r *github.Reactions) Reactions {
      return Reactions{
          ThumbsUp:   r.GetPlusOne(),
          ThumbsDown: r.GetMinusOne(),
          // ... 其余字段
      }
  }
  ```

#### 3. 评论获取逻辑重复（违反宪法第一条 1.3）

- **文件:** `github/client.go:86-121`（FetchIssue 中的评论获取）、`github/client.go:205-243`（FetchPullRequest 中的评论获取）
- **问题:** 两个方法中的评论分页获取 + Reactions 转换逻辑几乎完全一致。
- **建议:** 提取一个私有方法 `fetchComments(ctx, owner, repo, number) ([]Comment, error)` 到 `Client` 上。

#### 4. 死代码未清理

- **文件:** `converter/formatter.go:178-185`
- **问题:** `toConverterOptions` 函数接受 `interface{}` 参数但从未被调用，且内部硬编码返回 false。这是典型的死代码。
- **建议:** 直接删除此函数。

#### 5. `strings.Contains` 被手写替代（违反宪法第一条 1.2 标准库优先）

- **文件:** `converter/issue_test.go:309-321`
- **问题:** 手动实现了 `contains()` 和 `findSubstring()` 函数来替代标准库 `strings.Contains`。
  ```go
  func contains(s, substr string) bool {
      return len(s) >= len(substr) && findSubstring(s, substr)
  }
  ```
- **建议:** 替换为 `strings.Contains(s, substr)`。当前函数名也与标准库冲突，容易造成混淆。

---

### [中优先级]

#### 6. YAML Front Matter 手工拼接存在注入风险

- **文件:** `converter/formatter.go:13-43`、`converter/pullrequest.go:123-161`、`converter/discussion.go:28-37`
- **问题:** 使用 `fmt.Sprintf("title: %q\n", issue.Title)` 虽然用了 `%q` 加引号，但如果 title 本身包含 `"` 或换行符，生成的 YAML 可能格式不正确。更健壮的方式是使用专门的 YAML 库或对特殊字符做转义处理。
- **建议:** 至少对 title/body 中可能包含的特殊字符做检查或转义。考虑到宪法第一条"标准库优先"，可以写一个简单的 YAML 字符串转义函数。

#### 7. `Response` 返回值被 `_` 丢弃

- **文件:** `github/client.go:33`、`github/client.go:130`、`github/client.go:183`
- **问题:** `c.githubClient.Issues.Get()` 等调用返回 `(result, *Response, error)`，其中 Response 包含 rate limit 信息，一律用 `_` 丢弃。虽然不违反宪法第三条（error 已处理），但丢失了 rate limit 信息可能导致生产环境中遇到限流时无任何提示。
- **建议:** 至少记录 response 的 `Rate` 信息到日志中，或在接近限流时返回警告。

#### 8. `Discussion` 类型缺少 `UpdatedAt` 和 `Milestone` 字段

- **文件:** `github/types.go:60-77`
- **问题:** `Discussion` 结构体缺少 `UpdatedAt` 字段，也缺少 `Labels`、`Milestone`、`Reactions` 等字段，而 `Issue` 和 `PullRequest` 都有。这可能是功能未完成的表现。
- **建议:** 评估是否需要为 Discussion 添加这些字段。如果 GraphQL API 支持这些数据，应补齐。

#### 9. 退出码文档与实际定义不一致

- **文件:** `internal/config/config.go:142-148`（HelpText 中的 Exit Codes 描述）vs `internal/cli/exit_codes.go:8-22`
- **问题:** `HelpText()` 中写的退出码含义与 `exit_codes.go` 中定义的常量不匹配：
  - HelpText 说 2 = "GitHub API error"，但 `exit_codes.go` 定义 `ExitNotFound = 2`
  - HelpText 说 3 = "Resource not found"，但 `exit_codes.go` 定义 `ExitAPIFailed = 3`
- **建议:** 统一 `HelpText()` 中的描述与 `exit_codes.go` 的定义。

---

### [低优先级]

#### 10. `convertReactions` 中的 `converterOptions` 可简化

- **文件:** `converter/formatter.go:171-175`
- **问题:** `converterOptions` 结构体与 `config.Options` 完全一致（都有 `enableReactions` 和 `enableUserLinks`），增加了不必要的中间层。
- **建议:** 考虑直接使用 `config.Options`，消除 `converterOptions` 和 `toConverterOptions`。

#### 11. `cli_test.go` 中基于测试名的条件判断

- **文件:** `internal/cli/cli_test.go:96`
- **问题:** 使用 `strings.Contains(tt.name, "GITHUB_TOKEN")` 来决定是否跳过测试。这种基于名字的条件判断很脆弱。
  ```go
  if strings.Contains(tt.name, "GITHUB_TOKEN") && os.Getenv("GITHUB_TOKEN") == "" {
      t.Skip("GITHUB_TOKEN not set, skipping token-dependent test")
  }
  ```
- **建议:** 在测试结构体中增加一个 `skipWithoutToken bool` 字段来显式控制。

#### 12. `TestRun` 测试实际上是空操作

- **文件:** `internal/cli/cli_test.go:133-144`
- **问题:** `TestRun` 函数注释说"我们无法真正测试 Run() 因为它调用 os.Exit()"，但函数体几乎为空，只保存和恢复 `os.Args`，没有实际断言。
- **建议:** 要么删除此测试，要么使用子进程测试模式来真正验证 `Run()` 的行为。

#### 13. `generateTestPath` 使用 rune 转换不够直观

- **文件:** `converter/pullrequest_test.go:241`
- **问题:** `string(rune('a'+i%26))` 通过 Unicode 码点偏移生成字符，可读性不佳。
  ```go
  return "path/to/file" + string(rune('a'+i%26)) + ".go"
  ```
- **建议:** 使用 `fmt.Sprintf("path/to/file%d.go", i)` 更直观。

---

## 审查总结

| 优先级 | 数量 | 关键问题 |
|--------|------|----------|
| 高     | 5    | 硬编码占位符、代码重复（2处）、死代码、手写标准库函数 |
| 中     | 4    | YAML注入风险、Response丢弃、字段缺失、文档不一致 |
| 低     | 4    | 结构体冗余、脆弱测试、空测试、可读性 |
