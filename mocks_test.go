package server_mocks_test

import (
	"github.com/KevinJCross/server_mocks"
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
	return apitest.New("mock-test").Handler(server_mocks.New().Handler())
}
