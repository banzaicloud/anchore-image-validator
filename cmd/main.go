// Copyright © 2019 Banzai Cloud.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/banzaicloud/anchore-image-validator/internal/app"
	"github.com/banzaicloud/anchore-image-validator/internal/log"
	"github.com/banzaicloud/anchore-image-validator/pkg/apis/security/v1alpha1"
	"github.com/dgraph-io/ristretto"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes/scheme"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	crconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
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
		// nolint: forbidigo
		fmt.Printf("%s version %s (%s) built on %s\n", "anchore-image-validator", version, commitHash, buildDate)

		os.Exit(0)
	}

	err := viper.ReadInConfig()

	var configFileNotFound viper.ConfigFileNotFoundError

	if !errors.As(err, &configFileNotFound) {
		emperror.Panic(errors.Wrap(err, "failed to read configuration"))
	}

	var config Config
	err = viper.Unmarshal(&config)

	if err != nil {
		emperror.Panic(errors.Wrap(err, "failed to unmarshal configuration"))
	}

	if viper.GetBool("dump-config") {
		// nolint: forbidigo
		fmt.Printf("%+v\n", config)

		os.Exit(0)
	}

	// Create logger (first thing after configuration loading)
	logger := log.NewLogger(config.Log)

	// Provide some basic context to all log lines
	logger = log.WithFields(logger, map[string]interface{}{"service": "imagecheck"})

	k8sCfg := crconfig.GetConfigOrDie()

	logger.Debug("kubernetes config", map[string]interface{}{
		"k8sHost": k8sCfg.Host,
	})

	v1alpha1.AddToScheme(scheme.Scheme)

	client, err := crclient.New(k8sCfg, crclient.Options{})
	if err != nil {
		logger.Error("get clisntset failed", map[string]interface{}{
			"k8sHost": k8sCfg.Host,
		})
	}

	logger.Info("starting the webhook.", map[string]interface{}{
		"port":     ":" + config.App.Port,
		"certfile": config.App.CertFile,
		"keyfile":  config.App.KeyFile,
		"cacheTTL": config.App.CacheTTL,
	})

	// create cache object
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,     // number of keys to track frequency of (10k).
		MaxCost:     100000000, // maximum cost of cache (100MB).
		BufferItems: 64,        // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}

	var cacheTTL time.Duration

	cacheTTL, err = time.ParseDuration(config.App.CacheTTL)

	if err != nil {
		logger.Error("parse cacheTTL failed", map[string]interface{}{
			"cacheTTL": config.App.CacheTTL,
		},
		)

		cacheTTL = time.Duration(30) * time.Minute
	}

	if viper.GetBool("dev-http") {
		http.ListenAndServe(":"+config.App.Port, app.NewApp(logger, client, cache, cacheTTL))
	} else {
		http.ListenAndServeTLS(":"+config.App.Port,
			config.App.CertFile,
			config.App.KeyFile,
			app.NewApp(logger, client, cache, cacheTTL))
	}
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}

	return value
}
