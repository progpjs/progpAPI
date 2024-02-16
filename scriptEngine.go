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

	// Shutdown stop the engine. He can't be used anymore after that.
	// It mainly occurs after a fatal error or at script ends.
	Shutdown()

	// CreateNewScriptContext creates a new context which can be used
	// to execute a new script context from the others scripts.
	//
	CreateNewScriptContext(securityGroup string) ScriptContext

	// SetRuntimeErrorHandler allows to set a function which will manage runtime error.
	// The handler runtime true if the error is handler or false
	// to use the default behavior, which consist of printing the error and exit.
	//
	SetRuntimeErrorHandler(handler RuntimeErrorHandlerF)

	// SetScriptTerminatedHandler allows to add a function triggered when the script has finished his execution, with or without error.
	// It's call when all asynchronous function are executed, end before the end of the background tasks.
	// (mainly because this tasks can continue to executed without needing the javascript VM anymore)
	//
	SetScriptTerminatedHandler(handler ScriptTerminatedHandlerF)

	SetAllowedFunctionsChecker(handler CheckAllowedFunctionsF)
}

type RuntimeErrorHandlerF func(ctx ScriptContext, err *ScriptErrorMessage) bool
type ScriptTerminatedHandlerF func(ctx ScriptContext, scriptPath string, err *ScriptErrorMessage) *ScriptErrorMessage
type ScriptCallbackF func(error *ScriptErrorMessage)
type ScriptFileExecutorF func(ctx ScriptContext, scriptPath string) *ScriptErrorMessage
type CheckAllowedFunctionsF func(securityGroup string, functionGroup string, functionName string) bool

var gScriptFileExecutor ScriptFileExecutorF

type ScriptFunction interface {
	CallWithUndefined()

	CallWithError(err error)

	// KeepAlive allows to avoid destroying the function after the first call.
	// It must be used when you keep a reference to a function.
	//
	KeepAlive()

	// EnabledResourcesAutoDisposing allows the engine to automatically dispose the resources created while
	// calling this function. Without that you must call progpDispose on each disposable resources.
	// Here no, when activating this flag the engine release all the resource one the function call ends.
	// This includes all the async functions launched from this function and not only the main body of the function.
	//
	// If you are interested in this functionality, you can use the javascript function progpAutoDispose(() => { ... })
	//
	EnabledResourcesAutoDisposing(currentResourceContainer *SharedResourceContainer)

	// "2" means second argument.
	// It's used for callback for first argument is error message.

	CallWithArrayBuffer2(buffer []byte)
	CallWithString2(value string)
	CallWithStringBuffer2(value []byte)
	CallWithDouble2(value float64)

	CallWithBool2(value bool)
	CallWithResource2(value *SharedResource)
}

type ScriptContext interface {
	GetScriptEngine() ScriptEngine

	// GetSecurityGroup returns a group name which allows knowing the category of this context.
	// It's mainly used to allows / don't allow access to some functions groups.
	// For exemple you can use security group "unsafe" then the script will no be able to access to Go functions.
	//
	GetSecurityGroup() string

	// ExecuteScript executes a script inside this context.
	// It must be used once and don't allow executing more than one script.
	ExecuteScript(scriptContent string, compiledFilePath string, sourceScriptPath string) *ScriptErrorMessage

	// ExecuteScriptFile is like ExecuteScript but allows using a file (which can be typescript).
	ExecuteScriptFile(ctx ScriptContext, scriptPath string) *ScriptErrorMessage

	// TryDispose destroy the context and free his resources.
	// It's do nothing if this context can't be disposed, for
	// exemple if the engine only support one context.
	//
	TryDispose() bool

	// DisarmError remove the current error and allows continuing execution.
	// The error params allows to avoid case where a new error occurs since.
	DisarmError(error *ScriptErrorMessage)

	// IncreaseRefCount increase the ref counter for the context.
	// This avoid that the script exit, which is required the system is
	// keeping reference on some javascript functions.
	IncreaseRefCount()

	// DecreaseRefCount decrease the ref counter for the context.
	DecreaseRefCount()
}

func GetScriptFileExecutor() ScriptFileExecutorF {
	return gScriptFileExecutor
}

func SetScriptFileExecutor(executor ScriptFileExecutorF) {
	gScriptFileExecutor = executor
}

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
