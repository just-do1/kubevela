package appfile

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	corev1alpha2 "github.com/oam-dev/kubevela/apis/core.oam.dev"
	"github.com/oam-dev/kubevela/apis/core.oam.dev/v1alpha2"
	"github.com/oam-dev/kubevela/pkg/oam/util"
	"github.com/oam-dev/kubevela/pkg/utils/system"
	// +kubebuilder:scaffold:imports
)

var cfg *rest.Config
var scheme *runtime.Scheme
var k8sClient client.Client
var testEnv *envtest.Environment
var definitionDir string
var wd v1alpha2.WorkloadDefinition
var addonNamespace = "test-addon"

func TestAppFile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Cli Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.New(zap.UseDevMode(true), zap.WriteTo(GinkgoWriter)))
	ctx := context.Background()
	By("bootstrapping test environment")
	useExistCluster := false
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:  []string{filepath.Join("..", "..", "charts", "vela-core", "crds")},
		UseExistingCluster: &useExistCluster,
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())
	scheme = runtime.NewScheme()
	Expect(corev1alpha2.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(clientgoscheme.AddToScheme(scheme)).NotTo(HaveOccurred())
	Expect(v1beta1.AddToScheme(scheme)).NotTo(HaveOccurred())
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	definitionDir, err = system.GetCapabilityDir()
	Expect(err).Should(BeNil())
	Expect(os.MkdirAll(definitionDir, 0755)).Should(BeNil())

	Expect(k8sClient.Create(ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: addonNamespace}})).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	workloadData, err := ioutil.ReadFile("testdata/workloadDef.yaml")
	Expect(err).Should(BeNil())

	Expect(yaml.Unmarshal(workloadData, &wd)).Should(BeNil())

	wd.Namespace = addonNamespace
	logf.Log.Info("Creating workload definition", "data", wd)
	Expect(k8sClient.Create(ctx, &wd)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	def, err := ioutil.ReadFile("testdata/terraform-aliyun-oss-workloadDefinition.yaml")
	Expect(err).Should(BeNil())
	var terraformDefinition v1alpha2.WorkloadDefinition
	Expect(yaml.Unmarshal(def, &terraformDefinition)).Should(BeNil())
	terraformDefinition.Namespace = addonNamespace
	logf.Log.Info("Creating workload definition", "data", terraformDefinition)
	Expect(k8sClient.Create(ctx, &terraformDefinition)).Should(SatisfyAny(BeNil(), &util.AlreadyExistMatcher{}))

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	_ = k8sClient.Delete(context.Background(), &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: addonNamespace}})
	_ = k8sClient.Delete(context.Background(), &wd)
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
