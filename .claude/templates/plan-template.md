# [项目名称] 技术实现方案

**文档版本**: {{version}}
**创建日期**: {{date}}
**适用规格**: {{spec_path}}
**遵循宪法**: `constitution.md`

---

## 1. 技术上下文总结

### 1.1 技术选型

| 技术领域 | 选型 | 理由 |
|---------|------|------|
| **编程语言** | {{language}} | {{reason}} |
| **HTTP 客户端** | {{http_client}} | {{reason}} |
| **数据存储** | {{storage}} | {{reason}} |
| **配置管理** | {{config}} | {{reason}} |
| **日志记录** | {{logging}} | {{reason}} |

### 1.2 依赖管理

**最小外部依赖原则**：
```
require (
    {{dependency_list}}
)
```

### 1.3 性能目标

| 指标 | 目标值 | 实现策略 |
|------|--------|---------|
| {{metric_1}} | {{target_1}} | {{strategy_1}} |
| {{metric_2}} | {{target_2}} | {{strategy_2}} |

---

## 2. "合宪性"审查

### 2.1 第一条：简单性原则 (Simplicity First)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **1.1 YAGNI** | {{yagni_measure}} | {{status}} |
| **1.2 标准库优先** | {{stdlib_measure}} | {{status}} |
| **1.3 反过度工程** | {{anti_overengineering_measure}} | {{status}} |

**反模式对照**：
```go
// ❌ 过度工程：不必要的抽象
{{bad_example}}

// ✅ 简单实现：直接函数
{{good_example}}
```

### 2.2 第二条：测试先行铁律 (Test-First Imperative)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **2.1 TDD 循环** | {{tdd_measure}} | {{status}} |
| **2.2 表格驱动** | {{table_driven_measure}} | {{status}} |
| **2.3 拒绝 Mocks** | {{no_mock_measure}} | {{status}} |

**测试架构**：
```
{{test_architecture_diagram}}
```

### 2.3 第三条：明确性原则 (Clarity and Explicitness)

| 宪法条款 | 本方案措施 | 状态 |
|---------|-----------|------|
| **3.1 错误处理** | {{error_handling_measure}} | {{status}} |
| **3.2 无全局变量** | {{no_global_measure}} | {{status}} |

**错误处理示例**：
```go
{{error_handling_example}}
```

**依赖注入示例**：
```go
{{dependency_injection_example}}
```

---

## 3. 项目结构细化

### 3.1 目录树

```
{{directory_tree}}
```

### 3.2 包职责矩阵

| 包 | 职责 | 依赖 | 导出 |
|---|------|------|------|
| {{package_1}} | {{responsibility_1}} | {{deps_1}} | {{exports_1}} |
| {{package_2}} | {{responsibility_2}} | {{deps_2}} | {{exports_2}} |

### 3.3 依赖关系图

```
{{dependency_diagram}}
```

**依赖原则**：
- {{principle_1}}
- {{principle_2}}

---

## 4. 核心数据结构

### 4.1 主要类型定义

```go
// {{struct_name}} {{description}}
type {{struct_name}} struct {
    {{field_1}}    {{type_1}}    // {{comment_1}}
    {{field_2}}    {{type_2}}    // {{comment_2}}
}
```

### 4.2 类型关系图

```
{{type_relationship_diagram}}
```

---

## 5. 接口设计

### 5.1 包级接口

```go
// {{package}}/{{file}}.go

// {{function_name}} {{description}}
func {{function_name}}({{params}}) ({{returns}}, error) {
    // 实现细节
}
```

### 5.2 接口清单

| 接口 | 包 | 签名 | 说明 |
|------|---|------|------|
| {{interface_1}} | {{package_1}} | {{signature_1}} | {{desc_1}} |
| {{interface_2}} | {{package_2}} | {{signature_2}} | {{desc_2}} |

---

## 6. 实现路线图

### Phase {{phase_number}}: {{phase_name}}

| 任务 | 优先级 | 验收标准 |
|------|--------|---------|
| {{task_1}} | {{priority_1}} | {{acceptance_1}} |
| {{task_2}} | {{priority_2}} | {{acceptance_2}} |

---

## 7. 风险与缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| {{risk_1}} | {{impact_1}} | {{mitigation_1}} |
| {{risk_2}} | {{impact_2}} | {{mitigation_2}} |

---

## 8. 附录

### 8.1 关键配置示例

```yaml
{{config_example}}
```

### 8.2 API 端点映射

| 端点 | 方法 | 说明 |
|------|------|------|
| {{endpoint_1}} | {{method_1}} | {{desc_1}} |

### 8.3 错误码映射

| 错误码 | 含义 | 触发条件 |
|-------|------|---------|
| {{code_1}} | {{meaning_1}} | {{condition_1}} |
