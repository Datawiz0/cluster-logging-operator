package k8shandler

import (
	"context"
	"fmt"
	"reflect"
	"time"

	monitoringv1 "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/pkg/constants"
	"github.com/openshift/cluster-logging-operator/pkg/k8shandler/collector"
	"github.com/openshift/cluster-logging-operator/pkg/logger"
	"github.com/openshift/cluster-logging-operator/pkg/utils"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/retry"

	configv1 "github.com/openshift/api/config/v1"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const (
	fluentdAlertsFile        = "fluentd/fluentd_prometheus_alerts.yaml"
	fluentdName              = "fluentd"
	syslogName               = "syslog"
	fluentdRequiredESVersion = "6"
)

func (clusterRequest *ClusterLoggingRequest) removeFluentd() (err error) {
	if clusterRequest.isManaged() {

		if err = clusterRequest.RemoveService(fluentdName); err != nil {
			return
		}

		if err = clusterRequest.RemoveServiceMonitor(fluentdName); err != nil {
			return
		}

		if err = clusterRequest.RemovePrometheusRule(fluentdName); err != nil {
			return
		}

		if err = clusterRequest.RemoveConfigMap(fluentdName); err != nil {
			return
		}

		if err = clusterRequest.RemoveConfigMap(constants.FluentdTrustedCAName); err != nil {
			return
		}

		if err = clusterRequest.RemoveDaemonset(fluentdName); err != nil {
			return
		}

		// Wait longer than the terminationGracePeriodSeconds
		time.Sleep(12 * time.Second)

		if err = clusterRequest.RemoveSecret(fluentdName); err != nil {
			return
		}

	}

	return nil
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdService() error {
	service := NewService(
		fluentdName,
		clusterRequest.cluster.Namespace,
		fluentdName,
		[]v1.ServicePort{
			{
				Port:       metricsPort,
				TargetPort: intstr.FromString(metricsPortName),
				Name:       metricsPortName,
			},
		},
	)

	service.Annotations = map[string]string{
		"service.alpha.openshift.io/serving-cert-secret-name": "fluentd-metrics",
	}

	utils.AddOwnerRefToObject(service, utils.AsOwner(clusterRequest.cluster))

	err := clusterRequest.Create(service)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure creating the fluentd service: %v", err)
	}

	return nil
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdServiceMonitor() error {

	cluster := clusterRequest.cluster

	serviceMonitor := NewServiceMonitor(fluentdName, cluster.Namespace)

	endpoint := monitoringv1.Endpoint{
		Port:   metricsPortName,
		Path:   "/metrics",
		Scheme: "https",
		TLSConfig: &monitoringv1.TLSConfig{
			CAFile:     prometheusCAFile,
			ServerName: fmt.Sprintf("%s.%s.svc", fluentdName, cluster.Namespace),
			// ServerName can be e.g. fluentd.openshift-logging.svc
		},
	}

	labelSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"logging-infra": "support",
		},
	}

	serviceMonitor.Spec = monitoringv1.ServiceMonitorSpec{
		JobLabel:  "monitor-fluentd",
		Endpoints: []monitoringv1.Endpoint{endpoint},
		Selector:  labelSelector,
		NamespaceSelector: monitoringv1.NamespaceSelector{
			MatchNames: []string{cluster.Namespace},
		},
	}

	utils.AddOwnerRefToObject(serviceMonitor, utils.AsOwner(cluster))

	err := clusterRequest.Create(serviceMonitor)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure creating the fluentd ServiceMonitor: %v", err)
	}

	return nil
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdPrometheusRule() error {
	ctx := context.TODO()
	cluster := clusterRequest.cluster

	rule := NewPrometheusRule(fluentdName, cluster.Namespace)

	spec, err := NewPrometheusRuleSpecFrom(utils.GetShareDir() + "/" + fluentdAlertsFile)
	if err != nil {
		return fmt.Errorf("failure creating the fluentd PrometheusRule: %w", err)
	}

	rule.Spec = *spec

	utils.AddOwnerRefToObject(rule, utils.AsOwner(cluster))

	err = clusterRequest.Create(rule)
	if err == nil {
		return nil
	}
	if !errors.IsAlreadyExists(err) {
		return fmt.Errorf("failure creating the fluentd PrometheusRule: %w", err)
	}

	current := &monitoringv1.PrometheusRule{}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err = clusterRequest.client.Get(ctx, types.NamespacedName{Name: rule.Name, Namespace: rule.Namespace}, current)
		if err != nil {
			logrus.Debugf("could not get prometheus rule %q: %v", rule.Name, err)
			return err
		}
		current.Spec = rule.Spec
		if err = clusterRequest.client.Update(ctx, current); err != nil {
			return err
		}
		logrus.Debug("updated prometheus rules")
		return nil
	})
}

// includeLegacyForwardConfig to address Bug 1782566.
// To be removed when legacy forwarding is unsupported
func (clusterRequest *ClusterLoggingRequest) includeLegacyForwardConfig() bool {
	config := &v1.ConfigMap{
		Data: map[string]string{},
	}
	if err := clusterRequest.Get("secure-forward", config); err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		logger.Warnf("There was a non-critical error trying to fetch the secure-forward configmap: %v", err)
	}
	_, found := config.Data["secure-forward.conf"]
	return found
}

// includeLegacySyslogConfig to address Bug 1799024.
// To be removed when legacy syslog is no longer supported.
func (clusterRequest *ClusterLoggingRequest) includeLegacySyslogConfig() bool {
	config := &v1.ConfigMap{
		Data: map[string]string{},
	}
	if err := clusterRequest.Get(syslogName, config); err != nil {
		if errors.IsNotFound(err) {
			return false
		}
		logger.Warnf("There was a non-critical error trying to fetch the configmap: %v", err)
	}
	_, found := config.Data["syslog.conf"]
	return found
}

// useOldRemoteSyslogPlugin checks if old plugin (docebo/fluent-plugin-remote-syslog) is to be used for sending syslog or new plugin (dlackty/fluent-plugin-remote_syslog) is to be used
func (clusterRequest *ClusterLoggingRequest) useOldRemoteSyslogPlugin() bool {
	if clusterRequest.ForwarderRequest == nil {
		return false
	}
	enabled, found := clusterRequest.ForwarderRequest.Annotations[UseOldRemoteSyslogPlugin]
	return found && enabled == "enabled"
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdConfigMap(fluentConf string) error {
	logrus.Debug("createOrUpdateFluentdConfigMap...")
	fluentdConfigMap := NewConfigMap(
		fluentdName,
		clusterRequest.cluster.Namespace,
		map[string]string{
			"fluent.conf": fluentConf,
			"run.sh":      string(utils.GetFileContents(utils.GetShareDir() + "/fluentd/run.sh")),
		},
	)

	utils.AddOwnerRefToObject(fluentdConfigMap, utils.AsOwner(clusterRequest.cluster))

	err := clusterRequest.Create(fluentdConfigMap)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure constructing Fluentd configmap: %v", err)
	}
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		current := &v1.ConfigMap{}
		if err = clusterRequest.Get(fluentdConfigMap.Name, current); err != nil {
			if errors.IsNotFound(err) {
				logrus.Debugf("Returning nil. The configmap %q was not found even though create previously failed.  Was it culled?", fluentdConfigMap.Name)
				return nil
			}
			return fmt.Errorf("Failed to get %v configmap for %q: %v", fluentdConfigMap.Name, clusterRequest.cluster.Name, err)
		}
		if reflect.DeepEqual(fluentdConfigMap.Data, current.Data) {
			return nil
		}
		current.Data = fluentdConfigMap.Data
		return clusterRequest.Update(current)
	})

	return retryErr
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdSecret() error {

	fluentdSecret := NewSecret(
		fluentdName,
		clusterRequest.cluster.Namespace,
		map[string][]byte{
			"ca-bundle.crt": utils.GetWorkingDirFileContents("ca.crt"),
			"tls.key":       utils.GetWorkingDirFileContents("system.logging.fluentd.key"),
			"tls.crt":       utils.GetWorkingDirFileContents("system.logging.fluentd.crt"),
		})

	utils.AddOwnerRefToObject(fluentdSecret, utils.AsOwner(clusterRequest.cluster))

	err := clusterRequest.CreateOrUpdateSecret(fluentdSecret)
	if err != nil {
		return err
	}

	return nil
}

func newFluentdPodSpec(cluster *logging.ClusterLogging, elasticsearchAppName string, elasticsearchInfraName string, proxyConfig *configv1.Proxy, trustedCABundleCM *v1.ConfigMap, pipelineSpec logging.ClusterLogForwarderSpec) v1.PodSpec {
	collectionSpec := logging.CollectionSpec{}
	if cluster.Spec.Collection != nil {
		collectionSpec = *cluster.Spec.Collection
	}
	resources := collectionSpec.Logs.FluentdSpec.Resources
	if resources == nil {
		resources = &v1.ResourceRequirements{
			Limits: v1.ResourceList{v1.ResourceMemory: defaultFluentdMemory},
			Requests: v1.ResourceList{
				v1.ResourceMemory: defaultFluentdMemory,
				v1.ResourceCPU:    defaultFluentdCpuRequest,
			},
		}
	}
	fluentdContainer := NewContainer("fluentd", "fluentd", v1.PullIfNotPresent, *resources)

	fluentdContainer.Ports = []v1.ContainerPort{
		{
			Name:          metricsPortName,
			ContainerPort: metricsPort,
			Protocol:      v1.ProtocolTCP,
		},
	}

	fluentdContainer.Env = []v1.EnvVar{
		{Name: "NODE_NAME", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "spec.nodeName"}}},
		{Name: "METRICS_CERT", Value: "/etc/fluent/metrics/tls.crt"},
		{Name: "METRICS_KEY", Value: "/etc/fluent/metrics/tls.key"},
		{Name: "NODE_IPV4", ValueFrom: &v1.EnvVarSource{FieldRef: &v1.ObjectFieldSelector{APIVersion: "v1", FieldPath: "status.hostIP"}}},
	}

	proxyEnv := utils.SetProxyEnvVars(proxyConfig)
	fluentdContainer.Env = append(fluentdContainer.Env, proxyEnv...)

	fluentdContainer.VolumeMounts = []v1.VolumeMount{
		{Name: "runlogjournal", MountPath: "/run/log/journal"},
		{Name: "varlog", MountPath: "/var/log"},
		{Name: "varlibdockercontainers", ReadOnly: true, MountPath: "/var/lib/docker"},
		{Name: "config", ReadOnly: true, MountPath: "/etc/fluent/configs.d/user"},
		{Name: "secureforwardconfig", ReadOnly: true, MountPath: "/etc/fluent/configs.d/secure-forward"},
		{Name: "secureforwardcerts", ReadOnly: true, MountPath: "/etc/ocp-forward"},
		{Name: "syslogconfig", ReadOnly: true, MountPath: "/etc/fluent/configs.d/syslog"},
		{Name: "syslogcerts", ReadOnly: true, MountPath: "/etc/ocp-syslog"},
		{Name: "entrypoint", ReadOnly: true, MountPath: "/opt/app-root/src/run.sh", SubPath: "run.sh"},
		{Name: "certs", ReadOnly: true, MountPath: "/etc/fluent/keys"},
		{Name: "localtime", ReadOnly: true, MountPath: "/etc/localtime"},
		{Name: "dockercfg", ReadOnly: true, MountPath: "/etc/sysconfig/docker"},
		{Name: "dockerdaemoncfg", ReadOnly: true, MountPath: "/etc/docker"},
		{Name: "filebufferstorage", MountPath: "/var/lib/fluentd"},
		{Name: metricsVolumeName, MountPath: "/etc/fluent/metrics"},
	}
	for _, target := range pipelineSpec.Outputs {
		if target.Secret != nil && target.Secret.Name != "" {
			path := fmt.Sprintf("/var/run/ocp-collector/secrets/%s", target.Secret.Name)
			fluentdContainer.VolumeMounts = append(fluentdContainer.VolumeMounts, v1.VolumeMount{Name: target.Name, MountPath: path})
		}
	}

	addTrustedCAVolume := false
	// If trusted CA bundle ConfigMap exists and its hash value is non-zero, mount the bundle.
	if trustedCABundleCM != nil && hasTrustedCABundle(trustedCABundleCM) {
		addTrustedCAVolume = true
		fluentdContainer.VolumeMounts = append(fluentdContainer.VolumeMounts,
			v1.VolumeMount{
				Name:      constants.FluentdTrustedCAName,
				ReadOnly:  true,
				MountPath: constants.TrustedCABundleMountDir,
			})
	}

	fluentdContainer.SecurityContext = &v1.SecurityContext{
		Privileged: utils.GetBool(true),
	}

	tolerations := utils.AppendTolerations(
		collectionSpec.Logs.FluentdSpec.Tolerations,
		[]v1.Toleration{
			{
				Key:      "node-role.kubernetes.io/master",
				Operator: v1.TolerationOpExists,
				Effect:   v1.TaintEffectNoSchedule,
			},
			{
				Key:      "node.kubernetes.io/disk-pressure",
				Operator: v1.TolerationOpExists,
				Effect:   v1.TaintEffectNoSchedule,
			},
		},
	)

	collectorContainer := collector.NewContainer()

	fluentdPodSpec := NewPodSpec(
		"logcollector",
		[]v1.Container{fluentdContainer, collectorContainer},
		[]v1.Volume{
			{Name: "runlogjournal", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/run/log/journal"}}},
			{Name: "varlog", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/var/log"}}},
			{Name: "varlibdockercontainers", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/var/lib/docker"}}},
			{Name: "config", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: "fluentd"}}}},
			{Name: collector.CollectorConfName, VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: collector.CollectorName}}}},
			{Name: "secureforwardconfig", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: "secure-forward"}, Optional: utils.GetBool(true)}}},
			{Name: "secureforwardcerts", VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: "secure-forward", Optional: utils.GetBool(true)}}},
			{Name: "syslogconfig", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: syslogName}, Optional: utils.GetBool(true)}}},
			{Name: "syslogcerts", VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: syslogName, Optional: utils.GetBool(true)}}},
			{Name: "entrypoint", VolumeSource: v1.VolumeSource{ConfigMap: &v1.ConfigMapVolumeSource{LocalObjectReference: v1.LocalObjectReference{Name: "fluentd"}}}},
			{Name: "certs", VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: "fluentd"}}},
			{Name: "localtime", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/etc/localtime"}}},
			{Name: "dockercfg", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/etc/sysconfig/docker"}}},
			{Name: "dockerdaemoncfg", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/etc/docker"}}},
			{Name: "filebufferstorage", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/var/lib/fluentd"}}},
			{Name: metricsVolumeName, VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: "fluentd-metrics"}}},
		},
		collectionSpec.Logs.FluentdSpec.NodeSelector,
		tolerations,
	)
	for _, target := range pipelineSpec.Outputs {
		if target.Secret != nil && target.Secret.Name != "" {
			fluentdPodSpec.Volumes = append(fluentdPodSpec.Volumes, v1.Volume{Name: target.Name, VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: target.Secret.Name}}})
		}
	}

	if addTrustedCAVolume {
		fluentdPodSpec.Volumes = append(fluentdPodSpec.Volumes,
			v1.Volume{
				Name: constants.FluentdTrustedCAName,
				VolumeSource: v1.VolumeSource{
					ConfigMap: &v1.ConfigMapVolumeSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: constants.FluentdTrustedCAName,
						},
						Items: []v1.KeyToPath{
							{
								Key:  constants.TrustedCABundleKey,
								Path: constants.TrustedCABundleMountFile,
							},
						},
					},
				},
			})
	}

	fluentdPodSpec.PriorityClassName = clusterLoggingPriorityClassName
	// Shorten the termination grace period from the default 30 sec to 10 sec.
	fluentdPodSpec.TerminationGracePeriodSeconds = utils.GetInt64(10)

	// FIXME: The following conditional branch is a refactoring candidate
	// To address here https://issues.redhat.com/browse/LOG-833
	//
	// Why do we need this branch at all?
	//
	// The ClusterLogging CR supports three additional use cases without providing
	// a logStore stanza namely:
	// Case 1: A collection only stanza to use with secure forward
	// Case 2: A collection only stanza to use with syslog
	// Case 3: A collection only stanza to use with ClusterLogFowarder API
	//
	// Supporting all three cases implies:
	// 1. We don't want to add the init container if `cluster.Spec.LogStore == nil`.
	// 2. We don't want to add the init container if `cluster.Spec.LogStore == nil`
	//    and a ClusterLogFowarder CR and does not provide a default output (only
	//    relevant for case 3)
	//
	// What is the init container bound to a log store stanza?
	//
	// The init container implementation is only relevant for the default ES log store.
	// It checks for the version of an ES cluster being >=6. In upgrade scenarios
	// an ES cluster can eventually provide nodes with multiple version (5.x and 6.x)
	// until all nodes are upgraded. The data model used in 6.x is not backward
	// compatible to 5.x, thus preventing fluentd from starting avoids data model
	// inconsistencies during ES upgrades from 5.x to 6.x.
	//
	// Why does this work w/o a ClusterLogForwarder CR provided?
	//
	// ClusterLogging and ClusterLogFowarder CR reconciliation share the same logic
	// for collection config generation. If first is provided without the latter a
	// ClusterLogFowarder with default fields is assumed.
	// (See ClusterLoggingRequest#getLogForwarder)
	if pipelineSpec.HasDefaultOutput() && cluster.Spec.LogStore != nil {
		fluentdPodSpec.InitContainers = []v1.Container{
			newFluentdInitContainer(cluster),
		}
	} else {
		fluentdPodSpec.InitContainers = []v1.Container{}
	}

	return fluentdPodSpec
}

func newFluentdInitContainer(cluster *logging.ClusterLogging) v1.Container {
	collectionSpec := logging.CollectionSpec{}
	resources := collectionSpec.Logs.FluentdSpec.Resources
	if resources == nil {
		resources = &v1.ResourceRequirements{
			Limits: v1.ResourceList{v1.ResourceMemory: defaultFluentdMemory},
			Requests: v1.ResourceList{
				v1.ResourceMemory: defaultFluentdMemory,
				v1.ResourceCPU:    defaultFluentdCpuRequest,
			},
		}
	}
	initContainer := NewContainer("fluentd-init", "fluentd", v1.PullIfNotPresent, *resources)

	initContainer.VolumeMounts = []v1.VolumeMount{
		{Name: "certs", ReadOnly: true, MountPath: "/etc/fluent/keys"},
	}

	initContainer.Command = []string{
		"./wait_for_es_version.sh",
		fluentdRequiredESVersion,
		fmt.Sprintf("https://%s.%s.svc:9200", elasticsearchResourceName, cluster.Namespace),
	}

	return initContainer
}

func (clusterRequest *ClusterLoggingRequest) createOrUpdateFluentdDaemonset(pipelineConfHash, collectorConfHash string, proxyConfig *configv1.Proxy) (err error) {

	cluster := clusterRequest.cluster

	fluentdTrustBundle := &v1.ConfigMap{}
	// Create or update cluster proxy trusted CA bundle.
	if proxyConfig != nil {
		fluentdTrustBundle, err = clusterRequest.createOrGetTrustedCABundleConfigMap(constants.FluentdTrustedCAName)
		if err != nil {
			return
		}
	}

	fluentdPodSpec := newFluentdPodSpec(cluster, "elasticsearch", "elasticsearch", proxyConfig, fluentdTrustBundle, clusterRequest.ForwarderSpec)

	fluentdDaemonset := NewDaemonSet("fluentd", cluster.Namespace, "fluentd", "fluentd", fluentdPodSpec)
	fluentdDaemonset.Spec.Template.Spec.Containers[0].Env = updateEnvVar(v1.EnvVar{Name: "FLUENT_CONF_HASH", Value: pipelineConfHash}, fluentdDaemonset.Spec.Template.Spec.Containers[0].Env)
	fluentdDaemonset.Spec.Template.Spec.Containers[1].Env = updateEnvVar(v1.EnvVar{Name: "CONF_HASH", Value: collectorConfHash}, fluentdDaemonset.Spec.Template.Spec.Containers[1].Env)

	annotations, err := clusterRequest.getFluentdAnnotations(fluentdDaemonset)
	if err != nil {
		return err
	}

	fluentdDaemonset.Spec.Template.Annotations = annotations

	uid := getServiceAccountLogCollectorUID()
	if len(uid) == 0 {
		// There's no uid for logcollector serviceaccount; setting ClusterLogging for the ownerReference.
		utils.AddOwnerRefToObject(fluentdDaemonset, utils.AsOwner(cluster))
	} else {
		// There's a uid for logcollector serviceaccount; setting the ServiceAccount for the ownerReference with blockOwnerDeletion.
		utils.AddOwnerRefToObject(fluentdDaemonset, NewLogCollectorServiceAccountRef(uid))
	}

	err = clusterRequest.Create(fluentdDaemonset)
	if err != nil && !errors.IsAlreadyExists(err) {
		return fmt.Errorf("Failure creating Fluentd Daemonset %v", err)
	}

	if clusterRequest.isManaged() {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			return clusterRequest.updateFluentdDaemonsetIfRequired(fluentdDaemonset)
		})
		if retryErr != nil {
			return retryErr
		}
	}

	return nil
}

func (clusterRequest *ClusterLoggingRequest) updateFluentdDaemonsetIfRequired(desired *apps.DaemonSet) (err error) {
	logger.DebugObject("desired fluent update: %v", desired)
	current := &apps.DaemonSet{}

	if err = clusterRequest.Get(desired.Name, current); err != nil {
		if errors.IsNotFound(err) {
			// the object doesn't exist -- it was likely culled
			// recreate it on the next time through if necessary
			return nil
		}
		return fmt.Errorf("Failed to get Fluentd daemonset: %v", err)
	}

	flushBuffer := isBufferFlushRequired(current, desired)
	desired, different := isFluentdDaemonsetDifferent(current, desired)

	if different {
		current.Spec = desired.Spec
		if flushBuffer {
			current.Spec.Template.Spec.Containers[0].Env = updateEnvVar(v1.EnvVar{Name: "FLUSH_AT_SHUTDOWN", Value: "True"}, current.Spec.Template.Spec.Containers[0].Env)
			if err = clusterRequest.Update(current); err != nil {
				logrus.Debugf("Failed to prepare Fluentd daemonset to flush its buffers: %v", err)
				return err
			}

			// wait for pods to all restart then continue
			if err = clusterRequest.waitForDaemonSetReady(current); err != nil {
				return fmt.Errorf("Timed out waiting for Fluentd to be ready")
			}
		}
		logger.DebugObject("updating fluentd to: %v", current)
		if err = clusterRequest.Update(desired); err != nil {
			return err
		}
	}

	return nil
}

func (clusterRequest *ClusterLoggingRequest) getFluentdAnnotations(daemonset *apps.DaemonSet) (map[string]string, error) {

	if daemonset.Spec.Template.ObjectMeta.Annotations == nil {
		daemonset.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}

	annotations := daemonset.Spec.Template.ObjectMeta.Annotations

	fluentdTrustBundle := &v1.ConfigMap{}
	fluentdTrustBundleName := types.NamespacedName{Name: constants.FluentdTrustedCAName, Namespace: constants.OpenshiftNS}
	if err := clusterRequest.client.Get(context.TODO(), fluentdTrustBundleName, fluentdTrustBundle); err != nil {
		if !errors.IsNotFound(err) {
			return annotations, err
		}
	}

	if _, ok := fluentdTrustBundle.Data[constants.TrustedCABundleKey]; !ok {
		return annotations, fmt.Errorf("%v does not yet contain expected key %v", fluentdTrustBundle.Name, constants.TrustedCABundleKey)
	}

	trustedCAHashValue, err := calcTrustedCAHashValue(fluentdTrustBundle)
	if err != nil {
		return annotations, fmt.Errorf("unable to calculate trusted CA value. E: %s", err.Error())
	}

	if trustedCAHashValue == "" {
		return annotations, fmt.Errorf("Did not receive hashvalue for trusted CA value")
	}

	annotations[constants.TrustedCABundleHashName] = trustedCAHashValue

	return annotations, nil
}

func (clusterRequest *ClusterLoggingRequest) RestartFluentd(proxyConfig *configv1.Proxy) (err error) {

	normalizerConfig, err := clusterRequest.generateCollectorConfig()
	if err != nil {
		return err
	}

	logger.Debugf("Generated normalizer config: %s", normalizerConfig)
	normalizerConfHash, err := utils.CalculateMD5Hash(normalizerConfig)
	if err != nil {
		logger.Errorf("unable to calculate MD5 hash. E: %s", err.Error())
		return err
	}
	collectorConfig, err := collector.GenerateConfig()
	if err != nil {
		return err
	}
	logger.Debugf("Generated collector config: %s", collectorConfig)
	collectorConfHash, err := utils.CalculateMD5Hash(collectorConfig)
	if err != nil {
		logger.Errorf("unable to calculate MD5 hash. E: %s", err.Error())
		return err
	}
	if err = clusterRequest.createOrUpdateFluentdDaemonset(normalizerConfHash, collectorConfHash, proxyConfig); err != nil {
		return err
	}

	return clusterRequest.UpdateFluentdStatus()
}

//updateEnvar adds the value to the list or replaces it if it already existing
func updateEnvVar(value v1.EnvVar, values []v1.EnvVar) []v1.EnvVar {
	found := false
	for i, envvar := range values {
		if envvar.Name == value.Name {
			values[i] = value
			found = true
		}
	}
	if !found {
		values = append(values, value)
	}
	return values
}

func isFluentdDaemonsetDifferent(current, desired *apps.DaemonSet) (*apps.DaemonSet, bool) {
	current, different := isDaemonsetDifferent(current, desired)

	currentHash := current.Spec.Template.ObjectMeta.Annotations[constants.TrustedCABundleHashName]
	desiredHash := desired.Spec.Template.ObjectMeta.Annotations[constants.TrustedCABundleHashName]
	if currentHash != desiredHash {
		current.Spec.Template.ObjectMeta.Annotations[constants.TrustedCABundleHashName] = desiredHash
		different = true
	}

	return current, different
}
