package overseerrun

import (
	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// failedWithError filedc with error.
func failedWithError(status *v1alpha1.Status, err error) {
	if err == nil {
		return
	}

	status.Condition.Status = corev1.ConditionFalse
	status.Condition.Reason = "Error"
	status.Condition.Message = err.Error()
}
