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
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	lzv1alpha1 "github.com/rexagod/lazarus/api/v1alpha1"
)

// LTargetReconciler reconciles a LTarget object
type LTargetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=lz.rexa.god,resources=ltargets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lz.rexa.god,resources=ltargets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lz.rexa.god,resources=ltargets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=services,verbs=get;watch;create

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *LTargetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Reconciling LTarget")
	lzList := lzv1alpha1.LTargetList{}
	err := r.List(ctx, &lzList, client.InNamespace(req.Namespace))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list LTarget: %w", err)
	}
	if len(lzList.Items) == 0 {
		return ctrl.Result{}, nil
	}
	if len(lzList.Items) > 1 {
		return ctrl.Result{}, fmt.Errorf("found multiple LTarget objects, only one is allowed inside the " +
			"same namespace") // TODO(@rexagod): handle multiple LTarget objects (in different namespaces?)
	}
	lzTarget := lzList.Items[0]
	lzTarget.Status.ConnectionStatus = "cr found"
	l.Info("Reconciling LTarget", "name", lzTarget.Name)
	lzService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Lazarus Media Service",
			Namespace: req.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "lz.rexa.god/v1alpha1",
					Kind:       "LTarget",
					Name:       lzTarget.Name,
					UID:        lzTarget.UID,
				},
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "Lazarus Mediator Service Port",
					Port:       lzTarget.Spec.ExternalDelvePort, // TODO(@rexagod): multi-port service?
					TargetPort: lzTarget.Spec.InternalDelvePortOrName,
					Protocol:   "TCP",
				},
			},
			Type:     "LoadBalancer",
			Selector: lzTarget.Spec.LTargetLabel,
		},
	}
	existingLzService := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Lazarus Media Service",
			Namespace: req.Namespace,
		},
	}
	err = r.Get(ctx, req.NamespacedName, &existingLzService)
	if err == nil {
		l.Info("Lazarus Media Service already exists")
		lzTarget.Status.ConnectionStatus = "service found"
		lzService = existingLzService.DeepCopy()
	} else if err != nil && !errors.IsNotFound(err) {
		l.Info("Cannot get existing LTarget service")
		lzTarget.Status.ConnectionStatus = "cannot fetch service"
		return ctrl.Result{}, fmt.Errorf("failed to get LTarget service: %w", err)
	} else {
		l.Info("LTarget service not found, creating a new one")
		lzTarget.Status.ConnectionStatus = "creating service"
		err = r.Create(ctx, lzService)
		if err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to create LTarget service: %w", err)
		}
		lzTarget.Status.ConnectionStatus = "service created"
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LTargetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lzv1alpha1.LTarget{}).
		Complete(r)
}
