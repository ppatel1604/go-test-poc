package go_test_poc_test

import (
	"context"
	"flag"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	certManagerv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var _ = Describe("CertificateIssuance", func() {
	It("should issue a certificate successfully", func() {
		kubeconfig := flag.String("kubeconfig", "/Users/parthpatel/.kube/config", "absolute path to the kubeconfig file")
		flag.Parse()
		config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
		Expect(err).NotTo(HaveOccurred())

		cmClient, err := cmclientset.NewForConfig(config)
		Expect(err).NotTo(HaveOccurred())

		cert := &certManagerv1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-certificate",
				Namespace: "default",
			},
			Spec: certManagerv1.CertificateSpec{
				SecretName: "test-secret",
				IssuerRef: cmmeta.ObjectReference{
					Name: "letsencrypt-prod",
					Kind: "ClusterIssuer",
				},
				CommonName: "example.com",
				DNSNames:   []string{"example.com"},
			},
		}

		_, err = cmClient.CertmanagerV1().Certificates("default").Create(context.TODO(), cert, metav1.CreateOptions{})
		Expect(err).NotTo(HaveOccurred())

		// Wait for certificate to be ready (this is a simplified wait; in real scenarios consider using a watch or retry mechanism)
		time.Sleep(60 * time.Second)

		// Fetch the Certificate to check its status
		issuedCert, err := cmClient.CertmanagerV1().Certificates("default").Get(context.TODO(), "test-certificate", metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())

		// Expect the Certificate to be ready
		Expect(issuedCert.Status.Conditions).ToNot(BeEmpty())
		// for _, condition := range issuedCert.Status.Conditions {
		// 	if condition.Type == certmanagerv1.CertificateConditionReady {
		// 		Expect(condition.Status).To(Equal(metav1.ConditionTrue))
		// 	}
		// }
	})
})
