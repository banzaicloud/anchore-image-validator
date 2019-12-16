/*
Copyright 2019 Banzai Cloud.

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

package main

import (
	"fmt"
	"net/http"
	"os"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/banzaicloud/anchore-image-validator/internal/app"
	"github.com/banzaicloud/anchore-image-validator/internal/log"
	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	clientv1alpha1 "github.com/banzaicloud/anchore-image-validator/pkg/clientset/v1alpha1"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	crconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var securityClientSet *clientv1alpha1.SecurityV1Alpha1Client

const apiServiceResource = "imagecheck"

var (
	apiServiceGroup     = os.Getenv("ANCHORE_APISERVICE_GROUP")
	apiServiceVersion   = os.Getenv("ANCHORE_APISERVICE_VERSION")
	anchoreReleaseName  = os.Getenv("ANCHORE_RELEASE_NAME")
	kubernetesNameSpace = os.Getenv("KUBERNETES_NAMESPACE")
	namespaceSelector   = getEnv("NAMESPACE_SELECTOR", "exclude")
)

// nolint: gochecknoinits
func init() {
	pflag.Bool("version", false, "Show version information")
	pflag.Bool("dump-config", false, "Dump configuration to the console (and exit)")
	pflag.Bool("dev-http", false, "Developer mode use http for local testing")
}

func main() {

	configure(viper.GetViper(), pflag.CommandLine)

	pflag.Parse()

	if viper.GetBool("version") {
		fmt.Printf("%s version %s (%s) built on %s\n", "anchore-image-validator", version, commitHash, buildDate)

		os.Exit(0)
	}

	err := viper.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)
	if !configFileNotFound {
		emperror.Panic(errors.Wrap(err, "failed to read configuration"))
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		emperror.Panic(errors.Wrap(err, "failed to unmarshal configuration"))
	}

	if viper.GetBool("dump-config") {
		fmt.Printf("%+v\n", config)

		os.Exit(0)
	}

	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = log.WithFields(logger, map[string]interface{}{"service": "imagecheck"})

	k8sCfg := crconfig.GetConfigOrDie()

	logger.Debug("kubernetes config", map[string]interface{}{
		"k8sHost": k8sCfg.Host})

	v1alpha1.AddToScheme(scheme.Scheme)
	securityClientSet, err = clientv1alpha1.SecurityConfig(k8sCfg)
	if err != nil {
		logger.Error("error")
	}

	client, err := crclient.New(k8sCfg, crclient.Options{})
	if err != nil {
		logger.Error("get clisntset failed", map[string]interface{}{
			"k8sHost": k8sCfg.Host})
	}

	logger.Info("starting the webhook.", map[string]interface{}{
		"port":     ":" + config.App.Port,
		"certfile": config.App.CertFile,
		"keyfile":  config.App.KeyFile,
	})

	if viper.GetBool("dev-http") {
		http.ListenAndServe(":"+config.App.Port, app.NewApp(logger, client))
	} else {
		http.ListenAndServeTLS(":"+config.App.Port, config.App.CertFile, config.App.KeyFile, app.NewApp(logger, client))
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
