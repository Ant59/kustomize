// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestNamespacedGenerator(t *testing.T) {
	th := kusttest_test.MakeHarness(t)
	th.WriteK("/app", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
configMapGenerator:
- name: the-non-default-namespace-map
  namespace: non-default
  literals:
  - altGreeting=Good Morning from non-default namespace!
  - enableRisky="false"
- name: the-map
  literals:
  - altGreeting=Good Morning from default namespace!
  - enableRisky="false"

secretGenerator:
- name: the-non-default-namespace-secret
  namespace: non-default
  literals:
  - password.txt=verySecret
- name: the-secret
  literals:
  - password.txt=anotherSecret
`)
	m := th.Run("/app", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  altGreeting: Good Morning from non-default namespace!
  enableRisky: "false"
kind: ConfigMap
metadata:
  name: the-non-default-namespace-map-b6h49k7mt8
  namespace: non-default
---
apiVersion: v1
data:
  altGreeting: Good Morning from default namespace!
  enableRisky: "false"
kind: ConfigMap
metadata:
  name: the-map-4959m5tm6c
---
apiVersion: v1
data:
  password.txt: dmVyeVNlY3JldA==
kind: Secret
metadata:
  name: the-non-default-namespace-secret-h8d9hkgtb9
  namespace: non-default
type: Opaque
---
apiVersion: v1
data:
  password.txt: YW5vdGhlclNlY3JldA==
kind: Secret
metadata:
  name: the-secret-fgb45h45bh
type: Opaque
`)
}

func TestNamespacedGeneratorWithOverlays(t *testing.T) {
	th := kusttest_test.MakeHarness(t)
	th.WriteK("/app/base", `
namespace: base

configMapGenerator:
- name: testCase
  literals:
    - base=true
`)
	th.WriteK("/app/overlay", `
resources:
  - ../base

namespace: overlay

configMapGenerator:
  - name: testCase
    behavior: merge
    literals:
      - overlay=true
`)
	m := th.Run("/app/overlay", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  base: "true"
  overlay: "true"
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: testCase-4g75kbk6gm
  namespace: overlay
`)
}
