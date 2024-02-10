/*
 * (C) Copyright 2024 Johan Michel PIQUET, France (https://johanpiquet.fr/).
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package progpAPI

import (
	"fmt"
	"strings"
	"time"
)

func GoStringToQuotedString2(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	return "\"" + value + "\""
}

//region SafeGoRoutine

func SafeGoRoutine(f func()) {
	go func() {
		defer CatchFatalErrors()
		f()
	}()
}

func CatchFatalErrors() {
	if err := recover(); err != nil {
		fmt.Printf("CATCH FATAL ERROR: %s\n", err)
	}
}

//endregion

//region TaskQueue

// TaskQueue allows executing the C++ calls from only one thread.
// Without that, Go can be short on available threads which lead to a crash.
//
// This protection is only required where there is a lot of calls that can be blocked the thread.
// It's essentially when calling an event and calling a callback function.
type TaskQueue struct {
	channel  chan func()
	disposed bool
}

func NewTaskQueue() *TaskQueue {
	res := &TaskQueue{
		// The size avoid blocking.
		// Once exceeded, the system is locked and unexpected behaviors can occur.
		// It's why me make it big here.
		//
		channel: make(chan func(), 1024),
	}

	SafeGoRoutine(func() { res.start() })
	return res
}

func (m *TaskQueue) Push(f func()) {
	if m.disposed {
		return
	}

	m.channel <- f
}

func (m *TaskQueue) IsDisposed() bool {
	return m.disposed
}

func (m *TaskQueue) Dispose() {
	if m.disposed {
		return
	}

	m.disposed = true
	close(m.channel)
}

func (m *TaskQueue) start() {
	for {
		next := <-m.channel

		if m.disposed {
			break
		}

		next()
	}
}

//endregion

func PauseMs(timeInMs int) {
	duration := time.Millisecond * time.Duration(timeInMs)
	time.Sleep(duration)
}
