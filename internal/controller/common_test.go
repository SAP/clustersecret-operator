/*
SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package controller

import (
	"flag"

	"k8s.io/klog/v2"
)

func init() {
	klog.InitFlags(nil)
	flag.Set("v", "2")
}
