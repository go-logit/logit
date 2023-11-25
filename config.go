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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/FishGoddess/logit/defaults"
	"github.com/FishGoddess/logit/io/file"
	"github.com/FishGoddess/logit/io/size"
	"github.com/FishGoddess/logit/io/writer"
)

type WriterConfig struct {
	// Target is where the writer writes logs.
	// Values: "stdout", "stderr", "file".
	Target string `json:"target" yaml:"target" toml:"target" bson:"target"`

	// Mode is how the writer writes logs.
	// Values: "direct", "buffer", "batch".
	// Direct means writer writes logs without any buffer or batch, which is one log one writing operation.
	// Buffer means writer will keep logs in a buffer and write these logs once until the buffer is full.
	// Batch means writer will keep logs in a batch and write these logs once until the count of logs in this batch >= batch size.
	// Both of buffer and batch have better performance in writing logs.
	// However, they will lose some logs if the program crashed before syncing.
	Mode string `json:"mode" yaml:"mode" toml:"mode" bson:"mode"`

	// BufferSize is the size of a buffer.
	// You can use common words like "512B" or "4KB".
	// Only available when mode is "buffer".
	BufferSize string `json:"buffer_size" yaml:"buffer_size" toml:"buffer_size" bson:"buffer_size"`

	// BatchSize is the size of a batch.
	// Only available when mode is "batch".
	BatchSize uint `json:"batch_size" yaml:"batch_size" toml:"batch_size" bson:"batch_size"`

	// AutoSync is the frequency of syncing.
	// An empty string means syncing is manual.
	// You can use common words like "5m" or "60s".
	// See time.Duration and time.ParseDuration.
	// Only available when mode is "buffer" or "batch".
	AutoSync string `json:"auto_sync" yaml:"auto_sync" toml:"auto_sync" bson:"auto_sync"`
}

type FileConfig struct {
	// Path is the path (or prefix) of log file.
	Path string `json:"path" yaml:"path" toml:"path" bson:"path"`

	// Mode is the permission bits of log files.
	Mode os.FileMode `json:"mode" yaml:"mode" toml:"mode" bson:"mode"`

	// DirMode is the permission bits of directory storing log files.
	DirMode os.FileMode `json:"dir_mode" yaml:"dir_mode" toml:"dir_mode" bson:"dir_mode"`

	// Rotate is log file should split and backup when satisfy some conditions.
	// It's useful in production so we recommend you to set it to true.
	Rotate bool `json:"rotate" yaml:"rotate" toml:"rotate" bson:"rotate"`

	// MaxSize is the max size of a log file.
	// If size of data in one output operation is bigger than this value, then file will rotate before writing,
	// which means file and its backups may be bigger than this value in size.
	// You can use common words like "100MB" or "1GB".
	// Only available when rotate is true.
	MaxSize string `json:"max_size" yaml:"max_size" toml:"max_size" bson:"max_size"`

	// MaxAge is the time that backups will live.
	// All backups reach max age will be removed automatically.
	// You can use common words like "7d" or "24h".
	// See time.Duration and time.ParseDuration.
	// Only available when rotate is true.
	MaxAge string `json:"max_age" yaml:"max_age" toml:"max_age" bson:"max_age"`

	// MaxBackups is the max count of file backups.
	// Only available when rotate is true.
	MaxBackups int `json:"max_backups" yaml:"max_backups" toml:"max_backups" bson:"max_backups"`
}

type Config struct {
	// Level is the level of logger.
	// Values: debug, info, warn, error, print, off.
	Level string `json:"level" yaml:"level" toml:"level" bson:"level"`

	// Handler is how the handler handles the logs.
	// Values: "text", "json", "slog.text", "slog.json".
	// These handlers with "slog" prefix are from slog package of Go.
	// We recommend you to use our faster handlers, and feel free if you want to use slog's handlers.
	// Also, you can register your handlers to logit, see RegisterHandler.
	Handler string `json:"handler" yaml:"handler" toml:"handler" bson:"handler"`

	// Writer is the config of writer.
	Writer WriterConfig `json:"writer" yaml:"writer" toml:"writer" bson:"writer"`

	// File is the config of file.
	// Only available when the target of writer is "file".
	File FileConfig `json:"file" yaml:"file" toml:"file" bson:"file"`

	// WithSource adds source to logs if true.
	WithSource bool `json:"with_source" yaml:"with_source" toml:"with_source" bson:"with_source"`

	// WithPID adds pid to logs if true.
	WithPID bool `json:"with_pid" yaml:"with_pid" toml:"with_pid" bson:"with_pid"`

	replaceAttr func(groups []string, attr slog.Attr) slog.Attr
}

func NewDefaultConfig() *Config {
	conf := &Config{
		Level:   "debug",
		Handler: "text",
		Writer: WriterConfig{
			Target:     "stdout",
			Mode:       "direct",
			BufferSize: "16KB",
			BatchSize:  64,
			AutoSync:   "30s",
		},
		File: FileConfig{
			Path:       "./logit.log",
			Rotate:     false,
			Mode:       0644,
			DirMode:    0755,
			MaxSize:    "200MB",
			MaxAge:     "90d",
			MaxBackups: 100,
		},
		WithSource: false,
		WithPID:    true,
	}

	return conf
}

func (c *Config) WithReplaceAttr(replaceAttr func(groups []string, attr slog.Attr) slog.Attr) *Config {
	c.replaceAttr = replaceAttr
	return c
}

func (c *Config) parseTimeDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") || strings.HasSuffix(s, "D") {
		s = strings.TrimSuffix(s, "d")
		s = strings.TrimSuffix(s, "D")

		days, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, err
		}

		return time.Duration(days) * defaults.Day, nil
	}

	return time.ParseDuration(s)
}

func (c *Config) newFile() (io.Writer, error) {
	if c.File.Rotate {
		opts := []file.Option{
			file.WithMode(c.File.Mode),
			file.WithDirMode(c.File.DirMode),
			file.WithTimeFormat("20060102150405"),
		}

		if c.File.MaxSize != "" {
			maxSize, err := size.ParseByteSize(c.File.MaxSize)
			if err != nil {
				return nil, err
			}

			opts = append(opts, file.WithMaxSize(maxSize))
		}

		if c.File.MaxAge != "" {
			maxAge, err := c.parseTimeDuration(c.File.MaxAge)
			if err != nil {
				return nil, err
			}

			opts = append(opts, file.WithMaxAge(maxAge))
		}

		if c.File.MaxBackups > 0 {
			opts = append(opts, file.WithMaxBackups(c.File.MaxBackups))
		}

		return file.New(c.File.Path)
	}

	dir := filepath.Dir(c.File.Path)
	if err := os.MkdirAll(dir, c.File.DirMode); err != nil {
		return nil, err
	}

	return defaults.OpenFile(c.File.Path, c.File.Mode)
}

func (c *Config) newWriter() (io.Writer, error) {
	var w writer.Writer

	switch c.Writer.Target {
	case "stdout":
		w = os.Stdout
	case "stderr":
		w = os.Stderr
	case "file":
		f, err := c.newFile()
		if err != nil {
			return nil, err
		}

		w = writer.Wrap(f)
	default:
		return nil, fmt.Errorf("writer target %s invalid", c.Writer.Target)
	}

	switch c.Writer.Mode {
	case "direct":
		break
	case "buffer":
		bufferSize, err := size.ParseByteSize(c.Writer.BufferSize)
		if err != nil {
			return nil, err
		}

		w = writer.Buffer(w, bufferSize)
	case "batch":
		w = writer.Batch(w, c.Writer.BatchSize)
	default:
		return nil, fmt.Errorf("writer mode %s invalid", c.Writer.Mode)
	}

	if c.Writer.AutoSync != "" {
		frequency, err := time.ParseDuration(c.Writer.AutoSync)
		if err != nil {
			return nil, err
		}

		go func() {
			for {
				time.Sleep(frequency)

				if err := w.Sync(); err != nil {
					defaults.HandleError("writer.sync", err)
				}
			}
		}()
	}

	return w, nil
}

func (c *Config) newHandlerOptions() *slog.HandlerOptions {
	opts := &slog.HandlerOptions{
		Level:       parseLevel(c.Level),
		AddSource:   c.WithSource,
		ReplaceAttr: c.replaceAttr,
	}

	return opts
}

func (c *Config) NewHandler() (slog.Handler, error) {
	w, err := c.newWriter()
	if err != nil {
		return nil, err
	}

	opts := c.newHandlerOptions()
	return newHandler(c.Handler, w, opts)
}
