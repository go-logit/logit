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
	"io"
	"os"
	"path/filepath"

	"github.com/FishGoddess/logit/defaults"
	"github.com/FishGoddess/logit/handler"
	"github.com/FishGoddess/logit/rotate"
	"github.com/FishGoddess/logit/writer"
)

// Option sets some fields to config.
type Option func(conf *config)

func (o Option) applyTo(conf *config) {
	o(conf)
}

// WithDebugLevel sets debug level to config.
func WithDebugLevel() Option {
	return func(conf *config) {
		conf.level = LevelDebug
	}
}

// WithInfoLevel sets info level to config.
func WithInfoLevel() Option {
	return func(conf *config) {
		conf.level = LevelInfo
	}
}

// WithWarnLevel sets warn level to config.
func WithWarnLevel() Option {
	return func(conf *config) {
		conf.level = LevelWarn
	}
}

// WithErrorLevel sets error level to config.
func WithErrorLevel() Option {
	return func(conf *config) {
		conf.level = LevelError
	}
}

// WithPrintLevel sets print level to config.
func WithPrintLevel() Option {
	return func(conf *config) {
		conf.level = LevelPrint
	}
}

// WithOffLevel sets off level to config.
// All logs will be discarded.
func WithOffLevel() Option {
	return func(conf *config) {
		conf.level = LevelOff
	}
}

// WithWriter sets writer to config.
// The writer is for writing logs.
func WithWriter(w io.Writer) Option {
	newWriter := func() (io.Writer, error) {
		return w, nil
	}

	return func(conf *config) {
		conf.newWriter = newWriter
	}
}

// WithStdout sets os.Stdout to config.
// All logs will be written to stdout.
func WithStdout() Option {
	newWriter := func() (io.Writer, error) {
		return os.Stdout, nil
	}

	return func(conf *config) {
		conf.newWriter = newWriter
	}
}

// WithStderr sets os.Stderr to config.
// All logs will be written to stderr.
func WithStderr() Option {
	newWriter := func() (io.Writer, error) {
		return os.Stderr, nil
	}

	return func(conf *config) {
		conf.newWriter = newWriter
	}
}

// WithFile sets file to config.
// All logs will be written to a file in path.
// It will create all directories in path if not existed.
// The permission bits can be specified by defaults package.
// See defaults.FileDirMode and defaults.FileMode.
// If you want to customize the way open dir or file, see defaults.OpenFileDir and defaults.OpenFile.
func WithFile(path string) Option {
	newWriter := func() (io.Writer, error) {
		dir := filepath.Dir(path)
		if err := defaults.OpenFileDir(dir, defaults.FileDirMode); err != nil {
			return nil, err
		}

		return defaults.OpenFile(path, defaults.FileMode)
	}

	return func(conf *config) {
		conf.newWriter = newWriter
	}
}

// WithRotateFile sets rotate file to config.
// All logs will be written to a rotate file in path.
// A rotate file is useful in production, see rotate.File.
// The permission bits can be specified by defaults package.
// See defaults.FileDirMode and defaults.FileMode.
// Use rotate.Option to customize your rotate file.
func WithRotateFile(path string, opts ...rotate.Option) Option {
	newWriter := func() (io.Writer, error) {
		return rotate.New(path, opts...)
	}

	return func(conf *config) {
		conf.newWriter = newWriter
	}
}

// WithBuffer sets a buffer writer to config.
// You should specify a buffer size in bytes.
// The remained data in buffer may discard if you kill the process without syncing or closing the logger.
func WithBuffer(bufferSize uint64) Option {
	wrapWriter := func(w io.Writer) writer.Writer {
		return writer.Buffer(w, bufferSize)
	}

	return func(conf *config) {
		conf.wrapWriter = wrapWriter
	}
}

// WithBatch sets a batch writer to config.
// You should specify a batch size in count.
// The remained logs in batch may discard if you kill the process without syncing or closing the logger.
func WithBatch(batchSize uint64) Option {
	wrapWriter := func(w io.Writer) writer.Writer {
		return writer.Batch(w, batchSize)
	}

	return func(conf *config) {
		conf.wrapWriter = wrapWriter
	}
}

// WithHandler sets handler to config.
// It's a function returning a slog.Handler instance.
// You can return slog's handlers or your customizing handlers by this function.
func WithHandler(newHandler NewHandlerFunc) Option {
	return func(conf *config) {
		conf.newHandler = newHandler
	}
}

// WithTextHandler sets text handler to config.
func WithTextHandler() Option {
	return func(conf *config) {
		conf.newHandler = handler.NewTextHandler
	}
}

// WithJsonHandler sets json handler to config.
func WithJsonHandler() Option {
	return func(conf *config) {
		conf.newHandler = handler.NewJsonHandler
	}
}

// WithSource sets withSource=true to config.
// All logs will carry their caller information like file and line.
func WithSource() Option {
	return func(conf *config) {
		conf.withSource = true
	}
}

// WithPID sets withPID=true to config.
// All logs will carry the process id.
func WithPID() Option {
	return func(conf *config) {
		conf.withPID = true
	}
}

// ProductionOptions returns some options that we think they are useful in production.
// We recommend you to use them, so we provide this convenient way to create such a logger.
// The logger using these options will use rotate file and batch writer.
// Also, source and pid will be carried in all logs.
func ProductionOptions() []Option {
	opts := []Option{
		WithInfoLevel(), WithSource(), WithPID(),
		WithBatch(16), WithRotateFile("./logit.log"),
	}

	return opts
}
