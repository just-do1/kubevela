package apply

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"

	oamstd "github.com/oam-dev/kubevela/apis/standard.oam.dev/v1alpha1"
	oamutil "github.com/oam-dev/kubevela/pkg/oam/util"
)

var _ = Describe("Test apply", func() {
	var (
		int32_3   = int32(3)
		int32_5   = int32(5)
		ctx       = context.Background()
		deploy    *appsv1.Deployment
		podspec   *oamstd.PodSpecWorkload
		deployKey = types.NamespacedName{
			Name:      "testdeploy",
			Namespace: ns,
		}
		podspecKey = types.NamespacedName{
			Name:      "testpodspec",
			Namespace: ns,
		}
	)

	BeforeEach(func() {
		deploy = basicTestDeployment()
		podspec = basicTestPodSpecWorkload()

		Expect(k8sApplicator.Apply(ctx, deploy)).Should(Succeed())
		Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
	})

	AfterEach(func() {
		Expect(rawClient.Delete(ctx, deploy)).Should(SatisfyAny(Succeed(), &oamutil.NotFoundMatcher{}))
		Expect(rawClient.Delete(ctx, podspec)).Should(SatisfyAny(Succeed(), &oamutil.NotFoundMatcher{}))
	})

	Context("Test apply resources", func() {
		It("Test apply core resources", func() {
			deploy = basicTestDeployment()
			By("Set normal & array field")
			deploy.Spec.Replicas = &int32_3
			deploy.Spec.Template.Spec.Volumes = []corev1.Volume{{Name: "test"}}
			Expect(k8sApplicator.Apply(ctx, deploy)).Should(Succeed())
			resultDeploy := basicTestDeployment()
			Expect(rawClient.Get(ctx, deployKey, resultDeploy)).Should(Succeed())
			Expect(*resultDeploy.Spec.Replicas).Should(Equal(int32_3))
			Expect(len(resultDeploy.Spec.Template.Spec.Volumes)).Should(Equal(1))

			deploy = basicTestDeployment()
			By("Override normal & array field")
			deploy.Spec.Replicas = &int32_5
			deploy.Spec.Template.Spec.Volumes = []corev1.Volume{{Name: "test"}, {Name: "test2"}}
			Expect(k8sApplicator.Apply(ctx, deploy)).Should(Succeed())
			resultDeploy = basicTestDeployment()
			Expect(rawClient.Get(ctx, deployKey, resultDeploy)).Should(Succeed())
			Expect(*resultDeploy.Spec.Replicas).Should(Equal(int32_5))
			Expect(len(resultDeploy.Spec.Template.Spec.Volumes)).Should(Equal(2))

			deploy = basicTestDeployment()
			By("Unset normal & array field")
			deploy.Spec.Replicas = nil
			deploy.Spec.Template.Spec.Volumes = nil
			Expect(k8sApplicator.Apply(ctx, deploy)).Should(Succeed())
			resultDeploy = basicTestDeployment()
			Expect(rawClient.Get(ctx, deployKey, resultDeploy)).Should(Succeed())
			By("Unsetted fields shoulde be removed or set to default value")
			Expect(*resultDeploy.Spec.Replicas).Should(Equal(int32(1)))
			Expect(len(resultDeploy.Spec.Template.Spec.Volumes)).Should(Equal(0))
		})

		It("Test apply custom resources", func() {
			// use standard.oam.dev/v1alpha1 podSpecWorkload as sample CR
			podspec = basicTestPodSpecWorkload()
			By("Set normal & array field")
			podspec.Spec.Replicas = &int32_3
			podspec.Spec.PodSpec.Volumes = []corev1.Volume{{Name: "test"}}
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec := basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			Expect(*resultPodSpec.Spec.Replicas).Should(Equal(int32_3))
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(1))

			podspec = basicTestPodSpecWorkload()
			By("Override normal & array field")
			podspec.Spec.Replicas = &int32_5
			podspec.Spec.PodSpec.Volumes = []corev1.Volume{{Name: "test"}, {Name: "test2"}}
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec = basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			Expect(*resultPodSpec.Spec.Replicas).Should(Equal(int32_5))
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(2))

			podspec = basicTestPodSpecWorkload()
			By("Unset normal & array field")
			podspec.Spec.Replicas = nil
			podspec.Spec.PodSpec.Volumes = nil
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec = basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			By("Unsetted fields shoulde be removed")
			Expect(resultPodSpec.Spec.Replicas).Should(BeNil())
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(0))
		})

		It("Test multiple appliers", func() {
			deploy = basicTestDeployment()
			originalDeploy := deploy.DeepCopy()
			Expect(k8sApplicator.Apply(ctx, deploy)).Should(Succeed())

			modifiedDeploy := &appsv1.Deployment{}
			modifiedDeploy.SetGroupVersionKind(deploy.GroupVersionKind())
			Expect(rawClient.Get(ctx, deployKey, modifiedDeploy)).Should(Succeed())
			By("Other applier changed the deployment")
			modifiedDeploy.Spec.MinReadySeconds = 10
			modifiedDeploy.Spec.ProgressDeadlineSeconds = pointer.Int32Ptr(20)
			modifiedDeploy.Spec.Template.Spec.Volumes = []corev1.Volume{{Name: "test"}}
			Expect(rawClient.Update(ctx, modifiedDeploy)).Should(Succeed())

			By("Original applier apply again")
			Expect(k8sApplicator.Apply(ctx, originalDeploy)).Should(Succeed())
			resultDeploy := basicTestDeployment()
			Expect(rawClient.Get(ctx, deployKey, resultDeploy)).Should(Succeed())

			By("Check the changes from other applier are not effected")
			Expect(resultDeploy.Spec.MinReadySeconds).Should(Equal(int32(10)))
			Expect(*resultDeploy.Spec.ProgressDeadlineSeconds).Should(Equal(int32(20)))
			Expect(len(resultDeploy.Spec.Template.Spec.Volumes)).Should(Equal(1))
		})

		It("Test apply resources for rollout", func() {
			// use standard.oam.dev/v1alpha1 podSpecWorkload as sample CR
			podspec = basicTestPodSpecWorkload()
			By("Set normal & array field")
			podspec.Spec.Replicas = &int32_3
			podspec.Spec.PodSpec.Volumes = []corev1.Volume{{Name: "test"}}
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec := basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			Expect(*resultPodSpec.Spec.Replicas).Should(Equal(int32_3))
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(1))

			podspec = basicTestPodSpecWorkload()
			By("Override normal & array field")
			podspec.Spec.Replicas = &int32_5
			podspec.Spec.PodSpec.Volumes = []corev1.Volume{{Name: "test"}, {Name: "test2"}}
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec = basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			Expect(*resultPodSpec.Spec.Replicas).Should(Equal(int32_5))
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(2))

			podspec = basicTestPodSpecWorkload()
			By("Unset normal & array field")
			podspec.Spec.Replicas = nil
			podspec.Spec.PodSpec.Volumes = nil
			Expect(k8sApplicator.Apply(ctx, podspec)).Should(Succeed())
			resultPodSpec = basicTestPodSpecWorkload()
			Expect(rawClient.Get(ctx, podspecKey, resultPodSpec)).Should(Succeed())
			By("Unsetted fields shoulde be removed")
			Expect(resultPodSpec.Spec.Replicas).Should(BeNil())
			Expect(len(resultPodSpec.Spec.PodSpec.Volumes)).Should(Equal(0))
		})
	})
})

func basicTestDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testdeploy",
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			// Replicas: x  // normal field with default value
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{ // array field
						{
							Name:  "nginx",
							Image: "nginx:1.9.4", // normal field without default value
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "nginx"}},
			},
		},
	}
}

func basicTestPodSpecWorkload() *oamstd.PodSpecWorkload {
	return &oamstd.PodSpecWorkload{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodSpecWorkload",
			APIVersion: "standard.oam.dev/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testpodspec",
			Namespace: ns,
		},
		Spec: oamstd.PodSpecWorkloadSpec{
			// Replicas: x (normal field)
			PodSpec: corev1.PodSpec{
				Containers: []corev1.Container{ // array field
					{Name: "nginx",
						Image: "nginx:1.9.4"},
				},
			},
		},
	}
}
