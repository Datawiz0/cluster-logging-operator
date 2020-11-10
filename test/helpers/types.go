package helpers

import (
	"encoding/json"
	"strings"

	logger "github.com/ViaQ/logerr/log"
	"github.com/openshift/cluster-logging-operator/test"
)

type Logs []Log

type docker struct {
	ContainerID string `json:"container_id"`
}

type k8s struct {
	ContainerName    string            `json:"container_name"`
	ContainerImage   string            `json:"container_image"`
	ContainerImageID string            `json:"container_image_id"`
	PodName          string            `json:"pod_name"`
	PodID            string            `json:"pod_id"`
	Host             string            `json:"host"`
	Labels           map[string]string `json:"labels"`
	FlatLabels       []string          `json:"flat_labels"`
	MasterURL        string            `json:"master_url"`
	NamespaceName    string            `json:"namespace_name"`
	NamespaceID      string            `json:"namespace_id"`
}

type pipelineMetadata struct {
	Collector *struct {
		IPaddr4    string `json:"ipaddr4"`
		IPaddr6    string `json:"ipaddr6"`
		InputName  string `json:"inputname"`
		Name       string `json:"name"`
		ReceivedAt string `json:"received_at"`
		Version    string `json:"version"`
	} `json:"collector"`
}

type Log struct {
	Docker           *docker           `json:"docker"`
	Kubernetes       *k8s              `json:"kubernetes"`
	Message          string            `json:"message"`
	Level            string            `json:"level"`
	Hostname         string            `json:"hostname"`
	PipelineMetadata *pipelineMetadata `json:"pipeline_metadata"`
	Timestamp        string            `json:"@timestamp"`
	IndexName        string            `json:"viaq_index_name"`
	MessageID        string            `json:"viaq_msg_id"`
	OpenshiftLabels  openshiftMeta     `json:"openshift"`
	Timing
}

type Timing struct {
	//EpocIn is only added during benchmark testing
	EpocIn float64 `json:"epoc_in"`
	//EpocOut is only added during benchmark testing
	EpocOut float64 `json:"epoc_out"`
}

func (t *Timing) Difference() float64 {
	return t.EpocOut - t.EpocIn
}

//Bloat is the ratio of overall size / Message size
func (l *Log) Bloat() float64 {
	return float64(len(l.String())) / float64(len(l.Message))
}

type openshiftMeta struct {
	Labels map[string]string `json:"labels"`
}

func ParseLogs(in string) (Logs, error) {
	logger.V(3).Info("ParseLogs", "content", in)
	logs := []Log{}
	if in == "" {
		return logs, nil
	}

	err := json.Unmarshal([]byte(in), &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (l *Log) String() string {
	return test.JSONLine(l)
}

func (l Logs) ByIndex(prefix string) Logs {
	filtered := []Log{}
	for _, entry := range l {
		if strings.HasPrefix(entry.IndexName, prefix) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (l Logs) ByPod(name string) Logs {
	filtered := []Log{}
	for _, entry := range l {
		if entry.Kubernetes == nil {
			continue
		}
		if entry.Kubernetes.PodName == name {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (l Logs) NonEmpty() bool {
	if l == nil {
		return false
	}
	return len(l) > 0
}
