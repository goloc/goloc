// Copyright 2015 Mathieu MAST. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package goloc

import ()

type Result struct {
	Score  int
	Search string
	Id     string
	Type   string
	Name   string
	Number string
	Point
}

func CompareScoreResult(r1, r2 interface{}) int {
	id1 := r1.(*Result).Id
	id2 := r2.(*Result).Id
	if id1 == id2 {
		return 0
	}
	dif := r2.(*Result).Score - r1.(*Result).Score
	if dif != 0 {
		return dif
	}
	n1 := r1.(*Result).Name
	n2 := r2.(*Result).Name
	if n1 > n2 {
		return 1
	} else if n1 < n2 {
		return -1
	}
	if id1 > id2 {
		return 1
	} else if id1 < id2 {
		return -1
	}
	return 0
}
