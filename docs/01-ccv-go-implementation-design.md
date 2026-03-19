# ccv Go 实现结构设计

## 1. 目标

本文件把当前系统设计收敛成一版 Go 实现分层，重点说明：

- 项目目录结构
- 包职责边界
- 命令调用关系
- MVP 实现顺序

本文件不展开命令交互细节和扫描细则：

- `active` 细节见 `docs/02-ccv-active-design.md`
- 扫描与导出排除规则见 `docs/03-ccv-scan-design.md`

---

## 2. 实现原则

- 只管理命名环境：`~/.ccv/envs/<name>/`
- 不做项目目录自动发现
- 环境目录本体是事实来源
- `ccv.json` 是摘要与导出元数据，不是唯一事实来源
- 环境目录采用 Claude Code 原生兼容优先的布局
- CLI、Web 复用同一套环境读写逻辑
- 导出 / 导入只覆盖命名环境目录，不覆盖项目级 Claude 资源

---

## 3. 项目结构

建议并对齐当前代码的项目结构如下：

```text
cc-venv/
├── cmd/
│   └── ccv/
│       └── main.go
├── internal/
│   ├── archive/
│   │   ├── checksum.go
│   │   └── tarball.go
│   ├── cli/
│   │   ├── app.go
│   │   ├── active_cmd.go
│   │   ├── create_cmd.go
│   │   ├── export_cmd.go
│   │   ├── import_cmd.go
│   │   ├── list_cmd.go
│   │   ├── remove_cmd.go
│   │   ├── run_cmd.go
│   │   └── web_cmd.go
│   ├── config/
│   │   ├── ccvjson.go
│   │   ├── envvars.go
│   │   ├── manifest.go
│   │   └── scan.go
│   ├── env/
│   │   ├── create.go
│   │   ├── layout.go
│   │   ├── run.go
│   │   └── types.go
│   ├── exporter/
│   │   └── export.go
│   ├── importer/
│   │   └── import.go
│   ├── platform/
│   │   └── home.go
│   └── webui/
│       └── server.go
├── go.mod
└── go.sum
```

---

## 4. 包职责

### `cmd/ccv`

职责：

- CLI 入口
- 调用 `internal/cli`

### `internal/cli`

职责：

- 定义 Cobra 命令树
- 解析参数与 flag
- 调用内部服务
- 输出用户可读结果

约束：

- 不直接读写环境目录
- 不直接处理 tar.gz 和扫描细节

### `internal/env`

职责：

- 环境路径计算
- 环境创建、加载、列出、删除
- 为 `active` / `run` 准备进程执行上下文

建议核心类型保持简单：

```go
type Environment struct {
    Name            string
    RootPath        string
    ManifestPath    string
    ClaudeConfigDir string
}
```

这里的 `ClaudeConfigDir` 表示运行时传给 Claude 的配置根目录。其具体取值和注入策略，按 `docs/02-ccv-active-design.md` 收敛。

### `internal/config`

职责：

- 读写 `ccv.json`
- 读写导出包 `manifest.json`
- 读写环境变量配置
- 扫描环境目录中的受管资源

扫描范围应与系统设计一致，围绕原生兼容目录工作，例如：

- `.claude/CLAUDE.md`
- `.claude/settings.json`
- `.claude/commands/`
- `.claude/agents/`
- `.claude/skills/`
- `.claude/rules/`
- `.claude/plugins/cache/`
- `.claude.json`

### `internal/exporter`

职责：

- 编排 `ccv export <name>`
- 调用环境加载、扫描、manifest 生成、打包与校验逻辑

### `internal/importer`

职责：

- 编排 `ccv import <archive>`
- 调用解包、校验、目标路径恢复逻辑

### `internal/archive`

职责：

- tar.gz 打包 / 解包
- checksums 生成 / 校验

约束：

- 不依赖 `env`
- 不感知命令语义

### `internal/webui`

职责：

- 启动本地 HTTP 服务
- 提供环境列表、详情、删除、导入导出入口

### `internal/platform`

职责：

- home 目录、平台差异等基础能力

---

## 5. 目录模型对实现的要求

环境目录按原生兼容优先组织：

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

实现层需要满足：

- `create` 创建上述基础目录和最小初始化文件
- `active` / `run` 以该环境为运行上下文
- `scan` 基于上述目录生成摘要
- `export` 排除 `.claude/plugins/data/`
- `export` 排除 `.claude/plugins/cache/` 下 `temp*` 目录

---

## 6. 命令调用关系

```text
cli -> env
cli -> exporter -> env + config + archive
cli -> importer -> env + config + archive
cli -> webui -> env + config
```

补充约束：

- `cli` 不直接操作环境目录
- `config` 不负责启动进程
- `env` 不负责 tar.gz 归档
- `webui` 不重复实现环境读写逻辑

---

## 7. 关键流程

### `ccv create <name>`

流程：

1. 校验名称
2. 计算环境路径
3. 创建原生兼容目录骨架
4. 生成初始 `ccv.json`

### `ccv list`

流程：

1. 扫描 `~/.ccv/envs/`
2. 加载每个环境的 `ccv.json`
3. 输出名称、路径和摘要

第一版默认读已有摘要，不要求实时完整扫描。

### `ccv active <name>`

流程：

1. 加载环境
2. 准备环境变量与工作目录
3. 启动交互式 shell

具体注入和 shell 策略见 `docs/02-ccv-active-design.md`。

### `ccv run <name>`

流程：

1. 加载环境
2. 复用与 `active` 相同的执行上下文准备逻辑
3. 启动 `claude`

### `ccv remove <name>`

流程：

1. 加载环境
2. 做删除确认
3. 删除环境目录

### `ccv export <name>`

流程：

1. 加载环境
2. 扫描环境资源
3. 刷新 `ccv.json`
4. 生成导出 `manifest.json`
5. 应用导出排除规则
6. 打包归档

### `ccv import <archive>`

流程：

1. 解包到临时目录
2. 校验 `manifest.json`
3. 读取目标环境名
4. 恢复到 `~/.ccv/envs/<name>/`

### `ccv web`

流程：

1. 启动本地 HTTP 服务
2. 打开浏览器
3. 提供环境管理页面或接口

---

## 8. MVP 实现顺序

建议按下面顺序推进：

### 第一阶段

- `create`
- `list`
- `active`
- `run`
- `remove`

### 第二阶段

- `scan`
- `export`
- `import`

### 第三阶段

- `web`

---

## 9. 当前文档边界

当前建议把文档分工保持为：

- `00-ccv-system-design.md`
  - 讲系统边界、目录模型、命令模型
- `01-ccv-go-implementation-design.md`
  - 讲 Go 分层、包职责、调用关系、实现顺序
- `02-ccv-active-design.md`
  - 讲 `active` / `run` 的运行时细节
- `03-ccv-scan-design.md`
  - 讲扫描规则和导出排除规则
