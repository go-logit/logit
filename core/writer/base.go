// Copyright 2021 Ye Zi Jie. All Rights Reserved.
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
//
// Author: FishGoddess
// Email: fishgoddess@qq.com
// Created at 2021/07/01 23:40:06

package writer

import "io"

const (
	KB = 1024      // KB is the unit KB in size. 1 KB = 1024 Bytes.
	MB = 1024 * KB // MB is the unit MB in size. 1 MB = 1024*1024 Bytes.

	bufferSize = 16 * KB
)

type Flusher interface {
	Flush() (n int, err error)
}

type Writer interface {
	Flusher
	io.WriteCloser
}

type writerWrapper struct {
	writer io.Writer
}

func (ww *writerWrapper) Flush() (n int, err error) {
	if flusher, ok := ww.writer.(Flusher); ok {
		return flusher.Flush()
	}
	return 0, nil
}

func (ww *writerWrapper) Write(p []byte) (n int, err error) {
	return ww.writer.Write(p)
}

func (ww *writerWrapper) Close() error {
	if closer, ok := ww.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

func Wrapped(writer io.Writer) Writer {
	if w, ok := writer.(Writer); ok {
		return w
	}
	return &writerWrapper{writer: writer}
}

func Buffered(writer io.Writer) Writer {
	if w, ok := writer.(Writer); ok {
		return w
	}
	return newBufferedWriter(writer, bufferSize)
}
