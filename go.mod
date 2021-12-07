module github.com/quanxiang-cloud/overseer

go 1.16

require (
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/shipwright-io/build v0.6.0
	github.com/tektoncd/pipeline v0.30.0
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	knative.dev/pkg v0.0.0-20211101212339-96c0204a70dc
	knative.dev/serving v0.27.1
	sigs.k8s.io/controller-runtime v0.10.1
)

replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
