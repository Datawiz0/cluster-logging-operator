package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//LogSourceType is an explicitly defined log source
type LogSourceType string

const (

	//LogSourceTypeApp are container logs from non-infra structure containers
	LogSourceTypeApp LogSourceType = "logs.app"

	//LogSourceTypeInfra are logs from infra structure containers or node logs
	LogSourceTypeInfra LogSourceType = "logs.infra"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// LogForwading is the Schema for the logforwardings API
// +k8s:openapi-gen=true
type LogForwarding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ForwardingSpec    `json:"spec,omitempty"`
	Status *ForwardingStatus `json:"status,omitempty"`
}

//ForwardingSpec specifies log forwarding pipelines from a defined sources to dest outputs
// +k8s:openapi-gen=true
type ForwardingSpec struct {
	DisableDefaultForwarding bool           `json:"disableDefaultForwarding,omitempty"`
	Outputs                  []OutputSpec   `json:"outputs,omitempty"`
	Pipelines                []PipelineSpec `json:"pipelines,omitempty"`
}

//PipelineSpec is the sources spec to named targets
type PipelineSpec struct {
	Name       string        `json:"name,omitempty"`
	SourceType LogSourceType `json:"inputType,omitempty"`

	//OutputRefs is a list of  the names of outputs defined by forwarding.outputs
	OutputRefs []string `json:"outputRefs,omitempty"`
}

//OutputSpec specifies destination config for log message endpoints
type OutputSpec struct {
	Type     OutputType        `json:"type,omitempty"`
	Name     string            `json:"name,omitempty"`
	Endpoint string            `json:"endpoint,omitempty"`
	Secret   *OutputSecretSpec `json:"secret,omitempty"`
}

//OutputSecretSpec specifies secrets for pipelines
type OutputSecretSpec struct {

	//Name is the name of the secret to use with this output
	Name string `json:"name"`
}

//OutputType defines the type of endpoint that will receive messages
type OutputType string

const (
	//OutputTypeElasticsearch configures pipeline to send messages to elasticsearch
	OutputTypeElasticsearch OutputType = "elasticsearch"

	//OutputTypeForward configures the pipeline to send messages via Fluent's secure forward
	OutputTypeForward OutputType = "forward"
)

type LogForwardingStatus struct {
	LogForwardingState LogForwardingState `json:"state,omitempty"`
	Reason             string             `json:"reason,omitempty"`
}

type LogForwardingState string

const (
	LogForwardingStateAccepted LogForwardingState = "accepted"
	LogForwardingStateRejected LogForwardingState = "rejected"
)

//ForwardingStatus is the status of spec'd forwarding
// +k8s:openapi-gen=true
type ForwardingStatus struct {
	LogForwardingStatus LogForwardingStatus `json:"status,omitempty"`

	//LogSources lists the configured log sources
	LogSources []string `json:"sources,omitempty"`

	//Pipelines is the status of the outputs
	Pipelines []PipelineStatus `json:"pipelines,omitempty"`

	//Outputs is the status of the outputs
	Outputs []OutputStatus `json:"outputs,omitempty"`
}

//PipelineStatus is the status of a give pipeline
type PipelineStatus struct {
	//Name of the corresponding pipeline for this status
	Name string `json:"name,omitempty"`
	//State of the corresponding pipeline for this status
	State PipelineState `json:"state,omitempty"`
	//Reasons for the state of the corresponding pipeline for this status
	Reasons []PipelineStateReason `json:"reason,omitempty"`
	//Message about the corresponding pipeline
	Message string `json:"Message,omitempty"`
}

//PipelineState is the state of a spec'd pipeline
type PipelineState string

//PipelineStateReason provides a reason for the given pipeline state
type PipelineStateReason string

const (
	//PipelineStateAccepted is accepted by forwarding and includes all required fields to send messages to all spec'd outputs
	PipelineStateAccepted PipelineState = "Accepted"
	//PipelineStateDegraded is parially accepted by forwarding and includes some of the required fields to send messages to some of the spec'd outputs
	PipelineStateDegraded PipelineState = "Degraded"
	//PipelineStateDropped dropped by forwarding because its missing required fields to send messages to outputs
	PipelineStateDropped PipelineState = "Dropped"

	PipelineStateReasonUnrecognizedOutput   PipelineStateReason = "An outputRef does not correspond to a defined output"
	PipelineStateReasonUnrecognizedSource   PipelineStateReason = "The source type is unrecognized"
	PipelineStateReasonNonUniqueName        PipelineStateReason = "The pipeline is missing a unique name"
	PipelineStateReasonMissingName          PipelineStateReason = "The pipeline is missing a name"
	PipelineStateReasonMissingOutputs       PipelineStateReason = "The pipeline does not reference any outputs or referes more than exist"
	PipelineStateReasonMissingSource        PipelineStateReason = "The pipeline is missing a source type"
	PipelineStateReasonReservedNameConflict PipelineStateReason = "Name conflicts with an internally reserved name"
)

//OutputStatus of a given output
type OutputStatus struct {
	//Name of the corresponding output for this status
	Name string `json:"name,omitempty"`
	//State of the corresponding output for this status
	State OutputState `json:"state,omitempty"`
	//Reasons for the state of the corresponding output for this status
	Reasons []OutputStateReason `json:"reasons,omitempty"`
	//Message about the corresponding output
	Message string `json:"message,omitempty"`
}

//OutputState is the state of a spec'd output
type OutputState string

//OutputStateReason provides a reason for the given state
type OutputStateReason string

const (
	//OutputStateAccepted means the output is usuable by forwarding and is spec'd with all the required fields
	OutputStateAccepted OutputState = "Accepted"

	//OutputStateDropped means the output is unusuable by forwarding is missing some the required fields
	OutputStateDropped OutputState = "Dropped"

	//OutputStateNonUniqueName is not unique amoung all defined outputs
	OutputStateNonUniqueName OutputStateReason = "NonUnique Name"

	//OutputStateReservedNameConflict is not unique amoung all defined outputs
	OutputStateReservedNameConflict OutputStateReason = "Name conflicts with an internally reserved name"

	//OutputStateReasonMissingName is missing a name
	OutputStateReasonMissingName OutputStateReason = "Missing name"

	//OutputStateReasonMissingType is missing a name
	OutputStateReasonMissingType OutputStateReason = "Missing type"

	//OutputStateReasonMissingEndpoint is missing the endpoint information, it is empty, or is an invalid format
	OutputStateReasonMissingEndpoint OutputStateReason = "Missing endpoint"

	//OutputStateReasonMissingSecretName is missing the name of the secret
	OutputStateReasonMissingSecretName OutputStateReason = "Missing secret name"

	//OutputStateReasonSecretDoesNotExist for secrets which don't exist
	OutputStateReasonSecretDoesNotExist OutputStateReason = "Secret does not exist"

	//OutputStateReasonUnrecognizedType has an unrecognized or supported output type
	OutputStateReasonUnrecognizedType OutputStateReason = "Unrecognized type"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LogForwardingList contains a list of LogForwarding
type LogForwardingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ForwardingSpec `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LogForwarding{}, &LogForwardingList{})
}