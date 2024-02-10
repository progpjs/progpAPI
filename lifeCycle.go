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

type ScriptErrorMessage struct {
	isTranslated bool
	isLogged     bool

	ScriptEngine ScriptEngine

	Error      string
	ErrorLevel int

	StartColumn int
	EndColumn   int

	StartPosition int
	EndPosition   int

	SourceMapUrl string
	ScriptPath   string

	StackTraceFrameCount int
	StackTraceFrames     []StackTraceFrame
}

type StackTraceFrame struct {
	Line     int
	Column   int
	Function string
	Source   string
}

func (m *ScriptErrorMessage) Translate() {
	/*if !m.isTranslated && gScriptTransformer != nil {
		gScriptTransformer.TranslateErrorMessage(m)
	}

	m.isTranslated = true*/
}

func (m *ScriptErrorMessage) StackTrace() string {
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

func (m *ScriptErrorMessage) Print() {
	m.Translate()
	fmt.Printf("Javascript Error - %s\n", m.Error)
	print(m.StackTrace())
}

// DisarmError allows to continue after an un-catch error.
func (m *ScriptErrorMessage) DisarmError() {
	if m.ScriptEngine != nil {
		m.ScriptEngine.DisarmError(m)
	}
}

//endregion

func LogScriptError(error *ScriptErrorMessage) {
	if (error == nil) || (error.isLogged) {
		return
	}

	error.isLogged = true
	error.Print()
}

func OnUnCatchScriptError(error *ScriptErrorMessage) {
	if gErrorTranslator != nil {
		gErrorTranslator(error)
	}

	LogScriptError(error)
	error.ScriptEngine.Shutdown()
}

func SetErrorTranslator(handler ErrorTranslatorF) {
	gErrorTranslator = handler
}

type ErrorTranslatorF func(error *ScriptErrorMessage)

var gErrorTranslator ErrorTranslatorF

//endregion

//region ScriptExecResult

type ScriptExecResult struct {
	ScriptError *ScriptErrorMessage
	GoError     error
}

type V8ScriptCallback func(result ScriptExecResult)

func (m *ScriptExecResult) HasError() bool {
	return m.ScriptError != nil || m.GoError != nil
}

func (m *ScriptExecResult) PrintError() bool {
	if m.ScriptError != nil {
		m.ScriptError.Print()
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
	if gBackgroundTasksCount != 0 {
		gBackgroundTasksCount--
	}
	gBackgroundTasksCountMutex.Unlock()

	if gBackgroundTasksCount == 0 {
		close(gBackgroundTasksWaitChannel)
	}
}

// ForceExitingVM allows stopping the process without doing an os.exit.
// It's a requirement if profiling the memory, since without that, the log file isn't correctly flushed.
func ForceExitingVM() {
	gBackgroundTasksCount = 0
	close(gBackgroundTasksWaitChannel)

	ForEachScriptEngine(func(e ScriptEngine) {
		e.Shutdown()
	})
}

func EndOfAllBackgroundTasks() {
	<-gBackgroundTasksWaitChannel
}

//endregion

//region Executing script

func ExecuteScriptContent(scriptContent, scriptOrigin string, scriptEngine ScriptEngine) {
	err := scriptEngine.ExecuteScript(scriptContent, scriptOrigin)

	if err == nil {
		// The script exit but the VM must continue to execute
		// if a background task is executing.
		//
		EndOfAllBackgroundTasks()
	}
}

//endregion
