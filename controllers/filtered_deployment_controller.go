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
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// FilteredDeploymentController reconciles a Deployment object that has been filtered
// out using predicates.
type FilteredDeploymentController struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments/finalizers,verbs=update
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *FilteredDeploymentController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	instanceKey := req.NamespacedName
	l := log.FromContext(ctx).WithValues("deployment", instanceKey, "reconciler", "FilteredDeployment")
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

	// Instead of filtering here as in deployment_reconciler.go, this controller will never see
	// events that don't have the desired annotation.
	var four int32 = 4
	l.Info("Setting replica count to 4 if it isn't already 4")
	patchDiff := client.MergeFrom(existing.DeepCopy())
	existing.Spec.Replicas = &four
	if err = r.Patch(ctx, &existing, patchDiff); err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{Requeue: false}, nil // success
}

// SetupWithManager sets up the controller with the Manager.
func (r *FilteredDeploymentController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(&hasAskedNicelyInAnnotationsPredicate{}).
		Complete(r)
}

// hasAskedNicelyInAnnotationsPredicate is an event filter that that will only allow update and
// create events for resources that have the annotation "scalem" with the value "please".
type hasAskedNicelyInAnnotationsPredicate struct{}

var _ predicate.Predicate = &hasAskedNicelyInAnnotationsPredicate{}

// Update events will trigger this reconciler but only if the appropriate filter is set.
func (p *hasAskedNicelyInAnnotationsPredicate) Update(e event.UpdateEvent) bool {
	return hasAnnotationWithValue(e.ObjectNew.GetAnnotations(), "scaleme", "please")
}

// Create events will trigger this reconciler but only if the appropriate filter is set.
func (p *hasAskedNicelyInAnnotationsPredicate) Create(e event.CreateEvent) bool {
	return hasAnnotationWithValue(e.Object.GetAnnotations(), "scaleme", "please")
}

// Delete returns false - we do not care about Delete events.
func (p *hasAskedNicelyInAnnotationsPredicate) Delete(e event.DeleteEvent) bool {
	return false
}

// Generic returns false - we do not care about Generic events.
func (p *hasAskedNicelyInAnnotationsPredicate) Generic(e event.GenericEvent) bool {
	return false
}
