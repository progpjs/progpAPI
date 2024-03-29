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
	"sync"
)

//region Error management

//region ScriptErrorMessage

type JsErrorMessage struct {
	scriptContext JsContext

	isTranslated bool
	isLogged     bool
	isPrinted    bool

	Error      string
	ErrorLevel int

	StartColumn int
	EndColumn   int

	StartPosition int
	EndPosition   int

	SourceMapUrl string
	SourceMap    string
	ScriptPath   string

	StackTraceFrameCount int
	StackTraceFrames     []StackTraceFrame
}

func NewScriptErrorMessage(ctx JsContext) *JsErrorMessage {
	return &JsErrorMessage{scriptContext: ctx}
}

func (m *JsErrorMessage) GetScriptContext() JsContext {
	return m.scriptContext
}

type StackTraceFrame struct {
	Line     int
	Column   int
	Function string
	Source   string
}

func (m *JsErrorMessage) Translate() {
	if m.isTranslated {
		return
	}
	m.isTranslated = true

	if gErrorTranslator != nil {
		gErrorTranslator(m)
	}
}

func (m *JsErrorMessage) StackTrace() string {
	m.Translate()
	res := ""

	for _, frame := range m.StackTraceFrames {
		if frame.Function == "" {
			res += fmt.Sprintf("at %s:%d:%d\n", frame.Source, frame.Line, frame.Column)
		} else {
			res += fmt.Sprintf("at %s:%d:%d: %s\n", frame.Source, frame.Line, frame.Column, frame.Function)
		}
	}

	return res
}

func (m *JsErrorMessage) Print(forcePrinting bool) {
	if m.isPrinted && !forcePrinting {
		return
	}
	m.isPrinted = true

	m.Translate()
	fmt.Printf("Javascript Error - %s\n", m.Error)
	print(m.StackTrace())
}

func (m *JsErrorMessage) LogError() {
	if (m == nil) || m.isLogged {
		return
	}
	m.isLogged = true
	m.Print(false)
}

// DisarmError allows to continue after an un-catch error.
func (m *JsErrorMessage) DisarmError(ctx JsContext) {
	ctx.DisarmError(m)
}

//endregion

func OnUnCatchScriptError(error *JsErrorMessage) {
	error.LogError()
	error.scriptContext.GetScriptEngine().Shutdown()
}

func SetErrorTranslator(handler ErrorTranslatorF) {
	gErrorTranslator = handler
}

type ErrorTranslatorF func(error *JsErrorMessage)

var gErrorTranslator ErrorTranslatorF

//endregion

//region ScriptExecResult

type ScriptExecResult struct {
	ScriptError *JsErrorMessage
	GoError     error
}

type V8ScriptCallback func(result ScriptExecResult)

func (m *ScriptExecResult) HasError() bool {
	return m.ScriptError != nil || m.GoError != nil
}

func (m *ScriptExecResult) PrintError() bool {
	if m.ScriptError != nil {
		m.ScriptError.Print(false)
		return true
	} else if m.GoError != nil {
		println("GO ERROR - " + m.GoError.Error())
		return true
	}

	return false
}

//endregion

//region Background tasks

var gBackgroundTasksCount = 0
var gBackgroundTasksCountMutex sync.Mutex
var gBackgroundTasksWaitChannel = make(chan bool)

func DeclareBackgroundTaskStarted() {
	gBackgroundTasksCountMutex.Lock()
	defer gBackgroundTasksCountMutex.Unlock()
	gBackgroundTasksCount++
}

func DeclareBackgroundTaskEnded() {
	gBackgroundTasksCountMutex.Lock()
	defer gBackgroundTasksCountMutex.Unlock()

	if gBackgroundTasksCount != 0 {
		gBackgroundTasksCount--
	}

	if gBackgroundTasksCount == 0 {
		close(gBackgroundTasksWaitChannel)
		gBackgroundTasksWaitChannel = nil
	}
}

// ForceExitingVM allows stopping the process without doing an os.exit.
// It's a requirement if profiling the memory, since without that, the log file isn't correctly flushed.
func ForceExitingVM() {
	gBackgroundTasksCountMutex.Lock()
	defer gBackgroundTasksCountMutex.Unlock()

	gBackgroundTasksCount = 0

	if gBackgroundTasksWaitChannel != nil {
		close(gBackgroundTasksWaitChannel)
	}

	ForEachScriptEngine(func(e ScriptEngine) {
		e.Shutdown()
	})
}

// WaitTasksEnd wait until all background tasks are finished.
// It's used in order to know if the application can exit.
func WaitTasksEnd() {
	if gBackgroundTasksWaitChannel != nil {
		<-gBackgroundTasksWaitChannel
	}
}

//endregion
