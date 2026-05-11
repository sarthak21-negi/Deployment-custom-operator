/*
Copyright 2026.

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

package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	dpscalerv1alpha1 "github.com/sarthak21-negi/deployment-custom-operator/api/v1alpha1"
)

var _ = Describe("DeploymentCustomOperator Controller", func() {
	Context("When reconciling a resource", func() {
		const (
			resourceName   = "test-scaler"
			deploymentName = "test-nginx"
			namespace      = "default"
		)

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: namespace,
		}

		BeforeEach(func() {
			By("Creating test deployment")
			replicas := int32(1)
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      deploymentName,
					Namespace: namespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: &replicas,
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": deploymentName},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{"app": deploymentName},
						},
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "nginx",
									Image: "nginx:latest",
								},
							},
						},
					},
				},
			}
			err := k8sClient.Create(ctx, deployment)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}

			By("Creating DeploymentCustomOperator resource")
			resource := &dpscalerv1alpha1.DeploymentCustomOperator{
				ObjectMeta: metav1.ObjectMeta{
					Name:      resourceName,
					Namespace: namespace,
				},
				Spec: dpscalerv1alpha1.DeploymentCustomOperatorSpec{
					Targets: []dpscalerv1alpha1.DeploymentScaleTarget{
						{
							Name:      deploymentName,
							Namespace: namespace,
							Replicas:  5,
						},
					},
					Schedule: dpscalerv1alpha1.ScalingSchedule{
						StartHour:       0,
						EndHour:         23,
						DefaultReplicas: 1,
					},
				},
			}
			err = k8sClient.Create(ctx, resource)
			if err != nil && !apierrors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		})

		AfterEach(func() {
			By("Cleaning up DeploymentCustomOperator resource")
			resource := &dpscalerv1alpha1.DeploymentCustomOperator{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			if err == nil {
				Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
			}

			By("Cleaning up test deployment")
			deployment := &appsv1.Deployment{}
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      deploymentName,
				Namespace: namespace,
			}, deployment)
			if err == nil {
				Expect(k8sClient.Delete(ctx, deployment)).To(Succeed())
			}
		})

		It("should scale deployment successfully", func() {
			By("Reconciling the custom resource")
			controllerReconciler := &DeploymentCustomOperatorReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking deployment replicas are scaled to 5")
			Eventually(func() int32 {
				deployment := &appsv1.Deployment{}
				err := k8sClient.Get(ctx, types.NamespacedName{
					Name:      deploymentName,
					Namespace: namespace,
				}, deployment)
				if err != nil {
					return 0
				}
				if deployment.Spec.Replicas == nil {
					return 0
				}
				return *deployment.Spec.Replicas
			}, 10*time.Second, 1*time.Second).Should(Equal(int32(5)))

			By("Checking status is set to active")
			Eventually(func() bool {
				resource := &dpscalerv1alpha1.DeploymentCustomOperator{}
				err := k8sClient.Get(ctx, typeNamespacedName, resource)
				if err != nil {
					return false
				}
				return resource.Status.Active
			}, 10*time.Second, 1*time.Second).Should(BeTrue())
		})
	})
})
