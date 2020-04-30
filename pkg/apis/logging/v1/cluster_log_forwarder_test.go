package v1_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	. "github.com/openshift/cluster-logging-operator/test"
)

var _ = Describe("ClusterLogForwarder", func() {
	It("serializes with conditions correctly", func() {
		forwarder := ClusterLogForwarder{
			Spec: ClusterLogForwarderSpec{
				Pipelines: []PipelineSpec{
					{
						InputRefs:  []string{InputNameApplication},
						OutputRefs: []string{"X", "Y"},
					},
					{
						InputRefs:  []string{InputNameInfrastructure, InputNameAudit},
						OutputRefs: []string{"Y", "Z"},
					},
					{
						InputRefs:  []string{InputNameAudit},
						OutputRefs: []string{"X", "Z"},
					},
				},
			},
			Status: ClusterLogForwarderStatus{
				Conditions: NewConditions(Condition{
					Type:               "Bad",
					Reason:             "NotGood",
					Status:             "True",
					LastTransitionTime: v1.Date(1999, time.January, 1, 0, 0, 0, 0, time.UTC),
				}),
			},
		}
		Expect(YAMLString(forwarder)).To(EqualLines(`
  metadata:
    creationTimestamp: null
  spec:
    pipelines:
    - inputRefs:
      - application
      outputRefs:
      - X
      - "Y"
    - inputRefs:
      - infrastructure
      - audit
      outputRefs:
      - "Y"
      - Z
    - inputRefs:
      - audit
      outputRefs:
      - X
      - Z
  status:
    conditions:
    - lastTransitionTime: "1999-01-01T00:00:00Z"
      reason: NotGood
      status: "True"
      type: Bad
`))
		Expect(JSONString(forwarder)).To(EqualLines(`{
    "metadata": {
      "creationTimestamp": null
    },
    "spec": {
      "pipelines": [
        {
          "outputRefs": [
            "X",
            "Y"
          ],
          "inputRefs": [
            "application"
          ]
        },
        {
          "outputRefs": [
            "Y",
            "Z"
          ],
          "inputRefs": [
            "infrastructure",
            "audit"
          ]
        },
        {
          "outputRefs": [
            "X",
            "Z"
          ],
          "inputRefs": [
            "audit"
          ]
        }
      ]
    },
    "status": {
      "conditions": [
        {
          "type": "Bad",
          "status": "True",
          "reason": "NotGood",
          "lastTransitionTime": "1999-01-01T00:00:00Z"
        }
      ]
    }
  }`))
	})
})
