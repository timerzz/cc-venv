# ccv 系统总设计

## 1. 项目定位

`ccv` 是一个用于管理命名 Claude Code 虚拟环境的工具。

它的核心能力是：

- 创建环境
- 进入环境终端
- 在环境中运行 Claude Code
- 导出 / 导入环境快照
- 通过本地 Web 界面管理环境

一句话定义：

**ccv 是一个面向 Claude Code 的命名环境管理器与环境镜像工具。**

---

## 2. 设计边界

第一版设计边界如下：

- 只管理命名环境
- 不做项目目录自动发现
- 不做声明式依赖求解
- 以环境目录本身作为事实来源
- 支持导出 / 导入完整环境快照
- 提供 CLI 和本地 Web 管理入口

`ccv` 只管理命名环境目录中的资源，不管理项目级 Claude 资源。

明确不管理：

- 项目目录里的 `.claude/`
- 项目级 `.mcp.json`
- 项目级 commands / rules / skills / agents / hooks

---

## 3. 环境模型

第一版只保留一种环境：

- **全局命名环境**

统一路径：

```text
~/.ccv/envs/<name>/
```

每个环境彼此独立，按名称管理。

环境中的实际目录内容是事实来源，`ccv.json` 只是摘要与导出说明。

---

## 4. 目录模型

环境目录采用“原生兼容优先”的 Claude Code 目录布局：

```text
~/.ccv/envs/<name>/
├── ccv.json
├── .claude/
│   ├── CLAUDE.md
│   ├── settings.json
│   ├── settings.local.json
│   ├── commands/
│   ├── agents/
│   ├── skills/
│   ├── rules/
│   └── plugins/
│       ├── cache/
│       └── data/
└── .claude.json
```

关键含义：

- `.claude/CLAUDE.md`：环境级 Claude 记忆 / 指令
- `.claude/settings.json`：设置与插件启用配置
- `.claude/commands/`：环境级 commands
- `.claude/agents/`：环境级 agents
- `.claude/skills/`：环境级 skills
- `.claude/rules/`：环境级 rules
- `.claude/plugins/cache/`：插件本体缓存
- `.claude/plugins/data/`：插件持久数据
- `.claude.json`：环境级 MCP / 用户级配置镜像
- `ccv.json`：环境摘要与导出信息

第一版明确不做：

- `vendor`
- `state.json`
- lock 文件
- 外部 executable 安装器

---

## 5. 命令模型

第一版命令集合：

```bash
ccv create <name>
ccv list
ccv active <name>
ccv remove <name>
ccv run <name>
ccv web
ccv export <name>
ccv import <archive>
```

命令职责：

- `create`：创建一个新的命名环境
- `list`：列出所有环境及摘要
- `active`：进入环境终端，不自动执行 Claude Code
- `remove`：删除环境
- `run`：在环境中执行 Claude Code
- `web`：打开本地 Web 管理界面
- `export`：导出环境快照
- `import`：导入环境快照

---

## 6. 运行时原则

进入环境终端或执行 Claude Code 时，`ccv` 需要为进程注入该环境对应的变量。

这些变量包括：

- Claude 配置目录相关变量
- 环境标识变量
- 环境自定义变量

具体注入方式、shell 行为、工作目录策略：

- 在 `docs/02-ccv-active-design.md` 中单独定义

---

## 7. 导出 / 导入原则

导出 / 导入只覆盖命名环境目录中的资源。

不包含：

- 任何项目级 Claude 资源
- `.claude/plugins/data/`
- `.claude/plugins/cache/` 下以 `temp` 开头的目录

导出前必须基于当前环境目录刷新摘要。

具体扫描规则和导出排除规则：

- 在 `docs/03-ccv-scan-design.md` 中单独定义

---

## 8. Web 管理原则

`ccv web` 第一版建议收敛为：

- 启动本地 HTTP 服务
- 自动打开浏览器
- 只管理本机 `~/.ccv/envs/` 下的环境

第一版 Web UI 重点能力：

- 浏览环境列表
- 查看环境详情
- 删除环境
- 触发导入 / 导出

---

## 9. Go 模块分层

建议项目结构：

```text
cc-venv/
├── cmd/ccv/
├── internal/cli/
├── internal/env/
├── internal/config/
├── internal/exporter/
├── internal/importer/
├── internal/archive/
├── internal/webui/
└── internal/platform/
```

职责划分：

- `internal/cli`：命令树、参数解析、输出
- `internal/env`：环境路径、创建、删除、激活、运行
- `internal/config`：`ccv.json` 与扫描相关数据结构
- `internal/exporter`：导出编排
- `internal/importer`：导入编排
- `internal/archive`：tar.gz 和 checksums
- `internal/webui`：本地 Web 服务
- `internal/platform`：home 目录和平台差异

---

## 10. 当前实现顺序建议

建议按下面顺序推进：

1. `create`
2. `list`
3. `active`
4. `run`
5. `remove`
6. `scan`
7. `export`
8. `import`
9. `web`
