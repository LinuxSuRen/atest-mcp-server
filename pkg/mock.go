package pkg

import (
	"context"
	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
)

type MockStartRequest struct {
	Prefix string `json:"prefix" jsonschema:"the prefix of mock server, default is /mock"`
	Config string `json:"mockConfig" jsonschema:"the mock config content in YAML format"`
	Port   int    `json:"serverPort" jsonschema:"the port of the mock server, default is 9080"`
}

type MockServer interface {
	Start(ctx context.Context, request *mcp.CallToolRequest, args MockStartRequest) (
		result *mcp.CallToolResult, a any, err error)
	GetConfig(ctx context.Context, request *mcp.CallToolRequest, args any) (
		result *mcp.CallToolResult, a any, err error)
}

type remoteMockServer struct {
	Address string
}

func NewRemoteMockServer(address string) MockServer {
	return &remoteMockServer{
		Address: address,
	}
}

func (r *remoteMockServer) Start(ctx context.Context, request *mcp.CallToolRequest, args MockStartRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewMockClient(conn)

		mockConfig := &server.MockConfig{
			Prefix: args.Prefix,
			Config: args.Config,
			Port:   int32(args.Port),
		}

		_, err = runner.Reload(ctx, mockConfig)
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "success"},
				},
			}
		} else {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}
		}
	}
	return
}

func (r *remoteMockServer) GetConfig(ctx context.Context, request *mcp.CallToolRequest, args any) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewMockClient(conn)

		var config *server.MockConfig
		config, err = runner.GetConfig(ctx, &server.Empty{})
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: config.Config},
				},
			}
		} else {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
				},
			}
		}
	}
	return
}
