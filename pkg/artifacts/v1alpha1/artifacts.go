package v1alpha1

import (
	"strings"

	osv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	pipeline1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var depot map[schema.GroupVersionKind]Getter

func init() {
	depot = make(map[schema.GroupVersionKind]Getter)
	depot[corev1.SchemeGroupVersion.WithKind("Service")] = &Service{}
	depot[corev1.SchemeGroupVersion.WithKind("ConfigMap")] = &ConfigMap{}
	depot[corev1.SchemeGroupVersion.WithKind("PersistentVolumeClaim")] = &PersistentVolumeClaim{}
	depot[corev1.SchemeGroupVersion.WithKind("PersistentVolume")] = &PersistentVolume{}
	depot[pipeline1beta1.SchemeGroupVersion.WithKind("TaskRun")] = &TaskRun{}
	depot[pipeline1beta1.SchemeGroupVersion.WithKind("PipelineRun")] = &PipelineRun{}
	depot[appsv1.SchemeGroupVersion.WithKind("Deployment")] = &Deployment{}
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
	GetState(obj client.Object) osv1alpha1.StepCondition
	New() client.Object
}

type ConfigMap struct {
}

func (c *ConfigMap) New() client.Object {
	return &corev1.ConfigMap{}
}

func (c *ConfigMap) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*corev1.ConfigMap)
	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	sc.State = osv1alpha1.StepConditionSuccess
	return sc
}

type PersistentVolume struct {
}

func (p *PersistentVolume) New() client.Object {
	return &corev1.PersistentVolume{}
}

func (p *PersistentVolume) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*corev1.PersistentVolume)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	switch o.Status.Phase {
	case corev1.VolumeFailed:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.VolumePending:
		sc.State = osv1alpha1.StepConditionPending
	default:
		sc.State = osv1alpha1.StepConditionSuccess
	}

	sc.Message = o.Status.Message
	sc.Reason = o.Status.Reason
	return sc
}

type PersistentVolumeClaim struct {
}

func (p *PersistentVolumeClaim) New() client.Object {
	return &corev1.PersistentVolumeClaim{}
}

func (p *PersistentVolumeClaim) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*corev1.PersistentVolumeClaim)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	switch o.Status.Phase {
	case corev1.ClaimPending:
		sc.State = osv1alpha1.StepConditionPending
	case corev1.ClaimBound:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionFail
	}

	size := len(o.Status.Conditions)
	if size == 0 {
		return sc
	}

	condition := o.Status.Conditions[size-1]
	switch condition.Status {
	case corev1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason

	return sc
}

type TaskRun struct {
}

func (t *TaskRun) New() client.Object {
	return &pipeline1beta1.TaskRun{}
}

func (t *TaskRun) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*pipeline1beta1.TaskRun)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	size := len(o.Status.Conditions)
	if size == 0 {
		return sc
	}

	condition := o.Status.Conditions[size-1]

	switch condition.Status {
	case corev1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason
	return sc
}

type PipelineRun struct{}

func (p *PipelineRun) New() client.Object {
	return &pipeline1beta1.PipelineRun{}
}

func (p *PipelineRun) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*pipeline1beta1.PipelineRun)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	size := len(o.Status.Conditions)
	if size == 0 {
		return sc
	}

	condition := o.Status.Conditions[size-1]

	switch condition.Status {
	case corev1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason

	return sc
}

type Deployment struct{}

func (d *Deployment) New() client.Object {
	return &appsv1.Deployment{}
}

func (d *Deployment) GetState(obj client.Object) osv1alpha1.StepCondition {

	o := obj.(*appsv1.Deployment)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	size := len(o.Status.Conditions)
	if size == 0 {
		return sc
	}

	condition := o.Status.Conditions[size-1]

	switch condition.Status {
	case corev1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason

	return sc
}

type Service struct{}

func (s *Service) New() client.Object {
	return &corev1.Service{}
}

func (s *Service) GetState(obj client.Object) osv1alpha1.StepCondition {
	o := obj.(*corev1.Service)

	sc := osv1alpha1.StepCondition{
		GroupVersionKind: o.GroupVersionKind().String(),
	}

	size := len(o.Status.Conditions)
	if size == 0 {
		return sc
	}

	condition := o.Status.Conditions[size-1]

	switch condition.Status {
	case metav1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case metav1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuccess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason

	return sc
}
