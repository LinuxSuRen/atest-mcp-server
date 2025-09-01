package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"

	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/linuxsuren/atest-mcp-server/pkg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type serverOption struct {
	port          int
	runnerAddress string
	mode          string

	mockServer mock.DynamicServer
}

func newServerCommand() *cobra.Command {
	opt := serverOption{}
	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Run server",
		PreRunE: opt.preRunE,
		RunE:    opt.runE,
	}
	cmd.Flags().IntVarP(&opt.port, "port", "p", 7845, "The port to run server")
	cmd.Flags().StringVarP(&opt.runnerAddress, "runner-address", "", "", "The address of the runner")
	cmd.Flags().StringVarP(&opt.mode, "mode", "m", "http", "The mode: http, stdio or sse")
	return cmd
}

type Args struct {
	Name       string `json:"name" jsonschema:"the name to say hi to"`
	ServerPort int    `json:"serverPort" jsonschema:"the port of the mock server" default:"9080"`
	MockConfig string `json:"mockConfig" jsonschema:"the mock config content in YAML format"`
}

func (o *serverOption) preRunE(c *cobra.Command, args []string) (err error) {
	if o.runnerAddress == "" {
		err = fmt.Errorf("the runner-address is required")
	}
	return
}

//go:embed data/mainPrompt.txt
var mainPrompt string

func (o *serverOption) runE(c *cobra.Command, args []string) (err error) {
	opts := &mcp.ServerOptions{
		Instructions:      "ATest Server",
		CompletionHandler: complete,
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:  "atest-mcp-server",
		Title: "api-testing (aka atest) MCP Server",
	}, opts)

	server.AddPrompt(&mcp.Prompt{
		Name:        "create-test-case",
		Description: "Create a test case for HTTP testing",
	}, func(ctx context.Context, request *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Messages: []*mcp.PromptMessage{
				{
					Role: "user",
					Content: &mcp.TextContent{
						Text: mainPrompt,
					},
				},
			},
		}, nil
	})
	server.AddResource(&mcp.Resource{
		Name:        "atest-knowledge",
		Title:       "atest Knowledge",
		Description: "The knowledge of api-tesing (aka atest)",
		URI:         "https://linuxsuren.github.io/api-testing",
	}, func(ctx context.Context, request *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      "https://linuxsuren.github.io/api-testing/latest/tasks/template/",
				MIMEType: "text/html",
				Text:     "template functions",
			}, {
				URI:      "https://linuxsuren.github.io/api-testing/latest/tasks/verify/",
				MIMEType: "text/html",
				Text:     "verify functions",
			}, {
				URI:      "https://linuxsuren.github.io/api-testing/latest/tasks/mock/",
				MIMEType: "text/html",
				Text:     "mock server introduction",
			}, {
				URI:      "https://linuxsuren.github.io/api-testing/latest/tasks/extension/",
				MIMEType: "text/html",
				Text:     "extension introduction",
			}},
		}, nil
	})

	mockServer := pkg.NewRemoteMockServer(o.runnerAddress)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "start-mock-server",
		Description: "Start a mock server",
	}, mockServer.Start)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-mock-config",
		Description: "Get the mock config as YAML format",
	}, mockServer.GetConfig)

	runner := pkg.NewRunner(o.runnerAddress)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "run",
		Description: "Run a test case",
	}, runner.Run)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-suites",
		Description: "Get all test suites",
	}, runner.GetSuites)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-test-suite",
		Description: "Create a test suite for HTTP testing. Test suite is a collection of test cases.",
	}, runner.CreateTestSuite)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create-test-case",
		Description: "Create a test case for HTTP testing",
	}, runner.CreateTestCase)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-test-suite",
		Description: "Get a test suite for HTTP testing",
	}, runner.GetTestSuite)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-test-suite",
		Description: "Delete a test suite for HTTP testing",
	}, runner.DeleteTestSuite)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list-test-case",
		Description: "List all test cases",
	}, runner.ListTestCase)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-test-case",
		Description: "Get a test case for HTTP testing",
	}, runner.GetTestCase)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "run-test-case",
		Description: "Run a test case",
	}, runner.RunTestCase)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update-test-suite",
		Description: "Update a test suite for HTTP testing",
	}, runner.UpdateTestSuite)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "update-test-case",
		Description: "Update a test case for HTTP testing",
	}, runner.UpdateTestCase)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-suggested-apis",
		Description: "Get suggested APIs from swagger for HTTP testing",
	}, runner.GetSuggestedAPIs)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete-test-case",
		Description: "Delete a test case for HTTP testing",
	}, runner.DeleteTestCase)

	switch o.mode {
	case "sse":
		handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			return server
		})
		c.Println("Starting SSE server on port:", o.port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", o.port), handler)
	case "stdio":
		err = server.Run(c.Context(), &mcp.StdioTransport{})
	case "http":
		fallthrough
	default:
		handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
			return server
		}, nil)
		c.Println("Starting HTTP server on port:", o.port)
		err = http.ListenAndServe(fmt.Sprintf(":%d", o.port), handler)
	}
	return
}

func complete(ctx context.Context, req *mcp.CompleteRequest) (*mcp.CompleteResult, error) {
	return &mcp.CompleteResult{
		Completion: mcp.CompletionResultDetails{
			Total:  1,
			Values: []string{req.Params.Argument.Value + "x"},
		},
	}, nil
}
