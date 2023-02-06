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
	"encoding/json"
	"fmt"
	bootstrapv1 "github.com/charmed-kubernetes/cluster-api-bootstrap-provider-charmed-k8s/api/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	bsutil "sigs.k8s.io/cluster-api/bootstrap/util"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// CharmedK8sConfigReconciler reconciles a CharmedK8sConfig object
type CharmedK8sConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=charmedk8sconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=charmedk8sconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bootstrap.cluster.x-k8s.io,resources=charmedk8sconfigs/finalizers,verbs=update
//+kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status;machinesets;machines;machines/status;machinepools;machinepools/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=secrets;events;configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the CharmedK8sConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *CharmedK8sConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Lookup the charmed k8s config
	config := &bootstrapv1.CharmedK8sConfig{}
	if err := r.Client.Get(ctx, req.NamespacedName, config); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get config")
		return ctrl.Result{}, err
	}
	log.Info("Retrieved CharmedK8sConfig successfully")

	configOwner, err := bsutil.GetConfigOwner(ctx, r.Client, config)
	if apierrors.IsNotFound(err) {
		// Could not find the owner yet, this is not an error and will rereconcile when the owner gets set.
		log.Info("Config Owner could not be found yet ...")
		return ctrl.Result{}, nil
	}
	if err != nil {
		log.Error(err, "Failed to get Owner Config")
		return ctrl.Result{}, err
	}
	if configOwner == nil {
		log.Info("Config owner was nil")
		return ctrl.Result{}, nil
	}
	log.Info("Retrieved config owner successfully")

	log = log.WithValues("kind", configOwner.GetKind(), "version", configOwner.GetResourceVersion(), "name", configOwner.GetName())

	log = log.WithValues("Cluster", klog.KRef(configOwner.GetNamespace(), configOwner.ClusterName()))
	ctx = ctrl.LoggerInto(ctx, log)

	// Lookup the cluster the config owner is associated with
	cluster, err := util.GetClusterByName(ctx, r.Client, configOwner.GetNamespace(), configOwner.ClusterName())
	if err != nil {
		if errors.Cause(err) == util.ErrNoCluster {
			log.Info(fmt.Sprintf("%s does not belong to a cluster yet, waiting until it's part of a cluster", configOwner.GetKind()))
			return ctrl.Result{}, nil
		}

		if apierrors.IsNotFound(err) {
			log.Info("Cluster does not exist yet, waiting until it is created")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Could not get cluster with metadata")
		return ctrl.Result{}, err
	}
	log.Info("Retrieved cluster successfully", "cluster", cluster.ObjectMeta)

	if annotations.IsPaused(cluster, config) {
		log.Info("Reconciliation is paused for this object")
		return ctrl.Result{}, nil
	}

	dataSecretName := r.getDataSecretName(config)
	secret := &corev1.Secret{}
	err = r.Client.Get(ctx, client.ObjectKey{Namespace: config.Namespace, Name: dataSecretName}, secret)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// data secret doesn't exist yet, create it
			err := r.createDataSecret(ctx, config, configOwner, cluster)
			if err != nil {
				log.Error(err, "failed to create data secret")
				return ctrl.Result{}, err
			}
		} else {
			log.Error(err, "failed to get data secret")
			return ctrl.Result{}, err
		}
	}

	config.Status.DataSecretName = dataSecretName
	config.Status.Ready = true

	if err := r.Client.Status().Update(ctx, config); err != nil {
		log.Error(err, "failed to update CharmedK8sConfig status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CharmedK8sConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bootstrapv1.CharmedK8sConfig{}).
		Complete(r)
}

// create the data secret
func (r *CharmedK8sConfigReconciler) createDataSecret(ctx context.Context, config *bootstrapv1.CharmedK8sConfig, configOwner *bsutil.ConfigOwner, cluster *clusterv1.Cluster) error {
	log := ctrl.LoggerFrom(ctx)

	bootstrapData, err := r.newBootstrapData(config, configOwner)
	if err != nil {
		return errors.Wrapf(err, "failed to create bootstrap data")
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.getDataSecretName(config),
			Namespace: config.Namespace,
			Labels: map[string]string{
				clusterv1.ClusterLabelName: cluster.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: bootstrapv1.GroupVersion.String(),
					Kind:       "CharmedK8sConfig",
					Name:       config.Name,
					UID:        config.UID,
					Controller: pointer.Bool(true),
				},
			},
		},
		Data: map[string][]byte{
			"value":  []byte(bootstrapData),
			"format": []byte("juju"),
		},
		Type: clusterv1.ClusterSecretType,
	}

	if err := r.Client.Create(ctx, secret); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return errors.Wrapf(err, "failed to create data secret")
		}
		log.Info("data secret already exists, updating")
		if err := r.Client.Update(ctx, secret); err != nil {
			return errors.Wrapf(err, "failed to update data secret")
		}
	}

	return nil
}

func (r *CharmedK8sConfigReconciler) newBootstrapData(config *bootstrapv1.CharmedK8sConfig, configOwner *bsutil.ConfigOwner) ([]byte, error) {
	var applications []string

	if configOwner.IsControlPlaneMachine() {
		applications = config.Spec.ControlPlaneApplications
	} else {
		applications = config.Spec.WorkerApplications
	}

	bootstrapData, err := json.Marshal(applications)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to json encode bootstrap data")
	}

	return bootstrapData, nil
}

// get the data secret name that should go with this CharmedK8sConfig
// must be deterministic in case the CharmedK8sConfig.status is ever lost
func (r *CharmedK8sConfigReconciler) getDataSecretName(config *bootstrapv1.CharmedK8sConfig) string {
	return config.Name
}
