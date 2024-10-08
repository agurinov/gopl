// Copyright 2023-2024 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// Merge is a utility to resolve and combine errors resulting from
// evaluation. If ok is false, execution of validation should stop (either due
// to failFast or the result is not a ValidationError).
//
//nolint:errorlint
func Merge(dst, src error, failFast bool) (ok bool, err error) {
	if src == nil {
		return true, dst
	}

	srcValErrs, ok := src.(*ValidationError)
	if !ok {
		return false, src
	}

	if dst == nil {
		return !(failFast && len(srcValErrs.Violations) > 0), src
	}

	dstValErrs, ok := dst.(*ValidationError)
	if !ok {
		// what should we do here?
		return false, dst
	}

	dstValErrs.Violations = append(dstValErrs.Violations, srcValErrs.Violations...)
	return !(failFast && len(dstValErrs.Violations) > 0), dst
}

// PrefixErrorPaths prepends the formatted prefix to the violations of a
// ValidationError.
func PrefixErrorPaths(err error, format string, args ...any) {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		PrefixFieldPaths(valErr, format, args...)
	}
}

func MarkForKey(err error) {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		for _, violation := range valErr.Violations {
			violation.ForKey = proto.Bool(true)
		}
	}
}
