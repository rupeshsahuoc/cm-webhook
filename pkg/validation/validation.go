package validation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type Validator struct {
	Logger *logrus.Entry
}

func NewValidator(logger *logrus.Entry) *Validator {
	return &Validator{Logger: logger}
}

type configMapValidator interface {
	Validate(*corev1.ConfigMap) (validation, error)
}

type validation struct {
	Valid  bool
	Reason string
}

func (v *Validator) ValidateConfigMap(cm *corev1.ConfigMap) (validation, error) {
	var cmName string
	if cm.Name != "" {
		cmName = cm.Name
	} else {
		if cm.ObjectMeta.GenerateName != "" {
			cmName = cm.ObjectMeta.GenerateName
		}
	}
	log := logrus.WithField("configMap_name", cmName)
	log.Print("delete me")

	// list of all validations to be applied to the configMap
	validations := []configMapValidator{
		dataValidator{v.Logger},
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

	return validation{Valid: true, Reason: "valid configMap"}, nil
}
