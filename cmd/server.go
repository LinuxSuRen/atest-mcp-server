package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/linuxsuren/api-testing/docs"
	"github.com/linuxsuren/api-testing/pkg/mock"
	"github.com/linuxsuren/atest-mcp-server/pkg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

type serverOption struct {
	port          int
	runnerAddress string

	mockServer mock.DynamicServer
}

func newServerCommand() *cobra.Command {
	opt := serverOption{}
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run server",
		RunE:  opt.runE,
	}
	cmd.Flags().IntVarP(&opt.port, "port", "p", 7845, "The port to run server")
	cmd.Flags().StringVarP(&opt.runnerAddress, "runner-address", "", "", "The address of the runner")
	return cmd
}

type Args struct {
	Name       string `json:"name" jsonschema:"the name to say hi to"`
	ServerPort int    `json:"serverPort" jsonschema:"the port of the mock server" default:"9080"`
	MockConfig string `json:"mockConfig" jsonschema:"the mock config content in YAML format"`
}

func (o *serverOption) runE(c *cobra.Command, args []string) (err error) {
	opts := &mcp.ServerOptions{
		Instructions:      "ATest Server",
		CompletionHandler: complete,
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name: "atest-mcp-server",
	}, opts)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "start-mock-server",
		Description: "Start a mock server",
	}, func(ctx context.Context, request *mcp.CallToolRequest, args Args) (*mcp.CallToolResult, any, error) {
		if o.mockServer == nil {
			o.mockServer = mock.NewInMemoryServer(ctx, args.ServerPort)
		}

		go func() {
			reader := mock.NewInMemoryReader(args.MockConfig)

			err := o.mockServer.Start(reader, "/")
			if err != nil {
				fmt.Println("Failed to start mock server: ", err)
			}
		}()

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Mock Server started on port: " + o.mockServer.GetPort()},
			},
		}, nil, nil
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "stop-mock-server",
		Description: "Stop the mock server",
	}, func(ctx context.Context, request *mcp.CallToolRequest, args Args) (*mcp.CallToolResult, any, error) {
		if o.mockServer != nil {
			err := o.mockServer.Stop()
			msg := "Mock Server stopped"
			if err != nil {
				msg = "Failed to stop mock server: " + err.Error()
			}
			o.mockServer = nil
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: msg},
				},
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Mock Server not started"},
			},
		}, nil, nil
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get-mock-config-schema",
		Description: "Get the mock config schema",
	}, func(ctx context.Context, request *mcp.CallToolRequest, args Args) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: docs.MockSchema},
			},
		}, nil, err
	})

	if o.runnerAddress != "" {
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
			Description: "Create a test suite for HTTP testing",
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
			Name:        "delete-test-case",
			Description: "Delete a test case for HTTP testing",
		}, runner.DeleteTestCase)
	}

	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return server
	}, nil)
	c.Println("Starting server on port: ", o.port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", o.port), handler)
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
