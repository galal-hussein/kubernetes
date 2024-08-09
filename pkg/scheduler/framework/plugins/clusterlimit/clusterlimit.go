package clusterlimit

import (
	"context"
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const (
	Name = "ClusterLimit"
)

var scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

type ClusterLimit struct {
	sync.RWMutex
	fh         framework.Handle
	podLister  corelisters.PodLister
	nodeLister corelisters.NodeLister
}

func (c *ClusterLimit) Name() string {
	return Name
}

func New(ctx context.Context, obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {

	logger := klog.FromContext(ctx)

	// cast obj to clusterLimitArgs
	logger.V(5).Info("creating new ClusterLimit plugin")
	// clusterLimitArgs, err := getArgs(obj)
	// if err != nil {
	// 	return nil, err
	// }
	c := &ClusterLimit{
		fh:         handle,
		podLister:  handle.SharedInformerFactory().Core().V1().Pods().Lister(),
		nodeLister: handle.SharedInformerFactory().Core().V1().Nodes().Lister(),
	}

	return c, nil
}

// // getArgs : returns the arguments for the SySchedArg plugin.
// func getArgs(obj runtime.Object) (*pluginconfig.ClusterLimitArgs, error) {
// 	ClusterLimitArgs, ok := obj.(*pluginconfig.ClusterLimitArgs)
// 	if !ok {
// 		return nil, fmt.Errorf("want args to be of type ClusterLimitArgs, got %T", obj)
// 	}

// 	return ClusterLimitArgs, nil
// }

// implementing PreFilter plugin interface for k3k-scheduler
// 1. Check if sum of all (pod.limits + allPods.limits) < cluster.limit (configuration)
func (c *ClusterLimit) PreFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod) (*framework.PreFilterResult, *framework.Status) {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()
	podLimits := computePodResourceLimits(pod)

	allPodLimits := &framework.Resource{}
	// list pods across all namespaces
	podList, err := c.podLister.List(labels.NewSelector())
	if err != nil {
		return nil, framework.NewStatus(framework.Error, err.Error())
	}
	for _, obj := range podList {
		if obj.Name == pod.Name {
			continue
		}
		// pod havent yet scheduled
		if obj.Spec.NodeName == "" {
			continue
		}
		objLimits := computePodResourceLimits(obj)
		allPodLimits = addResourcePodLimit(allPodLimits, objLimits)
	}
	// klog.Infof("cluster limit args: %#v", c.clusterLimit.ControlPlane["memory"])
	klog.Infof("cluster limit prefilter for pod %s", pod.Name)
	klog.Infof("pod [%s] has memory limit %d MB, all pods limits %d MB", pod.Name, podLimits.Memory/1024/1024, allPodLimits.Memory/1024/1024)

	// for resource, limit := range c.clusterLimit.ControlPlane {
	// 	switch resource {
	// 	case string(v1.ResourceMemory):
	// 		klog.Infof("limit %d compared to cluster limit %d", podLimits.Memory+allPodLimits.Memory, limit.Value())
	// 		if podLimits.Memory+allPodLimits.Memory > limit.Value() {
	// 			return nil, framework.NewStatus(framework.Unschedulable, "Memory limit exceeded")
	// 		}
	// 	case string(v1.ResourceCPU):
	// 		if podLimits.MilliCPU+allPodLimits.MilliCPU > limit.Value() {
	// 			return nil, framework.NewStatus(framework.Unschedulable, "CPU limit exceeded")
	// 		}
	// 	case string(v1.ResourceEphemeralStorage):
	// 		if podLimits.EphemeralStorage+allPodLimits.EphemeralStorage > limit.Value() {
	// 			return nil, framework.NewStatus(framework.Unschedulable, "Ephermal storage limit exceeded")
	// 		}
	// 	}
	// }
	return nil, framework.NewStatus(framework.Success, "")
}

func (c *ClusterLimit) PreFilterExtensions() framework.PreFilterExtensions {
	return c
}

func (c *ClusterLimit) AddPod(ctx context.Context, cycleState *framework.CycleState, podToSchedule *v1.Pod, podToAdd *framework.PodInfo, nodeInfo *framework.NodeInfo) *framework.Status {
	return framework.NewStatus(framework.Success, "")
}

func (c *ClusterLimit) RemovePod(ctx context.Context, cycleState *framework.CycleState, podToSchedule *v1.Pod, podToRemove *framework.PodInfo, nodeInfo *framework.NodeInfo) *framework.Status {
	return framework.NewStatus(framework.Success, "")
}

func computePodResourceLimits(pod *v1.Pod) *framework.Resource {

	result := &framework.Resource{}
	for _, container := range pod.Spec.Containers {
		result.Add(container.Resources.Limits)
	}

	for _, container := range pod.Spec.InitContainers {
		result.Add(container.Resources.Limits)
	}
	return result
}

func addResourcePodLimit(sum, pod *framework.Resource) *framework.Resource {
	if pod == nil {
		return nil
	}
	sum.MilliCPU += pod.MilliCPU
	sum.Memory += pod.Memory
	sum.AllowedPodNumber += pod.AllowedPodNumber
	sum.EphemeralStorage += pod.EphemeralStorage

	return sum
}
