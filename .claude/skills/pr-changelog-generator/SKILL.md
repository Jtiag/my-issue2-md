---
name: pr-changelog-generator
description: 自动分析 Git 提交历史，生成规范的 PR Changelog。当用户准备提 PR、请求总结代码变更、或是要求“写一下更新日志”时，请务必优先触发本技能。
allowed-tools: Bash Read
metadata:
  author: tonybai
---

# PR Changelog Generator 技能指南

这个技能旨在帮助我们的开发者减轻提 PR 时的负担。我们团队非常看重 Review 的效率，一份结构清晰、重点突出的 Changelog 能极大节约 Reviewer 的时间，减少沟通摩擦。

## 执行步骤

当决定使用本技能时，请遵循以下流程：

### 1. 收集数据 (确定性执行)

请不要自己去猜测 Git 命令。请直接调用我们准备好的安全脚本来获取原始提交数据：

执行 `Bash(sh scripts/extract_commits.sh)`

### 2. 意图理解与提炼 (你的强项)

仔细阅读脚本返回的 commits。你需要发挥你的理解能力：

- 剔除那些无意义的提交（如 "fix typo", "update"）。
- 将连续相关的提交合并为一个功能点。
- **重点：** 请站在 Reviewer 的角度，思考他们最关心哪些变更？把高风险、大范围的改动排在最前面。

### 3. 结构化输出

请使用以下模板输出。**为什么要中英双语？因为我们是一个跨国协作团队，这能确保不同母语的同事都能快速理解。**

    ## 🚀 Changes (变更内容)
    * **[Feature/Fix]**: (English description) / (中文描述)
    * ...
