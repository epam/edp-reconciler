package cluster

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"strconv"
)

var log = ctrl.Log.WithName("cluster-util")

const (
	watchNamespaceEnvVar   = "WATCH_NAMESPACE"
	debugModeEnvVar        = "DEBUG_MODE"
	inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func GetOwnerReference(ownerKind string, ors []metav1.OwnerReference) *metav1.OwnerReference {
	log.Info("finding owner", "kind", ownerKind)
	if len(ors) == 0 {
		return nil
	}
	for _, o := range ors {
		if o.Kind == ownerKind {
			return &o
		}
	}
	return nil
}

// GetWatchNamespace returns the namespace the operator should be watching for changes
func GetWatchNamespace() (string, error) {
	ns, found := os.LookupEnv(watchNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", watchNamespaceEnvVar)
	}
	return ns, nil
}

// GetDebugMode returns the debug mode value
func GetDebugMode() (bool, error) {
	mode, found := os.LookupEnv(debugModeEnvVar)
	if !found {
		return false, nil
	}

	b, err := strconv.ParseBool(mode)
	if err != nil {
		return false, err
	}
	return b, nil
}

// Check whether the operator is running in cluster or locally
func RunningInCluster() bool {
	_, err := os.Stat(inClusterNamespacePath)
	return !os.IsNotExist(err)
}
