// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"fmt"
	"net/http"

	marathon "github.com/gambol99/go-marathon"
)

// MakeHealthHandler creates a handler to check health
func MakeHealthHandler(client marathon.Marathon) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, "OK") }
}
