package server_mocks_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServerMocks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ServerMocks Suite")
}
