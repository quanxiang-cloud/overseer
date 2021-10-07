package v1alpha1

import (
	"strings"

	osv1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	pipeline1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var depot map[schema.GroupVersionKind]Getter

func init() {
	depot = make(map[schema.GroupVersionKind]Getter)
	depot[corev1.SchemeGroupVersion.WithKind("PersistentVolume")] = &PersistentVolume{}
	depot[pipeline1beta1.SchemeGroupVersion.WithKind("TaskRun")] = &TaskRun{}
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
		sc.State = osv1alpha1.StepConditionSuceess
	}

	sc.Message = o.Status.Message
	sc.Reason = o.Status.Reason
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

	if len(o.Status.Conditions) == 0 {
		return sc
	}
	condition := o.Status.Conditions[0]

	switch condition.Status {
	case corev1.ConditionFalse:
		sc.State = osv1alpha1.StepConditionFail
	case corev1.ConditionTrue:
		sc.State = osv1alpha1.StepConditionSuceess
	default:
		sc.State = osv1alpha1.StepConditionUnknown
	}

	sc.Message = condition.Message
	sc.Reason = condition.Reason
	return sc
}
