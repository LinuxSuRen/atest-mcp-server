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
	UpdateTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteArgs) (
		result *mcp.CallToolResult, a any, err error)
	ListTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
	RunTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
	GetTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
	DeleteTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
	UpdateTestCase(ctx context.Context, request *mcp.CallToolRequest, args CreateTestCaseRequest) (
		result *mcp.CallToolResult, a any, err error)
	GetSuggestedAPIs(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
	DeleteTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
		result *mcp.CallToolResult, a any, err error)
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
	SuiteName string `json:"suiteName" jsonschema:"the name of test suite" mcp:"the name of test suite"`
	CaseName  string `json:"caseName" jsonschema:"the name of test case" mcp:"the name of test case"`
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
	Name string `json:"name" jsonschema:"the name of test suite" mcp:"the name of test suite"`
	API  string `json:"api" jsonschema:"the API path for test suite" mcp:"the API path for test suite"`
	Kind string `json:"kind" jsonschema:"the kind of test suite, such as swagger" mcp:"the kind of test suite, such as swagger"`
}

func (r *gRPCRunner) CreateTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
			Kind: args.Kind,
		}

		var reply *server.HelloReply
		reply, err = runner.CreateTestSuite(ctx, suite)
		if err == nil {
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "success"},
				},
			}
		} else {
			if reply == nil {
				reply = &server.HelloReply{}
			}
			result = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: err.Error()},
					&mcp.TextContent{Text: reply.Message},
				},
			}
		}
	}
	return
}

type CreateTestCaseRequest struct {
	SuiteName     string            `json:"suiteName" jsonschema:"the name of test suite" mcp:"the name of test suite"`
	CaseName      string            `json:"caseName" jsonschema:"the name of test case" mcp:"the name of test case"`
	API           string            `json:"api" jsonschema:"the API path for test case" mcp:"the API path for test case"`
	Method        string            `json:"method" jsonschema:"the HTTP method for test case" mcp:"the HTTP method for test case"`
	Body          string            `json:"body" jsonschema:"HTTP request payload body for test case" mcp:"HTTP request payload body for test case"`
	Headers       map[string]string `json:"headers" jsonschema:"HTTP request headers for test case" mcp:"HTTP request headers for test case"`
	QueryParams   map[string]string `json:"queryParams" jsonschema:"HTTP request query params for test case" mcp:"HTTP request query params for test case"`
	Cookies       map[string]string `json:"cookies" jsonschema:"HTTP request cookies for test case" mcp:"HTTP request cookies for test case"`
	FormParams    map[string]string `json:"formParams" jsonschema:"HTTP request form params for test case" mcp:"HTTP request form params for test case"`
	ExpectStatus  int32             `json:"expectStatus" jsonschema:"the expected HTTP status code for the HTTP response, such as 200" mcp:"the expected HTTP status code for the HTTP response, such as 200" mcp:"the expected HTTP status code for the HTTP response, such as 200"`
	ExpectBody    string            `json:"expectBody" jsonschema:"the expected HTTP response body for test case" mcp:"the expected HTTP response body for test case"`
	ExpectHeaders map[string]string `json:"expectHeaders" jsonschema:"the expected HTTP response headers for test case" mcp:"the expected HTTP response headers for test case"`
	ExpectSchema  string            `json:"expectSchema" jsonschema:"the expected HTTP response to verify as JSON schema for test case" mcp:"the expected HTTP response to verify as JSON schema for test case"`
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

type TestSuiteArgs struct {
	Name  string   `json:"name" jsonschema:"the name of test suite"`
	API   string   `json:"api" jsonschema:"the API path for test suite"`
	Param []*Pair  `json:"param" jsonschema:"the params for test suite"`
	Spec  *APISpec `json:"spec" jsonschema:"the API spec for test suite"`
}

type Pair struct {
	Key         string `json:"key" jsonschema:"the key of param"`
	Value       string `json:"value" jsonschema:"the value of param"`
	Description string `json:"description" jsonschema:"the description of param"`
}

type APISpec struct {
	Kind string `json:"kind" jsonschema:"the kind of API spec, such as swagger"`
	Url  string `json:"url" jsonschema:"the URL of API spec, such as http://localhost:8080/swagger.json"`
}

func (r *gRPCRunner) UpdateTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteArgs) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuite{
			Name: args.Name,
			Api:  args.API,
			Spec: &server.APISpec{
				Kind: args.Spec.Kind,
				Url:  args.Spec.Url,
			},
		}

		var reply *server.HelloReply
		reply, err = runner.UpdateTestSuite(ctx, suite)
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

func (r *gRPCRunner) DeleteTestSuite(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
		}

		var reply *server.HelloReply
		reply, err = runner.DeleteTestSuite(ctx, suite)
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

func (r *gRPCRunner) ListTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
		}

		var reply *server.Suite
		reply, err = runner.ListTestCase(ctx, suite)
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

type TestCaseIndentityRequest struct {
	Suite      string  `json:"suite" jsonschema:"the name of test suite"`
	Testcase   string  `json:"testcase" jsonschema:"the name of test case"`
	Parameters []*Pair `json:"parameters" jsonschema:"the params for test case"`
}

func (r *gRPCRunner) RunTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		testCase := &server.TestCaseIdentity{
			Suite:    args.Suite,
			Testcase: args.Testcase,
		}

		var reply *server.TestCaseResult
		reply, err = runner.RunTestCase(ctx, testCase)
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

func (r *gRPCRunner) GetTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		testCase := &server.TestCaseIdentity{
			Suite:    args.Suite,
			Testcase: args.Testcase,
		}

		var reply *server.TestCase
		reply, err = runner.GetTestCase(ctx, testCase)
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
					Header: convertMapToPairs(args.Headers),
					Query:  convertMapToPairs(args.QueryParams),
					Cookie: convertMapToPairs(args.Cookies),
					Form:   convertMapToPairs(args.FormParams),
				},
				Response: &server.Response{
					Body:       args.ExpectBody,
					StatusCode: args.ExpectStatus,
					Header:     convertMapToPairs(args.ExpectHeaders),
					Schema:     args.ExpectSchema,
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

func convertMapToPairs(data map[string]string) []*server.Pair {
	pairs := make([]*server.Pair, 0, len(data))
	for k, v := range data {
		pairs = append(pairs, &server.Pair{
			Key:   k,
			Value: v,
		})
	}
	return pairs
}

func (r *gRPCRunner) UpdateTestCase(ctx context.Context, request *mcp.CallToolRequest, args CreateTestCaseRequest) (
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
		reply, err = runner.UpdateTestCase(ctx, testCase)
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

func (r *gRPCRunner) GetSuggestedAPIs(ctx context.Context, request *mcp.CallToolRequest, args TestSuiteIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		suite := &server.TestSuiteIdentity{
			Name: args.Name,
			Api:  args.API,
		}

		var reply *server.TestCases
		reply, err = runner.GetSuggestedAPIs(ctx, suite)
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

func (r *gRPCRunner) DeleteTestCase(ctx context.Context, request *mcp.CallToolRequest, args TestCaseIndentityRequest) (
	result *mcp.CallToolResult, a any, err error) {
	var conn *grpc.ClientConn
	if conn, err = grpc.Dial(r.Address, grpc.WithInsecure()); err == nil {
		runner := server.NewRunnerClient(conn)

		testCase := &server.TestCaseIdentity{
			Suite:    args.Suite,
			Testcase: args.Testcase,
		}

		var reply *server.HelloReply
		reply, err = runner.DeleteTestCase(ctx, testCase)
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
