package forwarding

import (
	"fmt"

	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/pkg/generators/forwarding/fluentd"
)

//NewConfigGenerator create a config generator for a given collector type
func NewConfigGenerator(collector logging.LogCollectionType) (ConfigGenerator, error) {
	switch collector {
	case logging.LogCollectionTypeFluentd:
		return fluentd.NewConfigGenerator()
	}
	return nil, fmt.Errorf("Config generation not supported for collects of type %s", collector)
}

//ConfigGenerator is a config generator for a given ForwardingSpec
type ConfigGenerator interface {

	//Generate the config
	Generate(pipeline *logging.ForwardingSpec) (string, error)
}
