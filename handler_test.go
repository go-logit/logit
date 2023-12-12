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
	"fmt"
	"io"
	"log/slog"
	"testing"
)

// go test -v -cover -count=1 -test.cpu=1 -run=^TestPickHandler$
func TestPickHandler(t *testing.T) {
	newHandler := func(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
		return nil
	}

	newHandlers[t.Name()] = newHandler
	got, err := PickHandler(t.Name())
	if err != nil {
		t.Fatal(err)
	}

	if fmt.Sprintf("%p", got) != fmt.Sprintf("%p", newHandler) {
		t.Fatalf("got %p is wrong", got)
	}
}

// go test -v -cover -count=1 -test.cpu=1 -run=^TestRegisterHandler$
func TestRegisterHandler(t *testing.T) {
	if err := RegisterHandler("text", nil); err == nil {
		t.Fatal("register an existed handler func should be failed")
	}

	handler := "new"
	newHandlerFunc := func(w io.Writer, opts *slog.HandlerOptions) slog.Handler { return nil }

	if err := RegisterHandler(handler, newHandlerFunc); err != nil {
		t.Fatal(err)
	}

	newHandler, ok := newHandlers[handler]
	if !ok {
		t.Fatalf("handler %s not found", handler)
	}

	if fmt.Sprintf("%p", newHandler) != fmt.Sprintf("%p", newHandlerFunc) {
		t.Fatal("newHandler registered is wrong")
	}
}
