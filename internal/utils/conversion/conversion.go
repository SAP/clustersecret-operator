/*
SPDX-FileCopyrightText: 2023 SAP SE or an SAP affiliate company and clustersecret-operator contributors
SPDX-License-Identifier: Apache-2.0
*/

package conversion

import (
	"strconv"
)

// convert string to int64 and panic in case of any error
func Atoi(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return int64(i)
}

// convert int64 to string
func Itoa(i int64) string {
	return strconv.FormatInt(i, 10)
}
