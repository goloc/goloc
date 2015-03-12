// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"strings"
)

// Filter function type.
// If return value false, the result is removed.
type Filter func(*Result) bool

func DefaultFilter(result *Result) bool {
	if strings.HasPrefix(result.Type, "zone") {
		result.Score += 2
	} else if strings.HasPrefix(result.Type, "poi") {
		result.Score += 1
	} else if strings.HasPrefix(result.Type, "street") {
	}
	return true
}
