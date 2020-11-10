package functional

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/ViaQ/logerr/log"
	"github.com/openshift/cluster-logging-operator/internal/pkg/generator/forwarder"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/pkg/certificates"
	"github.com/openshift/cluster-logging-operator/pkg/components/fluentd"
	"github.com/openshift/cluster-logging-operator/pkg/constants"
	"github.com/openshift/cluster-logging-operator/pkg/utils"
	"github.com/openshift/cluster-logging-operator/test"
	"github.com/openshift/cluster-logging-operator/test/client"
	"github.com/openshift/cluster-logging-operator/test/helpers/cmd"
	"github.com/openshift/cluster-logging-operator/test/helpers/oc"
	"github.com/openshift/cluster-logging-operator/test/runtime"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	ApplicationLogFile = "/tmp/app-logs"
)

var (
	maxDuration          time.Duration
	defaultRetryInterval time.Duration
)

//FluentddFunctionalFramework deploys stand alone fluentd with the fluent.conf as generated by input ClusterLogForwarder CR
type FluentdFunctionalFramework struct {
	Name              string
	Namespace         string
	Conf              string
	image             string
	labels            map[string]string
	Forwarder         *logging.ClusterLogForwarder
	test              *client.Test
	pod               *corev1.Pod
	fluentContainerId string
	closeClient       func()
}

func init() {
	maxDuration, _ = time.ParseDuration("2m")
	defaultRetryInterval, _ = time.ParseDuration("1ms")
}

func NewFluentdFunctionalFramework() *FluentdFunctionalFramework {
	test := client.NewTest()
	return NewFluentdFunctionalFrameworkUsing(client.NewTest(), test.Close, 0)
}

func NewFluentdFunctionalFrameworkUsing(t *client.Test, fnClose func(), verbosity int) *FluentdFunctionalFramework {
	if level, found := os.LookupEnv("LOG_LEVEL"); found {
		if i, err := strconv.Atoi(level); err == nil {
			verbosity = int(i)
		}
	}

	log.MustInit("functional-framework")
	log.SetLogLevel(verbosity)
	testName := "functional"
	framework := &FluentdFunctionalFramework{
		Name:      testName,
		Namespace: t.NS.Name,
		image:     utils.GetComponentImage(constants.FluentdName),
		labels: map[string]string{
			"testtype": "functional",
			"testname": testName,
		},
		test:        t,
		Forwarder:   runtime.NewClusterLogForwarder(),
		closeClient: fnClose,
	}
	return framework
}

func (f *FluentdFunctionalFramework) Cleanup() {
	f.closeClient()
}

func (f *FluentdFunctionalFramework) RunCommand(container string, cmd ...string) (string, error) {
	log.V(2).Info("Running", "container", container, "cmd", cmd)
	out, err := runtime.ExecContainer(f.pod, container, cmd[0], cmd[1:]...).CombinedOutput()
	log.V(2).Info("Exec'd", "out", string(out), "err", err)
	return string(out), err
}

//Deploy the objects needed to functional test
func (f *FluentdFunctionalFramework) Deploy() (err error) {
	log.V(2).Info("Generating config", "forwarder", f.Forwarder)
	yaml, _ := yaml.Marshal(f.Forwarder)
	if f.Conf, err = forwarder.Generate(string(yaml), false); err != nil {
		return err
	}
	log.V(2).Info("Generating Certificates")
	if err, _ = certificates.GenerateCertificates(f.test.NS.Name,
		utils.GetScriptsDir(), "elasticsearch",
		utils.DefaultWorkingDir); err != nil {
		return err
	}
	log.V(2).Info("Creating config configmap")
	configmap := runtime.NewConfigMap(f.test.NS.Name, f.Name, map[string]string{})
	runtime.NewConfigMapBuilder(configmap).
		Add("fluent.conf", f.Conf).
		Add("run.sh", fluentd.RunScript)
	if err = f.test.Client.Create(configmap); err != nil {
		return err
	}

	log.V(2).Info("Creating certs configmap")
	certsName := "certs-" + f.Name
	certs := runtime.NewConfigMap(f.test.NS.Name, certsName, map[string]string{})
	runtime.NewConfigMapBuilder(certs).
		Add("tls.key", string(utils.GetWorkingDirFileContents("system.logging.fluentd.key"))).
		Add("tls.crt", string(utils.GetWorkingDirFileContents("system.logging.fluentd.crt")))
	if err = f.test.Client.Create(certs); err != nil {
		return err
	}

	log.V(2).Info("Creating service")
	service := runtime.NewService(f.test.NS.Name, f.Name)
	runtime.NewServiceBuilder(service).
		AddServicePort(24231, 24231).
		WithSelector(f.labels)
	if err = f.test.Client.Create(service); err != nil {
		return err
	}

	role := runtime.NewRole(f.test.NS.Name, f.Name,
		v1.PolicyRule{
			Verbs:     []string{"list", "get"},
			Resources: []string{"pods", "namespaces"},
			APIGroups: []string{""},
		},
	)
	if err = f.test.Client.Create(role); err != nil {
		return err
	}
	rolebinding := runtime.NewRoleBinding(f.test.NS.Name, f.Name,
		v1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     role.Name,
		},
		v1.Subject{
			Kind: "ServiceAccount",
			Name: "default",
		},
	)
	if err = f.test.Client.Create(rolebinding); err != nil {
		return err
	}

	log.V(2).Info("Defining pod...")
	containers := []corev1.Container{}
	f.pod = runtime.NewPod(f.test.NS.Name, f.Name, containers...)
	b := runtime.NewPodBuilder(f.pod).
		WithLabels(f.labels).
		AddConfigMapVolume("config", f.Name).
		AddConfigMapVolume("entrypoint", f.Name).
		AddConfigMapVolume("certs", certsName).
		AddContainer(constants.FluentdName, f.image).
		AddEnvVar("LOG_LEVEL", "debug").
		AddEnvVarFromFieldRef("POD_IP", "status.podIP").
		AddVolumeMount("config", "/etc/fluent/configs.d/user", "", true).
		AddVolumeMount("entrypoint", "/opt/app-root/src/run.sh", "run.sh", true).
		AddVolumeMount("certs", "/etc/fluent/metrics", "", true).
		End()
	if err = f.addOutputContainers(b, f.Forwarder.Spec.Outputs); err != nil {
		return err
	}
	log.V(2).Info("Creating pod", "pod", f.pod)
	if err = f.test.Client.Create(f.pod); err != nil {
		return err
	}

	log.V(2).Info("waiting for pod to be ready")
	if err = oc.Literal().From(fmt.Sprintf("oc wait -n %s pod/%s --timeout=60s --for=condition=Ready", f.test.NS.Name, f.Name)).Output(); err != nil {
		return err
	}
	if err = f.test.Client.Get(f.pod); err != nil {
		return err
	}
	log.V(2).Info("waiting for service endpoints to be ready")
	err = wait.PollImmediate(time.Second*2, time.Second*10, func() (bool, error) {
		ips, err := oc.Get().WithNamespace(f.test.NS.Name).Resource("endpoints", f.Name).OutputJsonpath("{.subsets[*].addresses[*].ip}").Run()
		if err != nil {
			return false, nil
		}
		// if there are IPs in the service endpoint, the service is available
		if strings.TrimSpace(ips) != "" {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return fmt.Errorf("service could not be started")
	}
	for _, cs := range f.pod.Status.ContainerStatuses {
		if cs.Name == constants.FluentdName {
			f.fluentContainerId = strings.TrimPrefix(cs.ContainerID, "cri-o://")
			break
		}
	}
	return nil
}

func (f *FluentdFunctionalFramework) addOutputContainers(b *runtime.PodBuilder, outputs []logging.OutputSpec) error {
	log.V(2).Info("Adding outputs", "outputs", outputs)
	for _, output := range outputs {
		switch output.Type {
		case logging.OutputTypeFluentdForward:
			if err := f.addForwardOutput(b, output); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *FluentdFunctionalFramework) WritesApplicationLogs(numOfLogs uint64) error {
	return f.WritesNApplicationLogsOfSize(numOfLogs, uint64(100))
}

func (f *FluentdFunctionalFramework) WritesNApplicationLogsOfSize(numOfLogs, size uint64) error {
	msg := "$(date -u +'%Y-%m-%dT%H:%M:%S.%N%:z') stdout F $msg "
	//podname_ns_containername-containerid.log
	//functional_testhack-16511862744968_fluentd-90a0f0a7578d254eec640f08dd155cc2184178e793d0289dff4e7772757bb4f8.log
	filepath := fmt.Sprintf("/var/log/containers/%s_%s_%s-%s.log", f.pod.Name, f.pod.Namespace, constants.FluentdName, f.fluentContainerId)
	result, err := f.RunCommand(constants.FluentdName, "bash", "-c", fmt.Sprintf("bash -c 'mkdir -p /var/log/containers;echo > %s;msg=$(cat /dev/urandom|tr -dc 'a-zA-Z0-9'|fold -w %d|head -n 1);for n in $(seq 1 %d);do echo %s >> %s; done'", filepath, size, numOfLogs, msg, filepath))
	log.V(3).Info("FluentdFunctionalFramework.WritesApplicationLogs", "result", result, "err", err)
	return err
}

func (f *FluentdFunctionalFramework) ReadApplicationLogsFrom(outputName string) (result string, err error) {
	err = wait.PollImmediate(defaultRetryInterval, maxDuration, func() (done bool, err error) {
		result, err = f.RunCommand(strings.ToLower(outputName), "cat", ApplicationLogFile)
		if err == nil {
			return true, nil
		}
		log.V(4).Error(err, "Polling application logs")
		return false, nil
	})
	if err == nil {
		result = fmt.Sprintf("[%s]", strings.Join(strings.Split(strings.TrimSpace(result), "\n"), ","))
	}
	log.V(4).Info("Returning", "logs", result)
	return result, err
}

func (f *FluentdFunctionalFramework) ReadNApplicationLogsFrom(n uint64, outputName string) ([]string, error) {
	lines := []string{}
	ctx, cancel := context.WithTimeout(context.Background(), test.SuccessTimeout())
	defer cancel()
	reader, err := cmd.TailReaderForContainer(f.pod, outputName, ApplicationLogFile)
	if err != nil {
		log.V(3).Error(err, "Error creating tail reader")
		return nil, err
	}
	for {
		line, err := reader.ReadLineContext(ctx)
		if err != nil {
			log.V(3).Error(err, "Error readlinecontext")
			return nil, err
		}
		lines = append(lines, line)
		n--
		if n == 0 {
			break
		}
	}
	return lines, err
}
