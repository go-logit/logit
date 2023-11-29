// Copyright 2023 FishGoddess. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logit

import (
	"context"
)

type contextKey struct{}

// NewContextWithKey wraps context with logger of key and returns a new context.
func NewContextWithKey(ctx context.Context, key interface{}, logger *Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

// FromContextWithKey gets logger from context of key and returns the default logger if missed.
func FromContextWithKey(ctx context.Context, key interface{}) *Logger {
	if logger, ok := ctx.Value(key).(*Logger); ok {
		return logger
	}

	// TODO 返回 default logger
	return New()
}

// NewContext wraps context with logger and returns a new context.
func NewContext(ctx context.Context, logger *Logger) context.Context {
	return NewContextWithKey(ctx, contextKey{}, logger)
}

// FromContext gets logger from context and returns the default logger if missed.
func FromContext(ctx context.Context) *Logger {
	return FromContextWithKey(ctx, contextKey{})
}
