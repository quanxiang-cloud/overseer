package overseerrun

import (
	v1alpha1 "github.com/quanxiang-cloud/overseer/pkg/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// failedWithError filedc with error.
func failedWithError(osrs *v1alpha1.OverseerRunStatus, err error) {
	if err == nil {
		return
	}

	osrs.Status = corev1.ConditionFalse
	osrs.Reason = "Error"
	osrs.Message = err.Error()
}
