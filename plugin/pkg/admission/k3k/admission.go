package k3k

import (
	"context"
	"fmt"
	"io"
	"time"

	"golang.org/x/sync/singleflight"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/utils/lru"
)

const (
	// PluginName indicates name of admission plugin.
	PluginName = "K3KCluster"
)

// K3KCluster enforces usage limits on a per resource basis in the namespace
type K3KCluster struct {
	*admission.Handler
	client kubernetes.Interface
	lister corev1listers.LimitRangeLister

	// liveLookups holds the last few live lookups we've done to help ammortize cost on repeated lookup failures.
	// This let's us handle the case of latent caches, by looking up actual results for a namespace on cache miss/no results.
	// We track the lookup result here so that for repeated requests, we don't look it up very often.
	liveLookupCache *lru.Cache
	group           singleflight.Group
	liveTTL         time.Duration
}

// Register registers a plugin
func Register(plugins *admission.Plugins) {
	plugins.Register(PluginName, func(config io.Reader) (admission.Interface, error) {
		return NewK3KCluster()
	})
}

// NewLimitRanger returns an object that enforces limits based on the supplied limit function
func NewK3KCluster() (*K3KCluster, error) {
	liveLookupCache := lru.New(10000)

	return &K3KCluster{
		Handler:         admission.NewHandler(admission.Create, admission.Update),
		liveLookupCache: liveLookupCache,
		liveTTL:         time.Duration(30 * time.Second),
	}, nil
}

var _ admission.ValidationInterface = &K3KCluster{}

//var _ admission.ValidationInterface = &K3KCluster{}

// Validate admits resources into cluster that do not violate any defined LimitRange in the namespace
func (k *K3KCluster) Validate(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) (err error) {
	if a.IsDryRun() {
		return nil
	}

	logger := klog.FromContext(ctx)
	if len(a.GetNamespace()) == 0 || a.GetKind().GroupKind() != api.Kind("Pod") {
		return nil
	}
	// make sure that schedulerName isnt set
	podObj := a.GetObject()
	pod := podObj.(*corev1.Pod)
	if pod.Spec.SchedulerName != "" && pod.Spec.SchedulerName != "default-scheduler" {
		logger.Info("Scheduler name is set to something other than default scheduler: %s", pod.Spec.SchedulerName)
		return admission.NewForbidden(a, fmt.Errorf("expected default scheduler, got: %T", pod.Spec.SchedulerName))
	}
	return nil
}
