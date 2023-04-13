package controllers

import (
	"context"
	"time"

	bootstrapv1 "github.com/charmed-kubernetes/cluster-api-bootstrap-provider-charmed-k8s/api/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

var _ = Describe("CharmedK8sConfig", func() {
	const (
		duration = time.Second * 10
		interval = time.Millisecond * 250
		timeout  = time.Second * 10
	)

	DescribeTable("With a machine owner",
		func(isControlPlaneOwner bool) {
			ctx := context.Background()

			cluster := &clusterv1.Cluster{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "cluster.x-k8s.io/v1beta1",
					Kind:       "Cluster",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cluster",
					Namespace: "default",
				},
				Spec: clusterv1.ClusterSpec{},
			}
			Expect(k8sClient.Create(ctx, cluster)).Should(Succeed())

			machineLabels := map[string]string{}
			if isControlPlaneOwner {
				machineLabels[clusterv1.MachineControlPlaneLabelName] = "true"
			}
			machine := &clusterv1.Machine{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "cluster.x-k8s.io/v1beta1",
					Kind:       "Machine",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-machine",
					Namespace: "default",
					Labels:    machineLabels,
				},
				Spec: clusterv1.MachineSpec{
					ClusterName: "test-cluster",
				},
			}
			Expect(k8sClient.Create(ctx, machine)).Should(Succeed())

			charmedK8sConfig := &bootstrapv1.CharmedK8sConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "bootstrap.cluster.x-k8s.io/v1beta1",
					Kind:       "CharmedK8sConfig",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-charmedk8sconfig",
					Namespace: "default",
					OwnerReferences: []metav1.OwnerReference{
						{
							APIVersion: "cluster.x-k8s.io/v1beta1",
							Kind:       "Machine",
							Name:       "test-machine",
							UID:        machine.GetUID(),
						},
					},
				},
				Spec: bootstrapv1.CharmedK8sConfigSpec{
					ControlPlaneApplications: []string{"etcd", "kubernetes-control-plane"},
					WorkerApplications:       []string{"kubernetes-worker"},
				},
			}
			Expect(k8sClient.Create(ctx, charmedK8sConfig)).Should(Succeed())

			Eventually(func() bool {
				lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
				err := k8sClient.Get(ctx, lookupKey, charmedK8sConfig)
				if err != nil {
					return false
				}
				return charmedK8sConfig.Status.Ready
			}, timeout, interval).Should(BeTrue())

			Expect(charmedK8sConfig.Status.DataSecretName).Should(Equal("test-charmedk8sconfig"))

			secret := &corev1.Secret{}
			lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
			Expect(k8sClient.Get(ctx, lookupKey, secret)).Should(Succeed())
			Expect(secret.Data["format"]).Should(Equal([]byte("juju")))
			if isControlPlaneOwner {
				Expect(secret.Data["value"]).Should(Equal([]byte("[\"etcd\",\"kubernetes-control-plane\"]")))
			} else {
				Expect(secret.Data["value"]).Should(Equal([]byte("[\"kubernetes-worker\"]")))
			}
			Expect(secret.ObjectMeta.OwnerReferences).Should(Equal([]metav1.OwnerReference{
				{
					APIVersion: bootstrapv1.GroupVersion.String(),
					Kind:       "CharmedK8sConfig",
					Name:       "test-charmedk8sconfig",
					UID:        charmedK8sConfig.UID,
					Controller: pointer.Bool(true),
				},
			}))

			Expect(k8sClient.Delete(ctx, charmedK8sConfig)).Should(Succeed())
			Eventually(func() bool {
				lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
				err := k8sClient.Get(ctx, lookupKey, charmedK8sConfig)
				return apierrors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())

			if useExistingCluster {
				// This is a real cluster, so the secret should get garbage collected
				Eventually(func() bool {
					lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
					err := k8sClient.Get(ctx, lookupKey, secret)
					return apierrors.IsNotFound(err)
				}, timeout, interval).Should(BeTrue())
			} else {
				// Not a real cluster, so there is no garbage collection for orphaned
				// secrets. Delete it.
				Expect(k8sClient.Delete(ctx, secret)).Should(Succeed())
			}
			Expect(k8sClient.Delete(ctx, machine)).Should(Succeed())
			Expect(k8sClient.Delete(ctx, cluster)).Should(Succeed())
		},
		Entry("Has a worker machine owner", false),
		Entry("Has a control plane machine owner", true),
	)

	Describe("With no owner", func() {
		It("Never becomes ready", func() {
			ctx := context.Background()

			charmedK8sConfig := &bootstrapv1.CharmedK8sConfig{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "bootstrap.cluster.x-k8s.io/v1beta1",
					Kind:       "CharmedK8sConfig",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-charmedk8sconfig",
					Namespace: "default",
				},
				Spec: bootstrapv1.CharmedK8sConfigSpec{
					ControlPlaneApplications: []string{"etcd", "kubernetes-control-plane"},
					WorkerApplications:       []string{"kubernetes-worker"},
				},
			}
			Expect(k8sClient.Create(ctx, charmedK8sConfig)).Should(Succeed())

			Eventually(func() bool {
				lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
				err := k8sClient.Get(ctx, lookupKey, charmedK8sConfig)
				return err == nil
			}, timeout, interval).Should(BeTrue())
			Consistently(func() bool {
				lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
				err := k8sClient.Get(ctx, lookupKey, charmedK8sConfig)
				if err != nil {
					return false
				}
				return !charmedK8sConfig.Status.Ready && charmedK8sConfig.Status.DataSecretName == ""
			}, duration, interval).Should(BeTrue())

			Expect(k8sClient.Delete(ctx, charmedK8sConfig)).Should(Succeed())
			Eventually(func() bool {
				lookupKey := types.NamespacedName{Name: "test-charmedk8sconfig", Namespace: "default"}
				err := k8sClient.Get(ctx, lookupKey, charmedK8sConfig)
				return apierrors.IsNotFound(err)
			}, timeout, interval).Should(BeTrue())
		})
	})
})
