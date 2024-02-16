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

package codegen

import "C"
import (
	"github.com/progpjs/progpAPI"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type ProgpV8CodeGenerator struct {
	// CurrentFunction is the function currently exported.
	// Is used by the data types.
	CurrentFunction *progpAPI.RegisteredFunction

	outputDir    string
	namespaces   map[string]bool
	functionList []*progpAPI.RegisteredFunction
	typeMap      map[string]IsTypeHandler

	cppImplInjectThis string
	goLangInjectThis  string
	fileCppImpl       string
	fileGoLang        string
}

func NewProgpV8Codegen() *ProgpV8CodeGenerator {
	var typeMap = make(map[string]IsTypeHandler)

	//region Register types handlers

	typeMap[""] = &TypeVoid{}
	typeMap["bool"] = &TypeBool{}
	typeMap["int"] = &TypeInt{}
	typeMap["float32"] = &TypeFloat32{}
	typeMap["float64"] = &TypeFloat64{}
	typeMap["string"] = &TypeString{}
	typeMap["[]uint8"] = &TypeUIntArray{}
	typeMap["unsafe.Pointer"] = &TypeUnsafePointer{}
	typeMap["progpAPI.ScriptFunction"] = &TypeV8Function{}
	typeMap["*progpAPI.SharedResource"] = &TypeSharedResource{}
	typeMap["*progpAPI.SharedResourceContainer"] = &TypeSharedResourceContainer{}
	typeMap["progpAPI.StringBuffer"] = &TypeStringBuffer{}

	//endregion

	fctRegistry := progpAPI.GetFunctionRegistry()
	namespaces := fctRegistry.GetNamespaces()

	// Function list must be sorted in order to always generate the same output.
	functionList := fctRegistry.GetAllFunctions(true)

	return &ProgpV8CodeGenerator{
		namespaces:   namespaces,
		functionList: functionList,
		typeMap:      typeMap,
	}
}

func (m *ProgpV8CodeGenerator) GenerateCode(autoUpdateDir string) {
	if autoUpdateDir == "" {
		return
	}

	if state, err := os.Stat(autoUpdateDir); err == nil {
		if !state.IsDir() {
			return
		}
	} else {
		return
	}

	m.outputDir = autoUpdateDir

	for _, f := range m.functionList {
		if err := m.glueCodeCreateBindingFunctionsFor(f); err != nil {
			log.Fatal(err)
			return
		}
	}

	m.createGroupFunctions()

	// Required for: "defer progpAPI.CatchFatalErrors()"
	m.AddNamespace("github.com/progpjs/progpAPI")

	nsList := ""
	for _, nsName := range m.getNamespaces() {
		nsList += "\n    " + progpAPI.GoStringToQuotedString2(nsName)
	}

	var template string

	//region File : generated.cpp

	template = `#ifndef PROGP_STANDALONE

#include "progpV8.h"
#include "_cgo_export.h"
#include <iostream>
#include <stdexcept>
%INJECT_HERE%

#endif // PROGP_STANDALONE
`

	template = strings.ReplaceAll(template, "%INJECT_HERE%", m.cppImplInjectThis)
	m.fileCppImpl = template

	//endregion

	//region file : generated.go

	template = `package libProgpV8
// #include <stdlib.h> // For C.free
// #include "progpAPI.h"
//
import "C"

import (%NAMESPACES%
)

%INJECT_HERE%
`
	template = strings.ReplaceAll(template, "%INJECT_HERE%", m.goLangInjectThis)
	template = strings.ReplaceAll(template, "%NAMESPACES%", nsList)
	m.fileGoLang = template

	//endregion

	hasUpdated := false

	if m.saveFileIfNotTheSame(path.Join(m.outputDir, "generated.cpp"), m.fileCppImpl) {
		hasUpdated = true
	}
	if m.saveFileIfNotTheSame(path.Join(m.outputDir, "generated.go"), m.fileGoLang) {
		hasUpdated = true
	}

	if hasUpdated {
		println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		println("!  Javascript binding code has been updated.  !")
		println("!  A restart is required.                     !")
		println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		os.Exit(1)
	}
}

func (m *ProgpV8CodeGenerator) AddNamespace(namespacePath string) {
	m.namespaces[namespacePath] = true
}

func (m *ProgpV8CodeGenerator) tryToCreateTypeHandler(typeName string) IsTypeHandler {
	return &CustomType{typeName: typeName}
}

func (m *ProgpV8CodeGenerator) getType(typeName string) IsTypeHandler {
	if res, ok := m.typeMap[typeName]; ok {
		return res
	}

	res := m.tryToCreateTypeHandler(typeName)

	if res != nil {
		m.typeMap[typeName] = res
		return res
	}

	log.Fatal("Type " + typeName + " not found by the code generator engine")
	return nil
}

func (m *ProgpV8CodeGenerator) getNamespaces() []string {
	var nsListArray []string

	for nsName := range m.namespaces {
		nsListArray = append(nsListArray, nsName)
	}

	// Required to always have the same order, without that
	// the generated code can have random content.
	//
	sort.Strings(nsListArray)

	return nsListArray
}

func (m *ProgpV8CodeGenerator) saveFileIfNotTheSame(filePath string, newContent string) bool {
	oldContentB, err := os.ReadFile(filePath)
	oldContent := string(oldContentB)

	if (err != nil) || (oldContent != newContent) {
		dirPath := path.Dir(filePath)

		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Fatal("Can't create directory " + dirPath)
		}

		err = os.WriteFile(filePath, []byte(newContent), os.ModePerm)
		if err != nil {
			log.Fatal("Can't write file " + filePath)
		}

		println("Codegen has updated file " + filePath)

		return true
	}

	return false
}

func (m *ProgpV8CodeGenerator) createGroupFunctions() {
	toInject := "\n\nvoid exposeGoFunctionsToV8(ProgpContext progpCtx, const std::string& group, v8::Local<v8::Object> v8Host) {"

	for _, f := range m.functionList {
		template := "\n    PROGP_BIND_FUNCTION(\"%FUNCTION_GROUP%\", \"%FUNCTION_NAME%\", (f_progp_v8_function)v8Function_%FUNCTION_FULL_NAME%);"
		template = strings.ReplaceAll(template, "%FUNCTION_GROUP%", f.Group)
		template = strings.ReplaceAll(template, "%FUNCTION_NAME%", f.JsFunctionName)
		template = strings.ReplaceAll(template, "%FUNCTION_FULL_NAME%", f.GoFunctionInfos.GeneratorUniqName)
		toInject += template
	}

	toInject += "\n}"

	m.cppImplInjectThis += toInject
}

func (m *ProgpV8CodeGenerator) glueCodeCreateBindingFunctionsFor(fct *progpAPI.RegisteredFunction) error {
	m.CurrentFunction = fct

	//region Data extracting

	cppAllParamsDecoding := ""
	cppCallParamsList := ""
	cppExtraBeforeCall := ""
	cppFreeResources := ""

	goParams := ""
	goAllParamsDecoding := ""
	goCallParamsList := ""

	returnTypeWrapper := m.getType(fct.GoFunctionInfos.ReturnType).ReturnTypeWrapper(m)
	returnTypeEncoder := m.getType(fct.GoFunctionInfos.ReturnType).ReturnTypeEncoder(m)

	if fct.IsAsync {
		cppExtraBeforeCall += "\n    progp_IncreaseContextRef();\n    resWrapper.isAsync = true;"
	}

	if len(fct.GoFunctionInfos.ParamTypes) != 0 {
		cppAllParamsDecoding = ""
		cppParamsCount := 0
		cppParamOffset := 0

		for offset, paramType := range fct.GoFunctionInfos.ParamTypes {
			argName := "p" + strconv.Itoa(offset)

			freeingResources := m.getType(paramType).CppArgResourcesFreeing(argName, m)
			if freeingResources != "" {
				cppFreeResources += "\n" + freeingResources
			}

			asCgoParam := m.getType(paramType).CgoFunctionParamType(m)
			if asCgoParam != "" {
				goParams += ", p" + strconv.Itoa(offset) + " " + asCgoParam
			}

			asV8ValueDecoder := m.getType(paramType).V8ToCppDecoder(m)

			if asV8ValueDecoder != "" {
				v := m.getType(paramType).CppToCgoParamCall(argName, m)

				cppCallParamsList += ", " + v
				cppAllParamsDecoding += "    " + asV8ValueDecoder + "(" + argName + ", " + strconv.Itoa(cppParamOffset) + ");\n"

				cppParamsCount++
				cppParamOffset++
			}

			cgoParamDecoding, cgoParamCall := m.getType(paramType).CgoToGoDecoding(argName, m)

			if cgoParamDecoding != "" {
				goAllParamsDecoding += "\n" + cgoParamDecoding + "\n"
			}

			if cgoParamCall != "" {
				goCallParamsList += ", " + cgoParamCall
			} else {
				goCallParamsList += ", " + argName
			}
		}

		if cppParamsCount > 0 {
			cppAllParamsDecoding = "    V8CALLARG_EXPECT_ARGCOUNT(" + strconv.Itoa(cppParamsCount) + ");\n" + cppAllParamsDecoding
		}
	}

	if cppExtraBeforeCall != "" {
		cppExtraBeforeCall = "\n" + cppExtraBeforeCall
	}

	goParams = "res *C." + returnTypeWrapper + goParams
	cppCallParamsList = "&resWrapper" + cppCallParamsList

	if goCallParamsList != "" {
		goCallParamsList = goCallParamsList[2:]
	}

	//endregion

	//region C++ functions

	template := `

void v8Function_%FUNCTION_FULL_NAME%(const v8::FunctionCallbackInfo<v8::Value> &callInfo) {
	PROGP_V8FUNCTION_BEFORE_PROGPCTX

%PARAMS_DECODING%
	%RETURN_TYPE_WRAPPER% resWrapper{};%EXTRA_BEFORE_CALL%
	resWrapper.currentEvent = progpCtx->event;
	progpCgoBinding__%FUNCTION_FULL_NAME%(%CALL_PARAMS_LIST%);
	%FREE_RESOURCES%
    if (resWrapper.errorMessage!=nullptr) {
		auto msg = std::string(resWrapper.errorMessage);
		delete(resWrapper.errorMessage);
        throw std::runtime_error(msg.c_str());
    } else if (resWrapper.constErrorMessage!= nullptr) {
		auto msg = std::string(resWrapper.errorMessage);
        throw std::runtime_error(resWrapper.errorMessage);
    }

    %RETURN_TYPE_ENCODER%
	PROGP_V8FUNCTION_AFTER
}`

	if fct.GoFunctionInfos.ReturnType != "" {
		returnTypeEncoder = "auto res = resWrapper.value;\n    " + returnTypeEncoder
	}

	template = strings.ReplaceAll(template, "%FUNCTION_FULL_NAME%", fct.GoFunctionInfos.GeneratorUniqName)
	template = strings.ReplaceAll(template, "%PARAMS_DECODING%", cppAllParamsDecoding)
	template = strings.ReplaceAll(template, "%CALL_PARAMS_LIST%", cppCallParamsList)
	template = strings.ReplaceAll(template, "%RETURN_TYPE_WRAPPER%", returnTypeWrapper)
	template = strings.ReplaceAll(template, "%RETURN_TYPE_ENCODER%", returnTypeEncoder)
	template = strings.ReplaceAll(template, "%EXTRA_BEFORE_CALL%", cppExtraBeforeCall)
	template = strings.ReplaceAll(template, "%FREE_RESOURCES%", cppFreeResources)

	m.cppImplInjectThis += template

	//endregion

	//region GoLang functions

	returnOutput := ""
	returnProcessing := ""

	if fct.GoFunctionInfos.ReturnType != "" {
		returnOutput = "goRes := "
		returnProcessing = m.getType(fct.GoFunctionInfos.ReturnType).GoValueToCgoValue(m)

		if fct.GoFunctionInfos.ReturnErrorOffset != -1 {
			if fct.GoFunctionInfos.ReturnErrorOffset == 0 {
				returnOutput = "err, goRes:= "
			} else {
				returnOutput = "goRes, err:= "
			}
		}
	} else if fct.GoFunctionInfos.ReturnErrorOffset != -1 {
		returnOutput = "err := "
	}

	if fct.GoFunctionInfos.ReturnErrorOffset != -1 {
		errorProcessing := `

	if err != nil {
		res.errorMessage = C.CString(err.Error())
		return
	}`

		if returnProcessing == "" {
			returnProcessing = errorProcessing
		} else {
			returnProcessing = errorProcessing + "\n" + returnProcessing
		}
	} else {
		returnProcessing = "\n" + returnProcessing
	}

	template = `
//export progpCgoBinding__%FUNCTION_FULL_NAME%
func progpCgoBinding__%FUNCTION_FULL_NAME%(%FUNCTION_PARAMS%) {
	defer progpAPI.CatchFatalErrors()
%PARAMS_DECODING%
	%RETURN_OUTPUT%%GO_FUNCTION_NAME%(%CALL_PARAMS_LIST%)%RETURN_PROCESSING%
}`

	template = strings.ReplaceAll(template, "%FUNCTION_FULL_NAME%", fct.GoFunctionInfos.GeneratorUniqName)
	template = strings.ReplaceAll(template, "%FUNCTION_NAME%", fct.JsFunctionName)
	template = strings.ReplaceAll(template, "%FUNCTION_PARAMS%", goParams)
	template = strings.ReplaceAll(template, "%PARAMS_DECODING%", goAllParamsDecoding)
	template = strings.ReplaceAll(template, "%CALL_PARAMS_LIST%", goCallParamsList)
	template = strings.ReplaceAll(template, "%GO_FUNCTION_NAME%", fct.GoFunctionName)

	template = strings.ReplaceAll(template, "%RETURN_OUTPUT%", returnOutput)
	template = strings.ReplaceAll(template, "%RETURN_PROCESSING%", returnProcessing)

	m.goLangInjectThis += template

	//endregion

	return nil
}
