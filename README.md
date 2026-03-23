# cc-venv
![GitHub License](https://img.shields.io/github/license/timerzz/cc-venv)
![Static Badge](https://img.shields.io/badge/github-repo-blue%3Flogo%3Dgithub?color=blue)


English | [简体中文](./README_ZH.md)

`cc-venv` is a named environment manager for Claude Code.

The goal is simple: keep different workflows in separate Claude environments, switch when needed, and leave your default Claude Code setup untouched.

## Why cc-venv

Many people do not use Claude Code for just one kind of work.

Sometimes you want a coding workflow:

- coding-focused `MCP` servers
- engineering or frontend-focused `Skills`
- project-oriented `Rules`

At other times you want a writing workflow:

- a different set of writing, translation, or summarization instructions
- different `MCP` tools, or fewer tools
- complete separation from the coding setup

If all of that lives in your default Claude Code configuration, it usually turns into this:

- mixed `Skills`
- `MCP` settings affecting each other
- `Rules` and `CLAUDE.md` instructions bleeding across workflows
- hesitation to experiment because it may break your default setup

`cc-venv` solves that by giving you isolated named environments:

- each environment has its own `CLAUDE.md`
- its own environment variables
- its own `MCP`
- its own `Skills`
- its own `Agents / Commands / Rules`

That makes it easy to maintain separate environments such as:

- `coding`
- `writing`
- `research`
- `prod-safe`

Switching environments does not affect your default Claude Code configuration.

## Compared with cc-mirror

[`cc-mirror`](https://github.com/numman-ali/cc-mirror) and `cc-venv` both provide Claude Code isolation, but they optimize for different things.

| Capability | `cc-venv` | `cc-mirror` |
| --- | --- | --- |
| Multiple isolated Claude Code environments | ✅ | ✅ |
| Embedded Web UI (`ccv web`) | ✅ | ❌ |
| Full environment import / export | ✅ | ❌ |
| Isolate `CLAUDE.md`, env vars, `Skills`, `MCP`, `Agents`, `Commands`, and `Rules` | ✅ | ✅ |
| Separate Claude Code binary per environment | ❌ | ✅ |
| Provider / prompt-pack / tweak-oriented variant management | ❌ | ✅ |

If your goal is to:

- maintain multiple Claude workflows for different purposes
- switch quickly
- keep the default Claude configuration clean
- move a complete environment to another machine

then `cc-venv` is the closer fit.

## Core capabilities

- create named environments
- enter an environment shell
- run Claude Code inside an environment
- manage environments through a Web UI
- isolate `LLM`, environment variables, `MCP`, `Skills`, `Agents`, `Commands`, and `Rules`
- export a full environment
- import and restore that environment on a new machine

## Quick start

### Option 1: Install globally with npm

```bash
npm install -g @timerzz/ccv
```

After installation:

```bash
ccv list
ccv create coding
ccv web
```

### Option 2: Download a Go binary directly

Download the binary for your platform from GitHub Releases and add it to your `PATH`:

- Linux x64: `ccv-linux-amd64`
- Linux arm64: `ccv-linux-arm64`
- macOS x64: `ccv-darwin-amd64`
- macOS arm64: `ccv-darwin-arm64`
- Windows x64: `ccv-windows-amd64.exe`
- Windows arm64: `ccv-windows-arm64.exe`

Common commands:

```bash
# Create a virtual environment named coding
ccv create coding
# Open the Web UI to manage and configure environments
ccv web
# List all current virtual environments
ccv list
# Run Claude Code inside the coding environment
ccv run coding
# Export the full coding environment, including Skills and MCP configuration
ccv export coding
```

<img width="2485" height="1268" alt="web" src="https://github.com/user-attachments/assets/ba6f3240-dac8-414a-8343-ff8c8e5aa488" />

## What is isolated inside an environment

Each environment maintains its own set of resources:

- `.claude/CLAUDE.md`
- `.claude/settings.json`
- `.claude/skills/`
- `.claude/agents/`
- `.claude/commands/`
- `.claude/rules/`
- `.claude.json`

So this is not just about switching an API key. It isolates a full Claude workflow.

## Web UI

Run:

```bash
ccv web
```

Then manage the following in the browser:

- environment list
- `LLM`
- `MCP`
- `Skills`
- `Env Vars`
- `CLAUDE.md`
- `Agents`
- `Commands`
- `Rules`
- import / export

## Migration and recovery

Once an environment is configured, you can export it:

```bash
ccv export coding
```

Then restore it on another machine:

```bash
ccv import coding-20260322-150405.tar.gz
```

This is especially useful when moving to a new machine, reinstalling a system, or syncing workflows across devices.

## Development

Build the frontend and sync embedded static assets:

```bash
make frontend
```

Build binaries for all platforms:

```bash
make all
```

Prepare npm distribution packages:

```bash
make npm-packages
```
