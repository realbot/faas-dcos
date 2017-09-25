// Copyright (c) Alex Ellis 2017, Alberto Quario 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import "strings"

func Function2ID(functionName string) string {
	return "/faas/functions/" + functionName
}

func ID2Function(ID string) string {
	return strings.TrimPrefix(ID, "/faas/functions/")
}

func Function2Endpoint(functionName string) string {
	return "functions-" + functionName
}
