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

package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgraph-io/ristretto"
	admissionv1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"logur.dev/logur"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const imageValidate = "/imagecheck"

// nolint: gochecknoglobals
var (
	runtimeScheme = runtime.NewScheme()
	codecs        = serializer.NewCodecFactory(runtimeScheme)
	deserializer  = codecs.UniversalDeserializer()

//	defaulter     = runtime.ObjectDefaulter(runtimeScheme)
)

// NewApp creates new application
func NewApp(logger logur.Logger, client client.Client, cache *ristretto.Cache, cacheTTL time.Duration) http.Handler {
	mux := http.NewServeMux()
	mux.Handle(imageValidate, newHTTPHandler(logger, client, cache, cacheTTL))
	logger.Info("newApp", map[string]interface{}{"app": imageValidate})

	return mux
}

// HTTPController collects the greeting use cases and exposes them as HTTP handlers.
type HTTPController struct {
	Logger   logur.Logger
	Client   client.Client
	Cache    *ristretto.Cache
	CacheTTL time.Duration
}

// NewHTTPHandler returns a new HTTP handler for the greeter.
func newHTTPHandler(logger logur.Logger,
	client client.Client,
	cache *ristretto.Cache,
	cacheTTL time.Duration) http.Handler {
	mux := http.NewServeMux()
	controller := NewHTTPController(logger, client, cache, cacheTTL)
	mux.HandleFunc(imageValidate, controller.webhookCTRL)
	logger.Info("newHTTPHandler", map[string]interface{}{"handler": imageValidate})

	return mux
}

// NewHTTPController returns a new HTTPController instance.
func NewHTTPController(logger logur.Logger,
	client client.Client,
	cache *ristretto.Cache,
	cacheTTL time.Duration) *HTTPController {
	return &HTTPController{
		Logger:   logger,
		Client:   client,
		Cache:    cache,
		CacheTTL: cacheTTL,
	}
}

func (a *HTTPController) webhookCTRL(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "reading request body failed", http.StatusInternalServerError)

		return
	}

	a.Logger.Debug("request body", map[string]interface{}{"body": body})

	if len(body) == 0 {
		http.Error(w, "empty body", http.StatusBadRequest)

		return
	}

	var admissionResponse *admissionv1beta1.AdmissionResponse

	ar := admissionv1beta1.AdmissionReview{}

	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		a.Logger.Error("Can't decode body")

		admissionResponse = &admissionv1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	} else {
		// nolint: forbidigo
		fmt.Println(r.URL.Path)
		if r.URL.Path == imageValidate {
			admissionResponse = validate(&ar, a.Logger, a.Client, a.Cache, a.CacheTTL)
		}
	}

	if admissionResponse != nil {
		ar.Response = admissionResponse
		if ar.Request != nil {
			ar.Response.UID = ar.Request.UID
		}
	}

	resp, err := json.Marshal(ar)
	if err != nil {
		a.Logger.Error("Can't encode response")
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}

	a.Logger.Info("Ready to write response ...")

	if _, err := w.Write(resp); err != nil {
		a.Logger.Error("Can't write response")
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}
