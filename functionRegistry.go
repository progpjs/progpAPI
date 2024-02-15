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

import "C"
import (
	"embed"
	"errors"
	"log"
	"path"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

var gFunctionRegistry *FunctionRegistry

//region ExposedFunction

type RegisteredFunction struct {
	IsAsync            bool
	Group              string
	JsFunctionName     string
	GoFunctionName     string
	GoFunctionFullName string
	GoFunctionRef      any
	GoFunctionInfos    ParsedGoFunction
}

//endregion

//region FunctionRegistry

type FunctionRegistry struct {
	modules           map[string]bool
	functionsArray    []*RegisteredFunction
	functionsMap      map[string]*RegisteredFunction
	extraGoNamespaces []string
	jsModulesTSX      map[string]EmbeddedFile
	useDynamicMode    bool
}

func GetFunctionRegistry() *FunctionRegistry {
	if gFunctionRegistry == nil {
		gFunctionRegistry = &FunctionRegistry{
			modules:      make(map[string]bool),
			jsModulesTSX: make(map[string]EmbeddedFile),
			functionsMap: make(map[string]*RegisteredFunction),
		}
	}

	return gFunctionRegistry
}

func (m *FunctionRegistry) GetNamespaces() map[string]bool {
	nsMap := make(map[string]bool)

	for _, e := range m.extraGoNamespaces {
		nsMap[e] = true
	}

	for _, fct := range m.functionsMap {
		fctNS := fct.GoFunctionInfos.CallParamNamespaces

		if fctNS != nil {
			for _, e := range fctNS {
				nsMap[e] = true
			}
		}
	}

	return nsMap
}

func (m *FunctionRegistry) UseGoNamespace(goNamespace string) *FunctionModule {
	nsParts := strings.Split(goNamespace, "/")
	modName := nsParts[len(nsParts)-1]
	m.extraGoNamespaces = append(m.extraGoNamespaces, goNamespace)

	return &FunctionModule{
		functionRegistry: m,
		namespacePath:    goNamespace,
		moduleName:       modName,
	}
}

func (m *FunctionRegistry) declareModuleAsNotEmpty(modName string) {
	m.modules[modName] = true
}

func (m *FunctionRegistry) addFunction(isAsync bool, group string, jsFunctionName string, goFunctionName string, goFunctionRef any) {
	fct := &RegisteredFunction{
		IsAsync:            isAsync,
		Group:              group,
		JsFunctionName:     jsFunctionName,
		GoFunctionName:     goFunctionName,
		GoFunctionFullName: group + "/" + goFunctionName,
		GoFunctionRef:      goFunctionRef,
	}

	parsed, err := ParseGoFunction(fct)

	if err != nil {
		log.Fatal("Error when parsing function " + goFunctionName + ".\nMessage: " + err.Error())
	}

	fct.GoFunctionInfos = parsed
	m.functionsArray = append(m.functionsArray, fct)
	m.functionsMap[fct.JsFunctionName] = fct
}

func (m *FunctionRegistry) GetAllFunctions(sortList bool) []*RegisteredFunction {
	if sortList {
		sort.Slice(m.functionsArray, func(i, j int) bool {
			return m.functionsArray[i].GoFunctionFullName < m.functionsArray[j].GoFunctionFullName
		})
	}

	return m.functionsArray
}

/*func (m *JsFunctionRegistry) installNodeModules(targetDir string) {
	for key, entry := range m.jsModulesTSX {
		asBytes, err := entry.Read()

		if err == nil {
			// Convert for filesystem where the path-separator is \
			parts := strings.Split(key, "/")
			key := path.Join(parts...)

			progpScripts.SafeWriteFile(path.Join(targetDir, key)+".tsx", asBytes)
		}
	}
}*/

func (m *FunctionRegistry) declareNodeModule(embedded embed.FS, innerPath string, modName string) {
	innerPath = path.Join("embedded", innerPath, modName) + ".tsx"
	modName = "@progp/" + modName

	m.jsModulesTSX[modName] = EmbeddedFile{FS: embedded, InnerPath: innerPath}
}

func (m *FunctionRegistry) GetEmbeddedModulesTSX() map[string]EmbeddedFile {
	return m.jsModulesTSX
}

func (m *FunctionRegistry) GetRefToFunction(jsFunctionName string) *RegisteredFunction {
	v, ok := m.functionsMap[jsFunctionName]
	if ok {
		return v
	}
	return nil
}

func (m *FunctionRegistry) EnableDynamicMode(enabled bool) {
	m.useDynamicMode = enabled
}

func (m *FunctionRegistry) IsUsingDynamicMode() bool {
	return m.useDynamicMode
}

//endregion

//region FunctionGroup

type FunctionGroup struct {
	jsGroupName string
	goModule    *FunctionModule
}

func (m *FunctionGroup) AddFunction(javascriptName string, goFunctionName string, goFunctionRef any) {
	m.goModule.addFunction(false, m.jsGroupName, javascriptName, goFunctionName, goFunctionRef)
}

func (m *FunctionGroup) AddAsyncFunction(jsName string, goFunctionName string, jsFunction any) {
	m.goModule.addFunction(true, m.jsGroupName, jsName, goFunctionName, jsFunction)
}

//endregion

//region FunctionModule

type FunctionModule struct {
	namespacePath    string
	moduleName       string
	isModuleInjected bool
	functionRegistry *FunctionRegistry
	modGroup         FunctionGroup
}

func (m *FunctionModule) ModName() string {
	return m.moduleName
}

// AddFunction add a function to a javascript group
// which name is the name of the go namespace last part.
func (m *FunctionModule) AddFunction(javascriptName string, goFunctionName string, goFunctionRef any) {
	m.addFunction(false, m.moduleName, javascriptName, goFunctionName, goFunctionRef)
}

// AddAsyncFunction add an async function to a javascript group
// which name is the name of the go namespace last part.
func (m *FunctionModule) AddAsyncFunction(jsName string, goFunctionName string, jsFunction any) {
	m.addFunction(true, m.moduleName, jsName, goFunctionName, jsFunction)
}

func (m *FunctionModule) GetFunctionRegistry() *FunctionRegistry {
	return m.functionRegistry
}

// UseCustomGroup allows using another javascript group than
// the default group for this go namespace.
func (m *FunctionModule) UseCustomGroup(jsGroupName string) *FunctionGroup {
	return &FunctionGroup{jsGroupName: jsGroupName, goModule: m}
}

// UseGroupGlobal allows using the global group
// where functionsArray directly accessible to javascript scripts without importing them.
func (m *FunctionModule) UseGroupGlobal() *FunctionGroup { return m.UseCustomGroup("global") }

func (m *FunctionModule) addFunction(isAsync bool, groupName string, javascriptName string, goFunctionName string, goFunctionRef any) {
	if groupName == "" {
		groupName = "global"
	}

	goFunctionName = m.moduleName + "." + goFunctionName
	endsWithAsync := strings.HasSuffix(goFunctionName, "Async")

	if isAsync {
		if !endsWithAsync {
			log.Fatalf("Go function `%s` is asynchrone and MUST ends with 'Async'", goFunctionName)
		}
	} else {
		if endsWithAsync {
			log.Fatalf("Go function `%s` is NOT asynchrone and MUST NOT ends with 'Async'", goFunctionName)
		}
	}

	if !m.isModuleInjected {
		m.isModuleInjected = true
		m.functionRegistry.declareModuleAsNotEmpty(m.moduleName)
	}

	m.functionRegistry.addFunction(isAsync, groupName, javascriptName, goFunctionName, goFunctionRef)
}

func (m *FunctionModule) DeclareNodeModule(embedded embed.FS, embeddedDirPath string, modName string) {
	m.functionRegistry.declareNodeModule(embedded, embeddedDirPath, modName)
}

//endregion

//region EmbeddedFile

type EmbeddedFile struct {
	FS        embed.FS
	InnerPath string
}

func (m *EmbeddedFile) Read() ([]byte, error) {
	return m.FS.ReadFile(m.InnerPath)
}

//endregion

//region GoFunctionInfos

type ParsedGoFunction struct {
	GeneratorUniqName string
	GoFunctionName    string

	ParamTypes          []string
	ParamTypeRefs       []reflect.Type
	CallParamNamespaces []string

	ReturnType        string
	ReturnErrorOffset int

	JsFunctionName string
	JsGroupName    string
}

func (m *ParsedGoFunction) GetJsFunctionName() string {
	return m.JsFunctionName
}

func (m *ParsedGoFunction) GetGoFunctionName() string {
	return m.GoFunctionName
}

//endregion

var gNextGoFunctionId = 0

func ParseGoFunction(fct *RegisteredFunction) (ParsedGoFunction, error) {
	reflectFct := reflect.TypeOf(fct.GoFunctionRef)

	sgn := reflectFct.String()
	res := ParsedGoFunction{ReturnErrorOffset: -1}

	if !strings.HasPrefix(sgn, "func(") {
		return res, errors.New("expect a function as second param")
	}

	// region Extract information from this function signature

	sgn = sgn[5:]

	idx := strings.Index(sgn, ")")
	returnInfos := strings.TrimSpace(sgn[idx+1:])
	sgn = sgn[0:idx]

	if sgn == "" {
		res.ParamTypes = []string{}
	} else {
		res.ParamTypes = strings.Split(sgn, ", ")
	}

	if returnInfos != "" {
		if returnInfos[0] == '(' {
			returnInfos = returnInfos[1 : len(returnInfos)-1]
		}

		returnTypes := strings.Split(returnInfos, ", ")

		size := len(returnTypes)

		if size >= 1 {
			if size == 1 {
				if returnTypes[0] == "error" {
					res.ReturnErrorOffset = 0
					returnTypes = nil
				} else {
					res.ReturnType = returnTypes[0]
				}
			} else if size == 2 {
				if returnTypes[0] == "error" {
					res.ReturnErrorOffset = 0
					returnTypes = returnTypes[1:]

					if returnTypes[0] == "error" {
						log.Fatalf("Function %s can return (error, error)", fct.GoFunctionFullName)
					}
				} else if returnTypes[1] == "error" {
					res.ReturnErrorOffset = 1
					returnTypes = returnTypes[0:1]
				} else {
					log.Fatalf("Function %s has more than 1 return type", fct.GoFunctionFullName)
				}

				res.ReturnType = returnTypes[0]
			} else {
				log.Fatalf("Function %s has more than 1 return type", fct.GoFunctionFullName)
			}
		}
	}

	//endregion

	//region Extract namespace for his call arguments

	paramsCount := reflectFct.NumIn()
	res.ParamTypeRefs = make([]reflect.Type, paramsCount)

	for i := 0; i < paramsCount; i++ {
		param := reflectFct.In(i)
		res.ParamTypeRefs[i] = param

		for {
			kind := param.Kind()

			if (kind == reflect.Pointer) || (kind == reflect.Array) || (kind == reflect.Slice) || (kind == reflect.Map) {
				param = param.Elem()
			} else {
				break
			}
		}

		pkgPath := param.PkgPath()

		if pkgPath != "" {
			if !slices.Contains(res.CallParamNamespaces, pkgPath) {
				res.CallParamNamespaces = append(res.CallParamNamespaces, pkgPath)
			}
		}
	}

	//endregion

	gNextGoFunctionId++
	res.GeneratorUniqName = strconv.Itoa(gNextGoFunctionId)

	return res, nil
}

// StringBuffer allows the code generator to known that we want this bytes
// to be send as if it was a string. Allows to avoid the cost of converting
// []byte to string before calling javascript.
type StringBuffer []byte
