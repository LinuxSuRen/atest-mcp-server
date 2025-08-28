package pkg

import (
	"context"
	"encoding/json"

	"github.com/linuxsuren/api-testing/pkg/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/grpc"
)

type Runner interface {
	Run(ctx context.Context, request *mcp.CallToolRequest, args RunRequest) (result *mcp.CallToolResult, a any, err error)
	GetSuites(ctx context.Context, request *mcp.CallToolRequest, args any) (result *mcp.CallToolResult, a any, err error)
	CreateTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (*mcp.CallToolResult, any, error)
	CreateTestCase(ctx context.Context, request *mcp.CallToolRequest, args CreateTestCaseRequest) (*mcp.CallToolResult, any, error)
	GetTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (result *mcp.CallToolResult, a any, err error)
}

type gRPCRunner struct {
	Address string
}

func NewRunner(address string) Runner {
	return &gRPCRunner{
		Address: address,
	}
}

type RunRequest struct {
	SuiteName string `json:"suiteName" jsonschema:"the name of test suite"`
	CaseName  string `json:"caseName" jsonschema:"the name of test case"`
}

func (r *gRPCRunner) Run(ctx context.Context, request *mcp.CallToolRequest, args RunRequest) (result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		runReq := &server.TestTask{
			CaseName: args.CaseName,
		}

		var reply *server.TestResult
		reply, err = runner.Run(ctx, runReq)
		if err == nil {
			data := reply.TestCaseResult

			dataAsStr, _ := json.Marshal(data)

			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(dataAsStr)},
				},
			}
		}
	}
	return
}

func (r *gRPCRunner) GetSuites(ctx context.Context, request *mcp.CallToolRequest, args any) (result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		var reply *server.Suites
		reply, err = runner.GetSuites(ctx, &server.Empty{})
		if err == nil {
			data := reply.Data

			dataAsStr, _ := json.Marshal(data)

			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: string(dataAsStr)},
				},
			}
		}
	}
	return
}

type TestSuiteIndentityRequest struct {
	Name string `json:"name" jsonschema:"the name of test suite"`
	API  string `json:"api" jsonschema:"the API path for test suite"`
}

func (r *gRPCRunner) CreateTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
			Kind: "http",
		}

		var reply *server.HelloReply
		reply, err = runner.CreateTestSuite(ctx, suite)
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: reply.Message},
				},
			}
		}
	}
	return
}

type CreateTestCaseRequest struct {
	SuiteName   string            `json:"suiteName" jsonschema:"the name of test suite"`
	CaseName    string            `json:"caseName" jsonschema:"the name of test case"`
	API         string            `json:"api" jsonschema:"the API path for test case"`
	Method      string            `json:"method" jsonschema:"the HTTP method for test case"`
	Body        string            `json:"body" jsonschema:"the body for test case"`
	Headers     map[string]string `json:"headers" jsonschema:"the headers for test case"`
	QueryParams map[string]string `json:"queryParams" jsonschema:"the query params for test case"`
}

func (r *gRPCRunner) CreateTestCase(ctx context.Context, request *mcp.CallToolRequest, args CreateTestCaseRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		testCase := &server.TestCaseWithSuite{
			SuiteName: args.SuiteName,
			Data: &server.TestCase{
				Name:      args.CaseName,
				SuiteName: args.SuiteName,
				Request: &server.Request{
					Api:    args.API,
					Method: args.Method,
					Body:   args.Body,
				},
			},
		}

		var reply *server.HelloReply
		reply, err = runner.CreateTestCase(ctx, testCase)
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: reply.Message},
				},
			}
		}
	}
	return
}

func (r *gRPCRunner) GetTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
			Kind: "http",
		}

		var reply *server.TestSuite
		reply, err = runner.GetTestSuite(ctx, suite)
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: reply.String()},
				},
			}
		}
	}
	return
}
