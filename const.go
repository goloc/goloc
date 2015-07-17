// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import (
	"time"
)

const (
	defaultLimit            = int(10)
	defaultTolerance        = float32(.5)
	defaultWorkLimit        = int(200)
	defaultMaxWaitAcquire   = 1 * time.Second
	defaultMaxWaitTraitment = 3 * time.Second

	// min/max value for int /unit
	maxUint = uint(1<<32 - 1)
	minUint = uint(0)
	maxInt  = int(1<<31 - 1)
	minInt  = int(-1 << 31)
)
