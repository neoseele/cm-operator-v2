package controller

import (
	"github.com/neoseele/cm-operator/pkg/controller/custommetric"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, custommetric.Add)
}
