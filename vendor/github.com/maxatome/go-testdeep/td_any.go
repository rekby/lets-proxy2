// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
)

type tdAny struct {
	tdList
}

var _ TestDeep = &tdAny{}

// Any operator compares data against several expected values. During
// a match, at least one of them has to match to succeed.
//
// TypeBehind method can return a non-nil reflect.Type if all items
// known non-interface types are equal, or if only interface types
// are found (mostly issued from Isa()) and they are equal.
func Any(expectedValues ...interface{}) TestDeep {
	return &tdAny{
		tdList: newList(expectedValues...),
	}
}

func (a *tdAny) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	for _, item := range a.items {
		if deepValueEqualOK(got, item) {
			return nil
		}
	}

	if ctx.BooleanError {
		return ctxerr.BooleanError
	}
	return ctx.CollectError(&ctxerr.Error{
		Message:  "comparing with Any",
		Got:      got,
		Expected: a,
	})
}

func (a *tdAny) TypeBehind() reflect.Type {
	return a.uniqTypeBehind()
}
