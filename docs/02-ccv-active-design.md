# ccv active 设计文档

## 1. 命令定位

`ccv active <name>` 的作用是：

- 进入指定命名环境的终端
- 为该终端会话注入环境变量
- 不自动执行 Claude Code

它的目标是给用户一个“已经切换到该环境”的交互式 shell。

---

## 2. 用户心智模型

用户执行：

```bash
ccv active team-base
```

预期行为是：

1. `ccv` 加载 `~/.ccv/envs/team-base/`
2. 以该环境为运行上下文准备变量
3. 启动一个新的交互式 shell
4. 用户随后可以在该 shell 中手动执行：
   - `claude`
   - `claude /xxx`
   - `git`
   - 其他任意命令

所以它更像：

- “进入一个已激活环境的 shell”

而不是：

- “执行一次命令后退出”

---

## 3. 输入与输出

命令格式：

```bash
ccv active <name>
```

成功时：

- 直接进入交互式 shell

失败时：

- 环境不存在
- shell 不可启动
- 环境目录严重损坏

---

## 4. 环境发现与校验

`active` 不做自动发现。

它只做显式加载：

```text
~/.ccv/envs/<name>/
```

第一版建议最少校验：

- 环境根目录存在

可选校验：

- `ccv.json` 是否存在
- `.claude/` 是否存在

设计原则是：

- 环境目录本体是真实状态来源
- `active` 不应做过重校验，以免阻断进入环境排查问题

---

## 5. 运行时环境变量

`active` 启动 shell 前，应基于两类来源构造进程环境变量。

### 5.1 来源一：父进程现有环境

先继承当前终端已有的环境变量。

### 5.2 来源二：环境定义变量

再加载该命名环境定义的变量。

这些变量用于承载：

- LLM / provider 相关变量
- 用户自定义环境变量
- MCP 或插件依赖变量

这里强调的是“有一份由 `ccv` 管理的环境变量配置”，而不是要求它必须属于 Claude Code 原生目录结构的一部分。该配置是 `ccv` 的运行时配置，不属于导出时要对齐的 Claude 原生资源目录。

### 5.3 来源三：ccv 强制注入变量

在前两层合并后，再由 `ccv` 覆盖写入：

- `CLAUDE_CONFIG_DIR`
- `CCV_ENV_NAME`
- `CCV_ENV_ROOT`
- `CCV_ACTIVE=1`

其中：

- `CLAUDE_CONFIG_DIR` 应指向该命名环境的 Claude 配置根目录
- 在当前原生兼容模型下，它应对应环境根目录，而不是旧的 `config/` 子目录

例如：

```text
CLAUDE_CONFIG_DIR=/home/user/.ccv/envs/team-base
CCV_ENV_NAME=team-base
CCV_ENV_ROOT=/home/user/.ccv/envs/team-base
CCV_ACTIVE=1
```

合并优先级建议为：

1. 父进程环境
2. 环境定义变量
3. `ccv` 强制变量

---

## 6. 注入机制

`active` 的“注入”不是修改当前父 shell。

它的实现机制应是：

- 创建一个新的子 shell 进程
- 把合并后的环境变量传给这个子进程
- 把 `stdin/stdout/stderr` 直接绑定到当前终端

这意味着：

- 用户进入的是一个新的 shell
- 用户输入会直接传给这个 shell
- 该 shell 中启动的所有子进程都会继承这些变量
- 退出该 shell 后，不会污染原来的父终端会话

---

## 7. shell 启动方式

第一版建议：

- 优先读取 `$SHELL`
- 为空时回退 `/bin/sh`

行为原则：

- 启动交互式 shell
- 继承当前标准输入输出
- 不做复杂的 shell 类型适配

更细的平台兼容策略可以后续单独补充，但不属于第一版阻塞项。

---

## 8. 当前工作目录

`active` 第一版采用：

- **保留当前工作目录**

也就是：

- 只切环境变量
- 不切 `cwd`

原因：

- 更符合“我在当前项目里临时切换 Claude 环境”的使用方式
- 不会打断当前项目上下文

如果用户需要查看环境目录，可以自行：

```bash
cd "$CCV_ENV_ROOT"
```

---

## 9. 提示信息

为了让用户明确知道已经进入哪个环境，第一版建议：

- 在进入 shell 前打印一行提示

例如：

```text
[ccv] active environment: team-base
```

第一版不强制修改 prompt。

原因：

- 简单
- 跨 shell 兼容性更好
- 不干扰用户已有 prompt 配置

---

## 10. 与 `run` 的关系

`active` 和 `run` 应共享同一套执行上下文准备逻辑。

两者共同部分：

- 加载命名环境
- 合并环境变量
- 注入 `CLAUDE_CONFIG_DIR` 和 `CCV_*`
- 保留当前工作目录

两者差异只在最后一步：

- `active` 启动交互式 shell
- `run` 直接启动 `claude`

因此实现上应尽量复用同一套“执行上下文准备”逻辑，避免行为漂移。
