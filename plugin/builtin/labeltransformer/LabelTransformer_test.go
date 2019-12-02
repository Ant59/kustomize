// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestLabelTransformer(t *testing.T) {
	tc := kusttest_test.NewPluginTestEnv(t).Set()
	defer tc.Reset()

	tc.BuildGoPlugin(
		"builtin", "", "LabelTransformer")

	th := kusttest_test.MakeHarnessEnhanced(t, "/app")

	rm := th.LoadAndRunTransformer(`
apiVersion: builtin
kind: LabelTransformer
metadata:
  name: notImportantHere
labels:
  app: myApp
  env: production
fieldSpecs:
  - path: metadata/labels
    create: true
`, `
apiVersion: v1
kind: Service
metadata:
  name: myService
spec:
  ports:
  - port: 7002
`)

	th.AssertActualEqualsExpected(rm, `
apiVersion: v1
kind: Service
metadata:
  labels:
    app: myApp
    env: production
  name: myService
spec:
  ports:
  - port: 7002
`)
}
