# Web 前端待办

## Phase 1: 工程基础

- [x] 建立 `web/` 前端工程骨架
- [x] 接入 `Vite + React + TypeScript`
- [x] 接入 `Tailwind CSS`
- [x] 接入 `i18n` 基础设施
- [x] 约定目录结构：`layouts`、`components`、`pages`、`lib`、`i18n`
- [ ] 引入图标方案
- [ ] 统一设计 token：颜色、圆角、阴影、间距、状态色
- [ ] 建立前端环境变量约定

## Phase 2: 应用壳层

- [x] 实现 Dashboard 初版页面壳
- [x] 拆出环境侧栏组件
- [x] 拆出状态卡组件
- [x] 接入语言切换入口
- [x] 抽出通用 `AppShell`
- [x] 引入路由层并预留页面级结构
- [ ] 增加全局空态、错误态、加载态组件
- [ ] 增加通知反馈机制

## Phase 3: 数据层

- [x] 建立最小 `api.ts`
- [x] 按领域拆分 API：`env.ts`、`llm.ts`、`mcp.ts`、`skills.ts`、`import-export.ts`
- [x] 建立统一请求错误处理
- [ ] 考虑接入 `TanStack Query`
- [ ] 统一 `Response<T>` 的前端类型和解包逻辑

## Phase 4: 页面实现

- [x] 环境总览页
- [x] LLM 配置页
- [x] MCP 管理页
- [x] Skills 管理页
- [x] 环境变量页
- [x] `CLAUDE.md` / Notes 编辑页
- [x] Agents 管理页
- [x] Commands 管理页
- [x] Rules 管理页
- [x] 导入导出页

## Phase 5: 表单与交互

- [x] 创建环境弹窗
- [x] 删除环境确认
- [x] 导入上传交互
- [x] 导出状态反馈
- [x] 导入成功后刷新环境列表并自动选中
- [ ] 表单校验和脏状态提示
- [ ] 敏感值默认隐藏 / 显示切换

## Phase 6: i18n 完善

- [x] 支持 `zh-CN`
- [x] 支持 `en`
- [x] 浏览器语言检测
- [x] 手动语言切换
- [x] 本地持久化语言偏好
- [ ] 统一技术词汇翻译策略
- [ ] 覆盖动态数量、错误消息、时间格式
- [ ] 梳理不翻译的专业名词白名单

## Phase 7: 质量收口

- [ ] 窄屏与移动端适配
- [ ] 键盘可访问性检查
- [ ] 颜色对比度检查
- [x] 前端构建产物接入 Go embed 流程
- [ ] 增加基本前端测试
- [ ] 增加开发文档和启动说明

## 当前实现顺序

1. `Tailwind CSS` 和 `i18n`
2. `Dashboard` 壳层重构
3. 数据层拆分
4. 环境总览页和 LLM/MCP/Skills 三个核心页面
5. 导入导出和危险操作交互
6. Agents / Commands / Rules 文件资源管理
