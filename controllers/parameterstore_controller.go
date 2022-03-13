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

	errs "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	logf "sigs.k8s.io/controller-runtime/pkg/log"

	ssmv1alpha1 "github.com/fr123k/aws-ssm-operator/api/v1alpha1"
)

var log = logf.Log.WithName("parameterstore-controller")

// ParameterStoreReconciler reconciles a ParameterStore object
type ParameterStoreReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	ssmc *SSMClient
}

//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ssm.aws,resources=parameterstores,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ssm.aws,resources=parameterstores/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ssm.aws,resources=parameterstores/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ParameterStore object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ParameterStoreReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	reqLogger.Info("Reconciling ParameterStore")

	// Fetch the ParameterStore instance
	instance := &ssmv1alpha1.ParameterStore{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile req.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the req.
		return reconcile.Result{}, err
	}

	// Define a new Secret object
	desired, err := r.newSecretForCR(instance)
	if err != nil {
		return reconcile.Result{}, errs.Wrap(err, "failed to compute secret for cr")
	}

	// Set ParameterStore instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, desired, r.Scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Secret already exists
	current := &corev1.Secret{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, current)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Creating a new Secret", "desired.Namespace", desired.Namespace, "desired.Name", desired.Name)
			err = r.Client.Create(context.TODO(), desired)
		}
	} else {
		reqLogger.Info("Updating an existing Secret", "desired.Namespace", desired.Namespace, "desired.Name", desired.Name)
		err = r.Client.Update(context.TODO(), desired)
	}

	return reconcile.Result{}, err
}

// newSecretForCR returns a Secret with the same name/namespace as the cr
func (r *ParameterStoreReconciler) newSecretForCR(cr *ssmv1alpha1.ParameterStore) (*corev1.Secret, error) {
	labels := map[string]string{
		"app": cr.Name,
	}
	if r.ssmc == nil {
		r.ssmc = newSSMClient(nil)
	}
	ref := cr.Spec.ValueFrom.ParameterStoreRef
	data1, err := r.ssmc.SSMParameterValueToSecret(ref)

	if err != nil {
		return nil, errs.Wrap(err, "failed to get json secret as map")
	}

	data2, anno, err := r.ssmc.SSMParametersValueToSecret(cr.Spec.ValueFrom.ParametersStoreRef)

	for k, v := range data1 {
		data2[k] = v
	}

	if err != nil {
		return nil, errs.Wrap(err, "failed to get json secret as map")
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        cr.Name,
			Namespace:   cr.Namespace,
			Labels:      labels,
			Annotations: anno,
		},
		StringData: data2,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ParameterStoreReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&ssmv1alpha1.ParameterStore{}).
		Complete(r)
}
