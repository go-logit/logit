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
	"log/slog"
	"testing"
)

// go test -v -cover -run=^TestConfigHandlerOptions$
func TestConfigHandlerOptions(t *testing.T) {
	replaceAttr := func(groups []string, attr slog.Attr) slog.Attr { return attr }

	conf := &config{
		level:       levelWarn,
		withSource:  true,
		replaceAttr: replaceAttr,
	}

	opts := conf.handlerOptions()

	if opts.Level != conf.level {
		t.Errorf("opts.Level %v != conf.level %v", opts.Level, conf.level)
	}

	if opts.AddSource != conf.withSource {
		t.Errorf("opts.AddSource %v != conf.withSource %v", opts.AddSource, conf.withSource)
	}

	if fmt.Sprintf("%p", opts.ReplaceAttr) != fmt.Sprintf("%p", conf.replaceAttr) {
		t.Errorf("opts.ReplaceAttr %p != conf.replaceAttr %p", opts.ReplaceAttr, conf.replaceAttr)
	}
}
