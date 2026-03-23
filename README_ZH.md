# cc-venv
![GitHub License](https://img.shields.io/github/license/timerzz/cc-venv)
![Static Badge](https://img.shields.io/badge/github-repo-blue%3Flogo%3Dgithub?color=blue)






[English](./README.md) | 简体中文

`cc-venv` 是一个面向 Claude Code 的命名环境管理器。

它的目标很直接：把不同工作流拆进不同 Claude 环境里，用的时候切换，不动你默认的 Claude Code 配置。

## 为什么需要 cc-venv

很多人用 Claude Code 不只做一件事。

有时候你希望它是“代码工作流”：

- 配好编程相关的 `MCP`
- 安装代码分析或前端设计相关 `Skills`
- 写偏工程的 `Rules`

但有时候你又希望它是“写作工作流”：

- 换一套更偏写作、翻译、总结的提示和规则
- 用不同的 `MCP` 或更少的工具
- 保持和代码环境完全隔离

如果这些内容都堆在默认的 Claude Code 配置里，最后通常会变成：

- `Skills` 混在一起
- `MCP` 配置互相影响
- `Rules` 和 `CLAUDE.md` 互相污染
- 想试新的工作流时，不敢动默认配置

`cc-venv` 的做法是给你一组独立环境：

- 每个环境有自己的 `CLAUDE.md`
- 每个环境有自己的环境变量
- 每个环境有自己的 `MCP`
- 每个环境有自己的 `Skills`
- 每个环境有自己的 `Agents / Commands / Rules`

这样你可以同时维护：

- `coding`
- `writing`
- `research`
- `prod-safe`

切换环境时，不会影响你默认的 Claude Code 配置。


## 快速使用

### 方式一：npm 全局安装

```bash
npm install -g @timerzz/ccv
```

安装完成后可以直接使用：

```bash
ccv list
ccv create coding
ccv web
```

### 方式二：直接下载 Go 二进制

从 GitHub Releases 下载对应平台的二进制，然后加入 `PATH`：

- Linux x64: `ccv-linux-amd64`
- Linux arm64: `ccv-linux-arm64`
- macOS x64: `ccv-darwin-amd64`
- macOS arm64: `ccv-darwin-arm64`
- Windows x64: `ccv-windows-amd64.exe`
- Windows arm64: `ccv-windows-arm64.exe`

常用命令：

```bash
# 创建一个名为coding的虚拟环境
ccv create coding
# 打开web页面对虚拟环境进行管理和配置
ccv web
# 查看当前所有虚拟环境
ccv list
# 进入coding虚拟环境中的claude code
ccv run coding
# 透传额外参数给 claude
ccv run coding --model claude-opus -p "总结一下这个仓库"
# 导出coding虚拟环境的所有配置、包括skills、mcp等
ccv exprot coding
```
<img width="2485" height="1268" alt="web" src="https://github.com/user-attachments/assets/ba6f3240-dac8-414a-8343-ff8c8e5aa488" />



## 和 cc-mirror 的对比

[`cc-mirror`](https://github.com/numman-ali/cc-mirror) 和 `cc-venv` 都能提供 Claude Code 隔离能力，但优化方向不同。

| 能力 | `cc-venv` | `cc-mirror` |
| --- | --- | --- |
| 多个隔离的 Claude Code 环境 | ✅ | ✅ |
| 内嵌 Web UI（`ccv web`） | ✅ | ❌ |
| 完整环境导入 / 导出 | ✅ | ❌ |
| 隔离 `CLAUDE.md`、环境变量、`Skills`、`MCP`、`Agents`、`Commands`、`Rules` | ✅ | ✅ |
| 每个环境独立 Claude Code 二进制 | ❌ | ✅ |
| 面向 provider / prompt-pack / tweak 的变体管理 | ❌ | ✅ |

如果你的目标是：

- 维护几套不同用途的 Claude 工作流
- 快速切换
- 保持默认 Claude 配置干净
- 能把整套环境迁移到另一台机器

那 `cc-venv` 更贴近这个使用方式。

## 核心能力

- 创建命名环境
- 进入环境 shell
- 在环境中运行 Claude Code
- 通过 Web 页面管理环境
- 隔离 `LLM`、环境变量、`MCP`、`Skills`、`Agents`、`Commands`、`Rules`
- 导出整个环境
- 在新机器导入环境并快速恢复

## 环境里隔离的内容

每个环境都维护自己的一套资源：

- `.claude/CLAUDE.md`
- `.claude/settings.json`
- `.claude/skills/`
- `.claude/agents/`
- `.claude/commands/`
- `.claude/rules/`
- `.claude.json`

因此它不只是“切一个 API Key”，而是真正把整套 Claude 工作流拆开。

## Web 管理界面

运行：

```bash
ccv web
```

就可以在浏览器里管理：

- 环境列表
- `LLM`
- `MCP`
- `Skills`
- `Env Vars`
- `CLAUDE.md`
- `Agents`
- `Commands`
- `Rules`
- 导入 / 导出

## 迁移与恢复

如果你已经调好一套环境，可以直接导出：

```bash
ccv export coding
```

在另一台机器恢复：

```bash
ccv import coding-20260322-150405.tar.gz
```

这点对于重装系统、换电脑、或者在多台机器之间同步工作流尤其有用。

## 开发

构建前端并产出可嵌入静态资源：

```bash
make frontend
```

构建全部平台二进制：

```bash
make all
```

构建 npm 分发目录：

```bash
make npm-packages
```
