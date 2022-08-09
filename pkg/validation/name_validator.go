package validation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// nameValidator is a container for validating the name of pods
type nameValidator struct {
	Logger logrus.FieldLogger
}

// nameValidator implements the podValidator interface
var _ podValidator = (*nameValidator)(nil)

// Name returns the name of nameValidator
func (n nameValidator) Name() string {
	return "name_validator"
}

// Validate inspects the name of a given pod and returns validation.
// The returned validation is only valid if the pod name does not contain some
// bad string.
func (n nameValidator) Validate(cm *corev1.ConfigMap) (validation, error) {
	badString := "offensive"
	data := cm.Data
	start, err := strconv.Atoi(data["start"])
	logrus.Print("Here data checking.......")
	logrus.Print(start)
	if err != nil {
		logrus.Print(err.Error())
	}
	if strings.Contains(cm.Name, badString) {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("pod name contains %q", badString),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid name"}, nil
}
