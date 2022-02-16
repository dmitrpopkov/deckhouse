/*
Copyright 2021 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/testing/hooks"
)

func Test(t *testing.T) {
	hooks.SetGinkgoParallelNodes()

	RegisterFailHandler(Fail)
	RunSpecs(t, "")
}

type HTTPMockResponse struct {
	Response string
	Code     int
}

type jwtPayload struct {
	Iss   string
	Sub   string
	Aud   string
	Scope string
	Nbf   int64
	Exp   int64
}
