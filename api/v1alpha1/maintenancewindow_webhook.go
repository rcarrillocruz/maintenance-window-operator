/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var maintenancewindowlog = logf.Log.WithName("maintenancewindow-resource")

func (r *MaintenanceWindow) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-window-open-cluster-management-io-v1alpha1-maintenancewindow,mutating=true,failurePolicy=fail,groups=window.open-cluster-management.io,resources=maintenancewindows,verbs=create;update,versions=v1alpha1,name=vmaintenancewindow.kb.io,sideEffects=None,admissionReviewVersions=v1

var _ webhook.Defaulter = &MaintenanceWindow{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *MaintenanceWindow) Default() {
	maintenancewindowlog.Info("default", "name", r.Name)

	if r.Spec.ChangeScope == "all" {
		labels := make(map[string]string)
		labels["maintenancewindows.window.open-cluster-management.io/scope"] = "all"
		r.SetLabels(labels)
	}

}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-window-open-cluster-management-io-v1alpha1-maintenancewindow,mutating=false,failurePolicy=fail,sideEffects=None,groups=window.open-cluster-management.io,resources=maintenancewindows,verbs=create;update;delete,versions=v1alpha1,name=vmaintenancewindow.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &MaintenanceWindow{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *MaintenanceWindow) ValidateCreate() error {
	maintenancewindowlog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *MaintenanceWindow) ValidateUpdate(old runtime.Object) error {
	maintenancewindowlog.Info("validate update", "name", r.Name)

	if r.Spec != old.(*MaintenanceWindow).Spec {
		return apierrors.NewBadRequest("MaintenanceWindow CR cannot be updated")
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *MaintenanceWindow) ValidateDelete() error {
	maintenancewindowlog.Info("validate delete", "name", r.Name)

	if r.Status.State == "OPENED" {
		return apierrors.NewBadRequest("MaintenanceWindow CR cannot be deleted while it is in OPENED state")
	}

	return nil
}
