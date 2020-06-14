// Copyright 2020 Ye Zi Jie. All Rights Reserved.
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
// Created at 2020/06/11 21:15:37

package files

import (
	"math/rand"
	"path/filepath"
	"strconv"
	"time"
)

// NameGenerator is the type for generating the name of every created file.
// You can customize your format of filename by implementing this function.
// The two parameters string and time.Time is useful. The string parameter is the directory
// stores all files created in this time and the time.Time parameter is the time of current moment.
type NameGenerator func(string, time.Time) string

// NextName is for code-readable.
// Directory stores all files created in this time and now is the time of current moment.
func (ng NameGenerator) NextName(directory string, now time.Time) string {
	return ng(directory, now)
}

// DefaultNameGenerator returns a name generator that creates a time-relative filename
// with given now time. Also, it uses random number to ensure this filename is available.
// The filename will be like "20200304-145246-45.log".
// Notice that directory stores all files created in this time and now is the time of current moment.
func DefaultNameGenerator() func(directory string, now time.Time) string {
	// 考虑使用原子计数器替换随机数
	// 这个 Seed 方法最好不要并发执行，但是这个方法有可能会被并发执行，这是个隐患
	// 在测试阶段就已经出现了随机数重复的情况，导致一个文件被写入多个文件的内容
	// issue: https://github.com/FishGoddess/logit/issues/7
	rand.Seed(time.Now().UnixNano())
	return func(directory string, now time.Time) string {
		name := now.Format("20060102-150405") + "-" + strconv.Itoa(rand.Intn(1000000)) + SuffixOfLogFile
		return filepath.Join(directory, name)
	}
}
