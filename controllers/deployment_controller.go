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

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instanceKey := req.NamespacedName
	l := log.FromContext(ctx).WithValues("deployment", instanceKey)
	l.Info("reconciliation started for Deployment")

	// If the deployment exists, get it and patch it
	var existing appsv1.Deployment
	err := r.Client.Get(ctx, instanceKey, &existing)

	if apierrors.IsNotFound(err) {
		// for this demo, do nothing.
		return ctrl.Result{}, nil
	}

	if err != nil {
		// couldn't get the deployment so try again.
		return ctrl.Result{Requeue: true}, err
	}

	if hasAnnotationWithValue(existing.Annotations, "foo", "bar") {
		l.Info("Deployment has the expected annotation, so setting replicas to 5 if not already done")
		patchDiff := client.MergeFrom(existing.DeepCopy())
		var five int32 = 5
		existing.Spec.Replicas = &five
		if err = r.Patch(ctx, &existing, patchDiff); err != nil {
			return ctrl.Result{Requeue: true}, err
		}
	} else {
		l.Info("Deployment did not have the expected annotation")
	}

	return ctrl.Result{Requeue: false}, nil // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		Complete(r)
}

// hasAnnotationWithValue returns true if the key exists in a and has value desiredValue.
func hasAnnotationWithValue(a map[string]string, key, desiredValue string) bool {
	if actualV, ok := a[key]; ok {
		return actualV == desiredValue
	}

	return false
}
