**Prompt 2: 生成任务列表**

方案非常完美。

现在，请扮演技术组长。请仔细阅读 `@./specs/001-core-functionality/spec.md` 和 `@./specs/001-core-functionality/plan.md`。

你的目标是将 `plan.md`中描述的技术方案，分解成一个**详尽的、原子化的、有依赖关系的、可被AI直接执行的任务列表**。

务必每个任务打上唯一的序号
eg:
1.2.1 某xx功能
1.创建 xx1
2.创建 xx2
3.创建 xx3

**关键要求：**
1.  **任务粒度：** 每个任务应该只涉及一个主要文件的修改或创建一个新文件。不要出现“实现所有功能”这种大任务。
2.  **TDD强制：** 根据`constitution.md`的“测试先行铁律”，**必须**先生成测试任务，后生成实现任务。
3.  **并行标记：** 对于没有依赖关系的任务，请标记 `[P]`。
4.  **阶段划分：** 即便`plan.md`中包含了粗略的阶段划分，也要以下面的为准。
    *   **Phase 1: Foundation** (数据结构定义)
    *   **Phase 2: GitHub Fetcher** (API交互逻辑，TDD)
    *   **Phase 3: Markdown Converter** (转换逻辑，TDD)
    *   **Phase 4: CLI Assembly** (命令行入口集成)

完成后，将生成的任务列表写入到`./specs/001-core-functionality/tasks.md`文件中。
