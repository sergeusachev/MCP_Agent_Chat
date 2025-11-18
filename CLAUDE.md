# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Model Context Protocol (MCP) client-server implementation using the official Go SDK (`github.com/modelcontextprotocol/go-sdk`). The project demonstrates bidirectional communication between an MCP client and server over stdio transport.

## Project Structure

```
MCP_Example/
├── mcp_server/
│   └── server_example.go    # MCP server with "greet" tool
├── mcp_client/
│   └── client_example.go    # MCP client that calls the server
├── go.mod                   # Module dependencies
└── go.sum                   # Dependency checksums
```

## Build Commands

**Build the server:**
```bash
cd mcp_server
go build -o myserver server_example.go
```

**Build the client:**
```bash
cd mcp_client
go build -o myclient client_example.go
```

**Run the server (standalone):**
```bash
cd mcp_server
./myserver
```

**Run the client (spawns and connects to server):**
```bash
cd mcp_client
./myclient
```

## Architecture

### MCP Server (`mcp_server/server_example.go`)

- **Server Name**: "greeter" (v1.0.0)
- **Transport**: StdioTransport (communicates via stdin/stdout)
- **Tool**: "greet" - Accepts first name and second name, returns a personalized greeting

The server uses typed input/output with automatic JSON schema generation:
- **Input type**: `Input{Name, SecondName}`
- **Output type**: `Output{Greeting, MetaData}`
- **Handler**: `SayHi()` function concatenates names and returns formatted greeting

### MCP Client (`mcp_client/client_example.go`)

- **Client Name**: "mcp-client" (v1.0.0)
- **Transport**: CommandTransport (spawns server as subprocess)
- **Server Binary Path**: `../mcp_server/myserver` (relative to client executable)

The client:
1. Creates a CommandTransport that executes the server binary
2. Establishes a session with the server
3. Calls the "greet" tool with name parameters
4. Logs the response content

## Key MCP Concepts in This Codebase

### Transport Layer
The project uses stdio-based transport where:
- Server runs with `StdioTransport{}` listening on stdin/stdout
- Client uses `CommandTransport` to spawn the server and communicate via pipes

### Tool Definition
Tools are defined with typed handlers:
```go
mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
```

The SDK automatically:
- Generates JSON schemas from Go struct tags
- Validates input arguments
- Serializes output to JSON

### Type Safety
Input/output types use `jsonschema` struct tags for schema metadata:
```go
type Input struct {
    Name string `json:"name" jsonschema:"the name of the person to greet"`
}
```

## Dependencies

- `github.com/modelcontextprotocol/go-sdk v1.1.0` - Official MCP SDK
- `github.com/google/jsonschema-go` - JSON schema generation
- `github.com/yosida95/uritemplate/v3` - URI template parsing
- `golang.org/x/oauth2` - OAuth2 support

## Development Workflow

1. **Modify the server**: Edit `mcp_server/server_example.go` to add/modify tools
2. **Rebuild the server**: `cd mcp_server && go build -o myserver server_example.go`
3. **Modify the client**: Edit `mcp_client/client_example.go` to call tools with different parameters
4. **Rebuild the client**: `cd mcp_client && go build -o myclient client_example.go`
5. **Test**: Run `./myclient` from the `mcp_client` directory

## Important Notes

- The client expects the server binary at `../mcp_server/myserver` - ensure this path is correct
- Both client and server must be built before the client can run successfully
- The server blocks on stdio, waiting for client requests
- Communication uses the MCP JSON-RPC protocol over stdio
