package server_mocks_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/KevinJCross/server_mocks"
	"github.com/aws/aws-lambda-go/events"
	. "github.com/onsi/ginkgo/v2"
	"github.com/steinfletcher/apitest"
)

var _ = Describe("Mock tests", func() {
	Context("Handler tests", func() {
		It("should respond on /", func() {
			api().
				Get("/").
				Expect(GinkgoT()).
				Body(`{"message":"hello"}`).
				End()
		})
	})
})

func api() *apitest.APITest {
	return apitest.New("mock-test").Handler(ToHttpHandler(server_mocks.New().Lambda()))
}

func ToHttpHandler(lambda func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		proxyResponse, err := lambda(r.Context(), events.APIGatewayProxyRequest{
			Resource:                        "/*",
			Path:                            r.URL.Path,
			HTTPMethod:                      r.Method,
			Headers:                         singleValue(r.Header),
			MultiValueHeaders:               r.Header,
			QueryStringParameters:           singleValue(r.URL.Query()),
			MultiValueQueryStringParameters: r.URL.Query(),
			PathParameters:                  parsePathParams("/*", r.URL.Path),
			StageVariables:                  map[string]string{},
			Body:                            string(body),
		})

		if err != nil {
			writeError(w, err)
			return
		}

		writeResponse(w, proxyResponse)
	})
}

func writeError(w http.ResponseWriter, err error) {
	// write a generic error, the same as API GW would if an error was returned by handler
	w.WriteHeader(http.StatusInternalServerError)
	escapedString, err := json.Marshal(err.Error())
	if err != nil {
		panic(err)
	}
	_, _ = w.Write([]byte(fmt.Sprintf(`{"error": ["%s"]\n`, escapedString)))
}

func singleValue(multiValueMap map[string][]string) map[string]string {
	singleValueMap := make(map[string]string)
	for k, mv := range multiValueMap {
		if len(mv) > 0 {
			singleValueMap[k] = mv[0]
		}
	}
	return singleValueMap
}

func parsePathParams(pathPattern string, path string) map[string]string {
	exp := regexp.MustCompile(`{(\w+)}`)
	pathPatternExp := regexp.MustCompile(exp.ReplaceAllString(pathPattern, `(?P<$1>\w+)`))

	subMatches := pathPatternExp.FindStringSubmatch(path)
	subMatchNames := pathPatternExp.SubexpNames()

	params := make(map[string]string)
	for i, paramName := range subMatchNames {
		if paramName == "" || len(subMatches) < i {
			continue
		}
		params[paramName] = subMatches[i]
	}

	return params
}

func writeResponse(w http.ResponseWriter, proxyResponse events.APIGatewayProxyResponse) {
	for k, v := range proxyResponse.Headers {
		w.Header().Add(k, v)
	}

	for k, vs := range proxyResponse.MultiValueHeaders {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(proxyResponse.StatusCode)
	_, _ = w.Write([]byte(proxyResponse.Body))
}
