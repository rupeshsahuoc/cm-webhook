package validation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// Validator is a container for mutation
type Validator struct {
	Logger *logrus.Entry
}

// NewValidator returns an initialised instance of Validator
func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

// podValidators is an interface used to group functions mutating pods
type podValidator interface {
	Validate(*corev1.ConfigMap) (validation, error)
	Name() string
}

type validation struct {
	Valid  bool
	Reason string
}

// ValidatePod returns true if a pod is valid
func (v *Validator) ValidateConfigMap(cm *corev1.ConfigMap) (validation, error) {
	var podName string
	if cm.Name != "" {
		podName = cm.Name
	} else {
		if cm.ObjectMeta.GenerateName != "" {
			podName = cm.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("pod_name", podName)
	log.Print("delete me")

	// list of all validations to be applied to the pod
	validations := []podValidator{
		nameValidator{v.Logger},
	}

	// apply all validations
	for _, v := range validations {
		var err error
		vp, err := v.Validate(cm)
		if err != nil {
			return validation{Valid: false, Reason: err.Error()}, err
		}
		if !vp.Valid {
			return validation{Valid: false, Reason: vp.Reason}, err
		}
	}

	return validation{Valid: true, Reason: "valid pod"}, nil
}
