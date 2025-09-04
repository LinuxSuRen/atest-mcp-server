package pkg

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type starter struct {
}

type Starter interface {
	Start(ctx context.Context, request *mcp.CallToolRequest, args any) (
		result *mcp.CallToolResult, a any, err error)
}

func NewStarter() Starter {
	return &starter{}
}

func (s *starter) Start(ctx context.Context, request *mcp.CallToolRequest, args any) (
	result *mcp.CallToolResult, a any, err error) {
	var startErr error
	switch runtime.GOOS {
	case "windows":
		startErr = exec.Command("cmd", "/C", "start atest-desktop").Run()
	case "darwin":
		startErr = exec.Command("open", "-a", "atest-desktop").Run()
	case "linux":
		startErr = exec.Command("xdg-open", "atest-desktop").Run()
	default:
		startErr = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if startErr == nil {
		result = &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "atest started successfully, please check the app"},
			},
		}
	} else {
		result = &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: startErr.Error()},
			},
		}
	}
	return
}
