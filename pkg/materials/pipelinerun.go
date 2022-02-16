package materials

import (
	"context"
	"sort"
	"strings"
	"time"

	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/apis/overseer/v1alpha1"
	pipelinev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	knativeapi "knative.dev/pkg/apis"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PipelineRun struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=tekton.dev,resources=pipelineruns;taskruns,verbs=get;list;watch;create;update;patch;delete

func (r *PipelineRun) Reconcile(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, pplrSpec *v1alpha1.PipelineRunSpec) error {
	logger := log.FromContext(ctx)
	logger = logger.WithName("reconcilePipelineRun")

	pplr := &pipelinev1beta1.PipelineRun{}

	name := genName(osr.Name, vst.Name)
	err := r.Get(ctx, types.NamespacedName{Namespace: osr.Namespace, Name: name}, pplr)
	if errors.IsNotFound(err) {
		// the pipeline has not been executed,
		// create and executed it.
		pplr, err = r.createPipeline(ctx, osr, vst, pplrSpec)
		if err != nil {
			return err
		}
	}
	if err != nil {
		logger.Error(err, "get pipelineRun", "name", name)
		return err
	}

	status := v1alpha1.VersatileStatus{}
	status.Ref = pplr.Name
	status.StartTime = metav1.NewTime(time.Now())
	status.Status = corev1.ConditionUnknown

	pipelineRuns := v1alpha1.PipelineRuns{}
	pipelineRuns.Ref = pplr.Name
	pipelineRuns.Status = corev1.ConditionUnknown

	cond := pplr.GetStatusCondition().GetCondition(knativeapi.ConditionSucceeded)
	if cond != nil {
		status.Conditions = v1alpha1.Conditions{
			{
				Type:               v1alpha1.Succeeded,
				Status:             cond.Status,
				LastTransitionTime: cond.LastTransitionTime.Inner,
				Reason:             cond.Reason,
				Message:            cond.Message,
			},
		}
		status.Status = cond.Status
		pipelineRuns.Status = cond.Status
	}

	if status.Status == corev1.ConditionFalse {
		osr.Status.SetFalse()
	}

	tasks := pplr.Status.PipelineRunStatusFields.TaskRuns
	for _, task := range tasks {
		taskRuns := v1alpha1.TaskRuns{
			Name:           task.PipelineTaskName,
			CompletionTime: task.Status.GetCondition(knativeapi.ConditionSucceeded).LastTransitionTime.Inner,
			Steps:          make([]v1alpha1.Steps, 0, len(task.Status.Steps)),
		}
		if task.Status != nil {
			taskRuns.StartTime = *task.Status.StartTime
		}

		taskCond := task.Status.GetCondition(knativeapi.ConditionSucceeded)
		if taskCond != nil {
			taskRuns.Status = taskCond.Status
		}
		for _, element := range task.Status.Steps {
			step := v1alpha1.Steps{Name: element.Name}
			if element.Terminated != nil {
				step.StartTime = element.Terminated.StartedAt
				step.CompletionTime = element.Terminated.FinishedAt
			}
			taskRuns.Steps = append(taskRuns.Steps, step)
		}

		pipelineRuns.CompletionTime = metav1.NewTime(time.Now())
		pipelineRuns.TaskRuns = append(pipelineRuns.TaskRuns, taskRuns)
	}

	sort.SliceStable(pipelineRuns.TaskRuns, func(i, j int) bool {
		return pipelineRuns.TaskRuns[i].StartTime.Before(&pipelineRuns.TaskRuns[j].StartTime)
	})

	osr.Status.SetVersatileStatus(status)
	osr.Status.SetPipelineRuns(pipelineRuns)
	return nil
}

func (r *PipelineRun) createPipeline(ctx context.Context, osr *v1alpha1.Overseer, vst *v1alpha1.Versatile, pplrSpec *v1alpha1.PipelineRunSpec) (*pipelinev1beta1.PipelineRun, error) {
	logger := log.FromContext(ctx)
	logger = logger.WithName("createPipeline")

	pplr := &pipelinev1beta1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			Name:      genName(osr.Name, vst.Name),
			Namespace: osr.Namespace,
			Labels:    osr.Labels,
		},
		Spec: pipelinev1beta1.PipelineRunSpec{
			PipelineRef: &pipelinev1beta1.PipelineRef{
				Name: pplrSpec.PipelineRef,
			},
			ServiceAccountName: osr.Spec.ServiceAccountName,
			Params:             make([]pipelinev1beta1.Param, 0, len(pplrSpec.Params)),
			Workspaces:         make([]pipelinev1beta1.WorkspaceBinding, 0),
		},
	}

	for _, param := range pplrSpec.Params {
		pplr.Spec.Params = append(pplr.Spec.Params, pipelinev1beta1.Param{
			Name:  param.Name,
			Value: *pipelinev1beta1.NewArrayOrString(param.Value),
		})
	}

	for _, workspace := range pplrSpec.Workspace {
		wsb := pipelinev1beta1.WorkspaceBinding{}
		for _, volume := range osr.Spec.Volumes {
			if volume.Name == workspace.Name {
				wsb.Name = workspace.Name
				wsb.SubPath = workspace.SubPath
				// wsb.VolumeClaimTemplate = volume.PersistentVolumeClaim
				wsb.PersistentVolumeClaim = volume.PersistentVolumeClaim
				wsb.EmptyDir = volume.EmptyDir
				wsb.ConfigMap = volume.ConfigMap
				wsb.Secret = volume.Secret
				break
			}
		}

		if wsb.Name == "" {
			osr.Status.SetFalse()
			return nil, nil
		}

		pplr.Spec.Workspaces = append(pplr.Spec.Workspaces, wsb)
	}

	if pplrSpec.Template != nil {
		pplr.Spec.PodTemplate = &pipelinev1beta1.PodTemplate{
			NodeSelector:                 pplrSpec.Template.NodeSelector,
			Tolerations:                  pplrSpec.Template.Tolerations,
			Affinity:                     pplrSpec.Template.Affinity,
			SecurityContext:              pplrSpec.Template.SecurityContext,
			Volumes:                      pplrSpec.Template.Volumes,
			RuntimeClassName:             pplrSpec.Template.RuntimeClassName,
			AutomountServiceAccountToken: pplrSpec.Template.AutomountServiceAccountToken,
			DNSPolicy:                    &pplrSpec.Template.DNSPolicy,
			DNSConfig:                    pplrSpec.Template.DNSConfig,
			EnableServiceLinks:           pplrSpec.Template.EnableServiceLinks,
			PriorityClassName:            &pplrSpec.Template.PriorityClassName,
			SchedulerName:                pplrSpec.Template.SchedulerName,
			ImagePullSecrets:             pplrSpec.Template.ImagePullSecrets,
			HostAliases:                  pplrSpec.Template.HostAliases,
			HostNetwork:                  pplrSpec.Template.HostNetwork,
		}
	}

	pplr.SetOwnerReferences(nil)
	if err := ctrl.SetControllerReference(osr, pplr, r.Scheme); err != nil {
		logger.Error(err, "Failed to SetControllerReference", "name", pplr.Name)
		return nil, err
	}

	if err := r.Create(ctx, pplr); err != nil {
		logger.Error(err, "Failed to create obj")
		return nil, err
	}

	return pplr, nil
}

func genName(str ...string) string {
	return strings.Join(str, "-")
}
