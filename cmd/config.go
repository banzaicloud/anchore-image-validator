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
	"os"
	"strings"

	"github.com/banzaicloud/anchore-image-validator/internal/app"
	"github.com/banzaicloud/anchore-image-validator/internal/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config contains configs
type Config struct {
	// App configuration
	App app.Config
	// Log configuration
	Log log.Config
}

func configure(v *viper.Viper, p *pflag.FlagSet) {
	p.Init("imagecheck", pflag.ExitOnError)
	pflag.Usage = func() {
		_, _ = fmt.Fprintln(os.Stderr, "Usage of imagecheck:")

		pflag.PrintDefaults()
	}
	_ = v.BindPFlags(p)
	// Log configuration
	v.SetDefault("log.format", "json")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.noColor", true)
	// App configuration
	v.SetDefault("app.port", 8443)
	v.SetDefault("app.certfile", "/webhook/certs/tls.crt")
	v.SetDefault("app.keyfile", "/webhook/certs/tls.key")
	v.SetDefault("app.cacheTTL", "30m")

	v.AllowEmptyEnv(true)
	v.SetEnvPrefix("imagecheck")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath(os.Getenv("CONFIG_DIR"))
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
