// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"time"
)

const (
	defaultTolerance        = 1
	defaultLocLimit         = 200
	defaultMaxWaitAcquire   = 1 * time.Second
	defaultMaxWaitTraitment = 3 * time.Second
)
