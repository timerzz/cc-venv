# ccv 文档索引

当前 `docs/` 保留核心文档与专项设计文档，避免历史草稿和现行规格混在一起。

## 文档分工

- `00-ccv-system-design.md`
  - 系统总设计
  - 包括边界、环境模型、目录模型、命令模型、导出导入原则、模块分层
- `01-ccv-go-implementation-design.md`
  - Go 项目结构、包边界、命令调用关系、MVP 实现顺序
- `02-ccv-active-design.md`
  - `ccv active` 的专项设计
  - 包括环境变量注入、shell 启动、cwd 策略、与 `run` 的关系
- `03-ccv-scan-design.md`
  - 资源扫描专项设计
  - 包括原生兼容目录、资源识别规则、导出排除规则、扫描时机
- `04-ccv-list-design.md`
  - `ccv list` 的专项设计
  - 包括数据来源、摘要展示、排序、降级策略、与扫描和 Web 的关系
- `05-ccv-web-design.md`
  - `ccv web` 的专项设计
  - 包括 Gin 后端、React 前端、embed 分发、API 边界、页面模型
- `06-web-frontend-todo.md`
  - Web 前端实现待办
  - 包括页面完成度、交互收口和剩余质量项

## 阅读顺序

建议按下面顺序看：

1. `00-ccv-system-design.md`
2. `01-ccv-go-implementation-design.md`
3. `02-ccv-active-design.md`
4. `03-ccv-scan-design.md`
5. `04-ccv-list-design.md`
6. `05-ccv-web-design.md`
7. `06-web-frontend-todo.md`
