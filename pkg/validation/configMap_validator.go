package validation

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

type dataValidator struct {
	Logger logrus.FieldLogger
}

var _ configMapValidator = (*dataValidator)(nil)

func (n dataValidator) Validate(cm *corev1.ConfigMap) (validation, error) {
	newData := cm.Data
	start, err1 := strconv.Atoi(newData["start"])
	end, err2 := strconv.Atoi(newData["end"])
	portsPerNode, err3 := strconv.Atoi(newData["ports-per-node"])
	if err1 != nil || err2 != nil || err3 != nil ||
		start < 5000 || end > 65000 || start > end || portsPerNode > end-start+1 {
		v := validation{
			Valid:  false,
			Reason: fmt.Sprintf("configMap data contains invalid values...start=%q, end=%q portsPerNode=%q", start, end, portsPerNode),
		}
		return v, nil
	}

	return validation{Valid: true, Reason: "valid data"}, nil
}
