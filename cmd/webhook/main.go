/*
Copyright (c) 2023 SAP SE

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
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/spf13/pflag"

	"k8s.io/klog/v2"

	"github.com/sap/clustersecret-operator/internal/admission"
)

var (
	bindAddress string
	tlsEnabled  bool
	tlsKeyFile  string
	tlsCertFile string
)

func main() {
	// initialize stderr logger
	errlog := log.New(os.Stderr, "", 0)

	// parse flags
	pflag.StringVar(&bindAddress, "bind_address", ":1080", "Bind address")
	pflag.BoolVar(&tlsEnabled, "tls_enabled", false, "Enable TlS")
	pflag.StringVar(&tlsKeyFile, "tls_key_file", "", "Path to TLS key")
	pflag.StringVar(&tlsCertFile, "tls_cert_file", "", "Path to TLS certificate")
	klog.InitFlags(nil)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.CommandLine.SortFlags = false
	pflag.Parse()

	if tlsEnabled {
		if tlsKeyFile == "" {
			errlog.Fatal("flag --tls_key_file is required")
		}
		if tlsCertFile == "" {
			errlog.Fatal("flag --tls_cert_file is required")
		}
	}

	// start webhooks
	klog.Infof("starting webhook on %s (TLS enabled: %v)", bindAddress, tlsEnabled)
	admissionHandler := admission.NewHandler()
	http.HandleFunc("/healthz", func(http.ResponseWriter, *http.Request) {})
	http.HandleFunc("/validation", admissionHandler.Validate)
	http.HandleFunc("/mutation", admissionHandler.Mutate)
	if tlsEnabled {
		klog.Fatalf("error running http listener: %s", http.ListenAndServeTLS(bindAddress, tlsCertFile, tlsKeyFile, nil))
	} else {
		klog.Fatalf("error running http listener: %s", http.ListenAndServe(bindAddress, nil))
	}
}
