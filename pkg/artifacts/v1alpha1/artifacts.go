package v1alpha1

import (
	"strings"

	osv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	pipeline1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	knativev1 "knative.dev/serving/pkg/apis/serving/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var depot map[schema.GroupVersionKind]Getter

func init() {
	depot = make(map[schema.GroupVersionKind]Getter)
	depot[corev1.SchemeGroupVersion.WithKind("Service")] = &Service{}
	depot[corev1.SchemeGroupVersion.WithKind("ConfigMap")] = &ConfigMap{}
	depot[corev1.SchemeGroupVersion.WithKind("PersistentVolume")] = &PersistentVolume{}
	depot[corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")] = &PersistentVolumeClaim{}

	depot[appsv1.SchemeGroupVersion.WithKind("Deployment")] = &Deployment{}

	depot[pipeline1beta1.SchemeGroupVersion.WithKind("TaskRun")] = &TaskRun{}
	depot[pipeline1beta1.SchemeGroupVersion.WithKind("PipelineRun")] = &PipelineRun{}
	depot[knativev1.SchemeGroupVersion.WithKind("Service")] = &KnativeService{}

}

func GetObj(gkv schema.GroupVersionKind) (client.Object, bool) {
	getter, ok := depot[gkv]
	if !ok {
		return nil, false
	}
	return getter.New(), true
}

func GetGetter(gkv string) (Getter, bool) {
	getter, ok := depot[ParseGroupVersionKind(gkv)]
	return getter, ok
}

func ParseGroupVersionKind(gvk string) schema.GroupVersionKind {
	if len(gvk) == 0 || gvk == "," {
		return schema.GroupVersionKind{}
	}

	i := strings.Index(gvk, ",")
	if i == 0 || i == len(gvk) {
		return schema.GroupVersionKind{}
	}

	gv, err := schema.ParseGroupVersion(gvk[:i])
	if err != nil {
		return schema.GroupVersionKind{}
	}

	i = strings.Index(gvk, "=")
	return gv.WithKind(gvk[i+1:])
}

type Getter interface {
	GetCondition(obj client.Object) osv1alpha1.RefCondition
	New() client.Object
}

type Service struct{}

func (s *Service) New() client.Object {
	return &corev1.Service{}
}

func (s *Service) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*corev1.Service)
	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
		Conditions:       make([]osv1alpha1.Condition, 0),
	}

	sc.Conditions = []osv1alpha1.Condition{
		{
			Status:             corev1.ConditionTrue,
			LastTransitionTime: o.CreationTimestamp,
		},
	}

	return sc
}

type ConfigMap struct {
}

func (c *ConfigMap) New() client.Object {
	return &corev1.ConfigMap{}
}

func (c *ConfigMap) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*corev1.ConfigMap)
	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	sc.Conditions = []osv1alpha1.Condition{
		{
			Status:             corev1.ConditionTrue,
			LastTransitionTime: o.CreationTimestamp,
		},
	}

	return sc
}

type PersistentVolume struct {
}

func (p *PersistentVolume) New() client.Object {
	return &corev1.PersistentVolume{}
}

func (p *PersistentVolume) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*corev1.PersistentVolume)

	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	var state corev1.ConditionStatus

	switch o.Status.Phase {
	case corev1.VolumeFailed:
		state = corev1.ConditionFalse
	case corev1.VolumePending:
		state = corev1.ConditionUnknown
	default:
		state = corev1.ConditionTrue
	}

	sc.Conditions = []osv1alpha1.Condition{
		{
			Status:             state,
			LastTransitionTime: o.CreationTimestamp,
			Message:            o.Status.Message,
			Reason:             o.Status.Reason,
		},
	}
	return sc
}

type PersistentVolumeClaim struct {
}

func (p *PersistentVolumeClaim) New() client.Object {
	return &corev1.PersistentVolumeClaim{}
}

func (p *PersistentVolumeClaim) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*corev1.PersistentVolumeClaim)

	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	var state corev1.ConditionStatus

	switch o.Status.Phase {
	case corev1.ClaimPending:
		state = corev1.ConditionUnknown
	case corev1.ClaimBound:
		state = corev1.ConditionTrue
	default:
		state = corev1.ConditionUnknown
	}

	sc.Conditions = []osv1alpha1.Condition{
		{
			Status:             state,
			LastTransitionTime: o.CreationTimestamp,
		},
	}

	for _, condition := range o.Status.Conditions {
		sc.Conditions = append(sc.Conditions, osv1alpha1.Condition{
			Status:             condition.Status,
			LastTransitionTime: condition.LastTransitionTime,
			Message:            condition.Message,
			Reason:             condition.Reason,
		})

	}

	return sc
}

type TaskRun struct {
}

func (t *TaskRun) New() client.Object {
	return &pipeline1beta1.TaskRun{}
}

func (t *TaskRun) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*pipeline1beta1.TaskRun)

	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
		Conditions:       make([]osv1alpha1.Condition, 0),
	}

	for _, condition := range o.Status.Conditions {
		sc.Conditions = append(sc.Conditions, osv1alpha1.Condition{
			Status:             condition.Status,
			LastTransitionTime: condition.LastTransitionTime.Inner,
			Message:            condition.Message,
			Reason:             condition.Reason,
		})

	}

	return sc
}

type PipelineRun struct{}

func (p *PipelineRun) New() client.Object {
	return &pipeline1beta1.PipelineRun{}
}

func (p *PipelineRun) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*pipeline1beta1.PipelineRun)

	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
		Conditions:       make([]osv1alpha1.Condition, 0),
	}

	for _, condition := range o.Status.Conditions {
		sc.Conditions = append(sc.Conditions, osv1alpha1.Condition{
			Status:             condition.Status,
			LastTransitionTime: condition.LastTransitionTime.Inner,
			Message:            condition.Message,
			Reason:             condition.Reason,
		})

	}
	return sc
}

type KnativeService struct{}

func (k *KnativeService) New() client.Object {
	return &knativev1.Service{}
}

func (k *KnativeService) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*knativev1.Service)

	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
		Conditions:       make([]osv1alpha1.Condition, 0),
	}

	for _, condition := range o.Status.Conditions {
		sc.Conditions = append(sc.Conditions, osv1alpha1.Condition{
			Status:             condition.Status,
			LastTransitionTime: condition.LastTransitionTime.Inner,
			Message:            condition.Message,
			Reason:             condition.Reason,
		})

	}

	return sc
}

type Deployment struct{}

func (d *Deployment) New() client.Object {
	return &appsv1.Deployment{}
}

func (d *Deployment) GetCondition(obj client.Object) osv1alpha1.RefCondition {
	o := obj.(*appsv1.Deployment)
	sc := osv1alpha1.RefCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
		Conditions:       make([]osv1alpha1.Condition, 0),
	}

	for _, condition := range o.Status.Conditions {
		sc.Conditions = append(sc.Conditions, osv1alpha1.Condition{
			Status:             condition.Status,
			LastTransitionTime: condition.LastTransitionTime,
			Message:            condition.Message,
			Reason:             condition.Reason,
		})

	}
	return sc
}
