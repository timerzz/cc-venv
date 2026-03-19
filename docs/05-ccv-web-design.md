# ccv web 设计文档

## 1. 目标

`ccv web` 的作用是：

- 提供本机命名环境的图形化管理入口
- 复用 CLI 已有能力，而不是重新定义一套环境模型
- 最终以单二进制方式分发

第一版重点是管理，不是远程协作或多用户系统。

---

## 2. 技术方案

第一版建议技术栈：

- 后端：`Gin`
- 前端：`Vite + React + shadcn/ui`
- 分发：前端构建产物通过 Go `embed` 嵌入后端

这套方案的目标是：

- 保持前后端开发体验清晰
- 最终运行时不依赖独立前端服务
- 继续以文件系统作为事实来源

第一版不引入数据库。

原因：

- 当前系统的事实来源已经是 `~/.ccv/envs/`
- 引入数据库会带来额外同步问题
- Web 只是另一层入口，不应变成新的状态源

---

## 3. 系统边界

`ccv web` 只管理本机命名环境：

```text
~/.ccv/envs/
```

不管理：

- 项目级 Claude 资源
- 远程主机环境
- 多用户权限模型

Web 的后端应直接复用：

- `internal/env`
- `internal/config`
- `internal/exporter`
- `internal/importer`

不能在 Web 层重复实现环境读写逻辑。

---

## 4. 前后端结构

建议结构：

```text
cc-venv/
├── internal/
│   └── webui/
│       ├── server.go                  # HTTP服务入口、embed嵌入
│       ├── routes.go                  # 路由注册
│       ├── handlers_env.go            # 环境CRUD handlers
│       ├── handlers_llm.go            # LLM配置 handlers
│       ├── handlers_mcp.go            # MCP管理 handlers
│       ├── handlers_skill.go          # Skills管理 handlers
│       ├── handlers_import_export.go  # 导入导出 handlers
│       └── response.go                # Response[T] 结构体
└── web/                               # 前端项目
    ├── index.html
    ├── package.json
    ├── vite.config.ts
    └── src/
        ├── main.tsx
        ├── app.tsx
        ├── lib/api.ts
        ├── pages/
        └── components/
```

其中：

- `internal/webui` 负责 HTTP 服务和 API
- `web/` 负责前端页面和构建

---

## 5. 统一响应结构

```go
type Response[T any] struct {
    Code int    `json:"code"`
    Data T      `json:"data,omitempty"`
    Msg  string `json:"msg,omitempty"`
}
```

- `code = 0`: 成功
- `code != 0`: 失败，错误信息在 `msg`

示例：

```json
// 成功
{"code": 0, "data": {...}}

// 失败
{"code": 1, "msg": "environment not found"}
```

---

## 6. API 设计

### 6.1 环境基础

```
GET    /api/envs                      # 列出所有环境
POST   /api/envs                      # 创建环境
GET    /api/envs/:name                # 获取环境详情
PUT    /api/envs/:name                # 修改环境基础信息
DELETE /api/envs/:name                # 删除环境
POST   /api/envs/:name/export         # 导出环境
POST   /api/envs/import               # 导入环境
```

#### GET /api/envs

响应：

```json
{
  "code": 0,
  "data": {
    "envs": [
      {
        "name": "my-env",
        "path": "/home/user/.ccv/envs/my-env",
        "resources": {
          "skills": 2,
          "agents": 1,
          "commands": 0,
          "rules": 0,
          "hooks": 0,
          "mcpServers": 1
        }
      }
    ]
  }
}
```

#### POST /api/envs

请求：

```json
{
  "name": "new-env"
}
```

#### GET /api/envs/:name

响应：

```json
{
  "code": 0,
  "data": {
    "name": "my-env",
    "path": "/home/user/.ccv/envs/my-env",
    "claudeMd": "# My env memo...",
    "settings": {...},
    "envVars": {...},
    "mcpServers": {...},
    "resources": {
      "skills": ["skill1", "skill2"],
      "agents": ["agent1"],
      "commands": [],
      "rules": [],
      "hooks": [],
      "plugins": []
    }
  }
}
```

#### PUT /api/envs/:name

请求：

```json
{
  "name": "new-name",              // 可选，重命名环境
  "claudeMd": "# My env memo...",  // 可选，CLAUDE.md 内容
  "envVars": {                     // 可选，环境变量
    "CUSTOM_VAR": "value"
  }
}
```

### 6.2 LLM 供应商管理

```
GET    /api/envs/:name/llm            # 获取LLM配置
PUT    /api/envs/:name/llm            # 更新LLM配置
GET    /api/llm/providers             # 获取支持的供应商列表
```

#### GET /api/envs/:name/llm

响应：

```json
{
  "code": 0,
  "data": {
    "provider": "anthropic",
    "apiKey": "sk-***",              // 脱敏显示
    "baseUrl": "",
    "models": {
      "default": "claude-sonnet-4-6",
      "sonnet": "claude-sonnet-4-6",
      "opus": "claude-opus-4-6",
      "haiku": "claude-haiku-4-5"
    }
  }
}
```

#### PUT /api/envs/:name/llm

请求：

```json
{
  "provider": "anthropic",
  "apiKey": "sk-xxx",
  "baseUrl": "https://api.anthropic.com",
  "models": {
    "default": "claude-sonnet-4-6",
    "sonnet": "claude-sonnet-4-6",
    "opus": "claude-opus-4-6",
    "haiku": "claude-haiku-4-5"
  }
}
```

#### GET /api/llm/providers

响应：

```json
{
  "code": 0,
  "data": {
    "providers": [
      {
        "id": "anthropic",
        "name": "Anthropic",
        "baseUrl": "https://api.anthropic.com",
        "models": ["claude-sonnet-4-6", "claude-opus-4-6", "claude-haiku-4-5"]
      },
      {
        "id": "openai",
        "name": "OpenAI",
        "baseUrl": "https://api.openai.com",
        "models": ["gpt-4", "gpt-4o", "gpt-3.5-turbo"]
      }
    ]
  }
}
```

### 6.3 MCP 管理

```
GET    /api/envs/:name/mcp            # 获取MCP服务器列表
POST   /api/envs/:name/mcp            # 添加MCP服务器
PUT    /api/envs/:name/mcp/:server    # 更新指定MCP服务器
DELETE /api/envs/:name/mcp/:server    # 删除指定MCP服务器
```

#### GET /api/envs/:name/mcp

响应：

```json
{
  "code": 0,
  "data": {
    "servers": {
      "filesystem": {
        "command": "npx",
        "args": ["-y", "@anthropic-ai/mcp-server-filesystem", "/path"],
        "env": {}
      },
      "my-server": {
        "command": "node",
        "args": ["server.js"],
        "env": {"API_KEY": "xxx"}
      }
    }
  }
}
```

#### POST /api/envs/:name/mcp

请求：

```json
{
  "name": "filesystem",
  "config": {
    "command": "npx",
    "args": ["-y", "@anthropic-ai/mcp-server-filesystem", "/path"],
    "env": {}
  }
}
```

#### PUT /api/envs/:name/mcp/:server

请求：

```json
{
  "config": {
    "command": "npx",
    "args": ["-y", "@anthropic-ai/mcp-server-filesystem", "/new-path"],
    "env": {}
  }
}
```

### 6.4 Skills 管理

```
GET    /api/envs/:name/skills                 # 获取skills列表
POST   /api/envs/:name/skills                 # 添加skill（zip上传或URL下载）
DELETE /api/envs/:name/skills/:skill          # 删除skill
```

#### GET /api/envs/:name/skills

响应：

```json
{
  "code": 0,
  "data": {
    "skills": [
      {
        "name": "skill1",
        "path": ".claude/skills/skill1"
      },
      {
        "name": "skill2",
        "path": ".claude/skills/skill2"
      }
    ]
  }
}
```

#### POST /api/envs/:name/skills

支持两种方式：

**方式1: URL下载**

请求 (application/json):

```json
{
  "url": "https://github.com/user/skill-repo/archive/main.zip"
}
```

**方式2: 文件上传**

请求 (multipart/form-data):

```
file: skill.zip
```

响应：

```json
{
  "code": 0,
  "data": {
    "name": "skill1",
    "path": ".claude/skills/skill1"
  }
}
```

### 6.5 导入导出

```
POST   /api/envs/:name/export         # 导出环境
POST   /api/envs/import               # 导入环境
```

#### POST /api/envs/:name/export

响应：

```json
{
  "code": 0,
  "data": {
    "downloadUrl": "/api/downloads/xxx.tar.gz"
  }
}
```

#### POST /api/envs/import

请求 (multipart/form-data):

```
file: env.tar.gz
```

---

## 7. 静态资源嵌入

前端构建产物嵌入到 Go 二进制中：

```go
//go:embed all:web/dist
var staticFS embed.FS
```

路由策略：

- `/api/*` 路由到 Gin handlers
- 其他路由回退到 `index.html` (SPA)

---

## 8. 开发模式

开发阶段前后端分离运行：

- 前端: `npm run dev` (Vite dev server, 代理 `/api` 到后端)
- 后端: `go run ./cmd/ccv web --dev` (仅启动 API 服务)

生产模式：

- 前端构建: `npm run build` -> `web/dist/`
- 后端编译: `go build ./cmd/ccv`
- 运行: `ccv web` (单二进制，嵌入前端)

---

## 9. 命令行参数

```bash
ccv web              # 启动服务器并打开浏览器
ccv web --port 8080  # 指定端口（默认 3000）
ccv web --no-open    # 不自动打开浏览器
ccv web --dev        # 开发模式（不嵌入前端，仅API）
```

---

## 10. 实现边界

`webui` 包的职责：

- 启动 HTTP 服务
- 注册路由
- 调用 `internal/env`、`internal/config`、`internal/exporter`、`internal/importer`
- 处理 HTTP 请求/响应

约束：

- 不在 `webui` 层重复实现环境读写逻辑
- 继续以文件系统作为事实来源
- 第一版不引入数据库

---

## 11. 实现顺序

1. **Phase 1: 核心 API**
   - 实现 `Response[T]` 结构体
   - 实现环境 CRUD API
   - 复用 `internal/env` 模块

2. **Phase 2: LLM & MCP**
   - 实现 LLM 配置 API
   - 实现 MCP 管理 API
   - 复用 `internal/config` 模块

3. **Phase 3: Skills & 导入导出**
   - 实现 Skills 管理 API
   - 实现导入导出 API
   - 复用 `internal/exporter`、`internal/importer`

4. **Phase 4: 前端集成**
   - 实现静态资源嵌入
   - 实现浏览器自动打开
   - 前端页面开发

---

## 12. 与现有设计文档的关系

- 系统边界以 `00-ccv-system-design.md` 为准
- 环境列表口径应与 `04-ccv-list-design.md` 一致
- 环境详情中的扫描口径应与 `03-ccv-scan-design.md` 一致

Web 不应引入独立的环境模型、扫描口径或状态源。
