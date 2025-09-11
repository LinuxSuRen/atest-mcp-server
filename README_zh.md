这是 atest 的 MCP Server

## 开始
请用以下命令启动MCP服务器： 一个测试 gRPC端口。

```
atest-store-mcp server --runner-address 127.0.0.1:64385
```

或在 Docker 中启动

```
docker run -p 7845:7845 ghcr.io/linuxsuren/atest-mcp-server --runner-address 127.0.0.1:64385
```

或从 npx 开始

```
npx atest-mcp-server-launcher@latest server --mode=stdio --runner-address=localhost:64385
```

你也可以通过以下方式设置MCP服务器模式： --mode 旗帜。

```
atest-store-mcp server --runner-address 127.0.0.1:64385 --mode=[sse|stdio]
```

##

MCP服务器

```json
{
  "mcpServers": {
    "atest": {
      "name": "atest",
      "type": "streamableHttp",
      "description": "The MCP server of atest",
      "isActive": true,
      "baseUrl": "http://localhost:7845",
      "disabledAutoApproveTools": []
    }
  }
}
```

或作为标准输入输出模式:

```json
{
  "mcpServers": {
    "atest": {
      "name": "atest-mcp-stdio",
      "type": "stdio",
      "description": "",
      "isActive": true,
      "command": "atest-store-mcp",
      "args": [
        "server",
        "-m=stdio",
        "--runner-address=localhost:64385"
      ]
    }
  }
}
```

## 如何构建
你可以使用以下命令构建二进制文件:

```
make build
```

你可以构建 .dxt 使用以下命令打包:

```shell
npm install -g @anthropic-ai/dxt
dxt pack
```
