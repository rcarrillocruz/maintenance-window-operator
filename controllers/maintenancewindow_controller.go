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

package controllers

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	windowv1alpha1 "github.com/stolostron/maintenance-window-operator/api/v1alpha1"
)

// MaintenanceWindowReconciler reconciles a MaintenanceWindow object
type MaintenanceWindowReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=window.open-cluster-management.io,resources=maintenancewindows,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=window.open-cluster-management.io,resources=maintenancewindows/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=window.open-cluster-management.io,resources=maintenancewindows/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MaintenanceWindow object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MaintenanceWindowReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var maintenanceWindow windowv1alpha1.MaintenanceWindow
	err := r.Get(ctx, req.NamespacedName, &maintenanceWindow)
	if err != nil {
		log.Log.Error(err, "unable to get MaintenanceWindow")
		return ctrl.Result{}, nil
	}

	startDate, err := time.Parse("2006-01-02", maintenanceWindow.Spec.Date)
	if err != nil {
		log.Log.Error(err, "unable to parse startDate")
		return ctrl.Result{}, nil
	}
	startTime, err := time.Parse(time.Kitchen, maintenanceWindow.Spec.Time)
	if err != nil {
		log.Log.Error(err, "unable to parse startTime")
		return ctrl.Result{}, nil
	}
	location, err := time.LoadLocation(maintenanceWindow.Spec.TimeZone)
	if err != nil {
		log.Log.Error(err, "unable to load location for given timezone")
	}
	builtTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), startTime.Hour(), startTime.Minute(), 0, 0, location)
	log.Log.Info("DEBUG", "built time", builtTime.String())
	log.Log.Info("DEBUG", "time now", metav1.Now())

	diff := builtTime.UTC().Sub(time.Now().UTC())
	log.Log.Info("DEBUG", "time till maintenance window", diff)
	if diff > 0 {
		log.Log.Info("DEBUG: Maintenance window has not yet started")
		log.Log.Info("DEBUG", "diff", diff)
		return ctrl.Result{RequeueAfter: diff}, nil
	} else {
		if time.Since(builtTime) > time.Duration(*maintenanceWindow.Spec.Duration)*time.Second {
			log.Log.Info("DEBUG", "fistParam", time.Since(builtTime).String())
			log.Log.Info("DEBUG", "secondParam", time.Duration(*maintenanceWindow.Spec.Duration)*time.Second)
			log.Log.Info("DEBUG: Maintenance window is closed")
		} else {
			log.Log.Info("DEBUG: Maintenance window now in place")
			return ctrl.Result{RequeueAfter: time.Duration(*maintenanceWindow.Spec.Duration) * time.Second}, nil
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MaintenanceWindowReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&windowv1alpha1.MaintenanceWindow{}).
		Complete(r)
}
