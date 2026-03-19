# ccv 资源扫描设计

## 1. 目标

资源扫描用于识别命名环境中的实际资源状态，为以下能力提供数据：

- `ccv export <name>`
- `ccv list`
- `ccv web`
- `ccv.json` 摘要刷新

扫描只负责读取当前状态，不负责安装或同步资源。

---

## 2. 扫描边界

扫描只面向命名环境目录：

```text
~/.ccv/envs/<name>/
```

只扫描该目录中的受管资源，不扫描：

- 项目目录 `.claude/`
- 项目级 `.mcp.json`
- 项目级 commands / rules / skills / agents / hooks
- 环境目录外的任意文件

---

## 3. 扫描模型

扫描器基于原生兼容目录工作：

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

第一版扫描输出分两类：

- 轻量摘要：给 `list` 和 Web 列表页使用
- 完整摘要：给 `export`、Web 详情页和 `ccv.json` 刷新使用

---

## 4. 识别规则

第一版只保留稳定、容易解释的识别规则。

- `CLAUDE.md`
  - 位置：`.claude/CLAUDE.md`
  - 规则：存在即记为有
- `skills`
  - 位置：`.claude/skills/`
  - 规则：每个一级子目录视为一个 skill
- `agents`
  - 位置：`.claude/agents/`
  - 规则：每个一级 `*.md` 文件视为一个 agent
- `commands`
  - 位置：`.claude/commands/`
  - 规则：每个一级 `*.md` 文件视为一个 command
- `rules`
  - 位置：`.claude/rules/`
  - 规则：递归扫描 `*.md`，名称取相对路径去掉扩展名
- `plugins`
  - 位置：`.claude/plugins/cache/`
  - 规则：每个一级目录视为一个插件缓存项，但忽略 `temp*`
- `MCP`
  - 位置：`.claude.json`
  - 规则：从该文件读取 MCP server 定义

---

## 5. 扫描时机

- `create`
  - 不扫描
- `active`
  - 不扫描
- `run`
  - 不扫描
- `export`
  - 必须扫描
- `import`
  - 不扫描目标环境现状，但要校验导入包
- `list`
  - 第一版默认读 `ccv.json`
- `web`
  - 列表页读 `ccv.json`
  - 详情页可实时扫描

---

## 6. 与 `ccv.json` 的关系

- 环境目录是事实来源
- `ccv.json` 是摘要和导出元数据
- `export` 前必须基于扫描结果刷新 `ccv.json`
- `list` 默认读取已有 `ccv.json`

---

## 7. 导出排除规则

导出时默认不包含：

- `.claude/plugins/data/`
- `.claude/plugins/cache/` 下所有 `temp*` 目录

这些规则同时影响扫描展示和打包边界：

- `data/` 不作为可导出资源
- `temp*` 不作为稳定插件资源
