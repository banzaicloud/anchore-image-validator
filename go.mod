module github.com/banzaicloud/anchore-image-validator

go 1.16

require (
	emperror.dev/emperror v0.32.0
	emperror.dev/errors v0.7.0
	github.com/dgraph-io/ristretto v0.0.3
	github.com/docker/distribution v0.0.0-20200213211116-66809646d941
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver v0.0.0-20190918201827-3de75813f604 // indirect
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v0.17.1
	logur.dev/adapter/logrus v0.4.1
	logur.dev/logur v0.16.2
	sigs.k8s.io/controller-runtime v0.4.0
)

replace (
	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.1
	k8s.io/client-go => k8s.io/client-go v0.17.1
)
