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
	"sync"
)

var gScriptEngines = make(map[string]ScriptEngine)
var gScriptEnginesMutex sync.RWMutex
var gScriptEngineBuilder = make(map[string]ScriptEngineBuilder)

type ScriptEngineBuilder = func() ScriptEngine

type ScriptEngine interface {
	// Start the engine, which is call one all is initialized in the Go side.
	Start()

	// GetEngineLanguage allows to know if it' a "javascript" engine or a "python" engine.
	GetEngineLanguage() string

	// GetEngineName the name of the underlying engine. For exemple "progpv8".
	GetEngineName() string

	WaitDebuggerReady()

	// GetInternalEngineVersion the version of the engine used internally.
	// For exemple if it's using Google V8, then return the v8 engine version.
	GetInternalEngineVersion() string

	// ExecuteScript execute a script by giving his content.
	ExecuteScript(scriptContent string, compiledFilePath string) *ScriptErrorMessage

	// Shutdown stop the engine. He can't be used anymore after that.
	// It mainly occurs after a fatal error or at script ends.
	Shutdown()

	// DisarmError remove the current error and allows continuing execution.
	// The error params allows to avoid case where a new error occurs since.
	DisarmError(error *ScriptErrorMessage)
}

type ScriptFunction interface {
	CallWithArrayBuffer(buffer []byte)
	CallWithString(value string)
	CallWithStringBuffer(value []byte)
	CallWithDouble(value float64)
	CallWithUndefined()
	CallWithError(err error)
	CallWithBool(value bool)
	CallWithResource(value *SharedResource)
	KeepAlive()
}

type ScriptCallback func(error *ScriptErrorMessage)

func ConfigRegisterScriptEngineBuilder(engineName string, builder ScriptEngineBuilder) {
	gScriptEngineBuilder[engineName] = builder
}

func GetScriptEngine(engineName string) ScriptEngine {
	gScriptEnginesMutex.RLock()
	engine := gScriptEngines[engineName]
	gScriptEnginesMutex.RUnlock()

	if engine != nil {
		return engine
	}

	gScriptEnginesMutex.Lock()
	defer gScriptEnginesMutex.Unlock()

	builder := gScriptEngineBuilder[engineName]
	if builder == nil {
		return nil
	}

	engine = builder()
	gScriptEngines[engineName] = engine

	return engine
}

func ForEachScriptEngine(f func(engine ScriptEngine)) {
	for _, e := range gScriptEngines {
		f(e)
	}
}
