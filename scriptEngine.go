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
	CreateNewScriptContext(securityGroup string, mustDebug bool) JsContext

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

type OnScriptCompilationErrorF func(scriptPath string, err error) bool
type RuntimeErrorHandlerF func(ctx JsContext, err *JsErrorMessage) bool
type ScriptTerminatedHandlerF func(ctx JsContext, scriptPath string, err *JsErrorMessage) *JsErrorMessage
type ScriptCallbackF func(error *JsErrorMessage)
type ScriptFileExecutorF func(ctx JsContext, scriptPath string) *JsErrorMessage
type ScriptFileCompilerF func(scriptPath string) (string, string, error)
type CheckAllowedFunctionsF func(securityGroup string, functionGroup string, functionName string) bool
type ListenProgpSignalF func(ctx JsContext, signal string) error

var gScriptFileExecutor ScriptFileExecutorF
var getScriptFileCompiler ScriptFileCompilerF

type JsFunction interface {
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

	CallWithDouble1(value float64)
	CallWithDouble2(value float64)

	CallWithBool2(value bool)

	CallWithResource1(value *SharedResource)
	CallWithResource2(value *SharedResource)
}

type JsContext interface {
	GetScriptEngine() ScriptEngine

	// GetSecurityGroup returns a group name which allows knowing the category of this context.
	// It's mainly used to allows / don't allow access to some functions groups.
	// For exemple you can use security group "unsafe" then the script will no be able to access to Go functions.
	//
	GetSecurityGroup() string

	// ExecuteScript executes a script inside this context.
	// It must be used once and don't allow executing more than one script.
	ExecuteScript(scriptContent string, compiledFilePath string, sourceScriptPath string) *JsErrorMessage

	// ExecuteScriptFile is like ExecuteScript but allows using a file (which can be typescript).
	ExecuteScriptFile(scriptPath string) *JsErrorMessage

	// ExecuteChildScriptFile execute a script from the inside of another script.
	ExecuteChildScriptFile(scriptPath string) error

	// TryDispose destroy the context and free his resources.
	// It's do nothing if this context can't be disposed, for
	// exemple if the engine only support one context.
	//
	TryDispose() bool

	// DisarmError remove the current error and allows continuing execution.
	// The error params allows to avoid case where a new error occurs since.
	DisarmError(error *JsErrorMessage)

	// IncreaseRefCount increase the ref counter for the context.
	// This avoids that the script exit, which is required the system is
	// keeping reference on some javascript functions.
	IncreaseRefCount()

	// DecreaseRefCount decrease the ref counter for the context.
	DecreaseRefCount()
}

func GetScriptFileExecutor() ScriptFileExecutorF {
	return gScriptFileExecutor
}

func GetScriptFileCompiler() ScriptFileCompilerF {
	return getScriptFileCompiler
}

func SetScriptFileExecutor(executor ScriptFileExecutorF) {
	gScriptFileExecutor = executor
}

func SetScriptFileCompiler(compiler ScriptFileCompilerF) {
	getScriptFileCompiler = compiler
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
