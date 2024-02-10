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

type CustomType struct {
	typeName string
}

func (m *CustomType) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return "&" + paramName
}

func (m *CustomType) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_V8OBJECT_TOSTRING"
}

func (m *CustomType) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "*C.s_progp_goStringOut"
}

func (m *CustomType) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnArrayBuffer"
}

func (m *CustomType) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "v8::Local<v8::Value> v8Res;\n    V8VALUE_FROM_GOCUSTOM(v8Res, res, resWrapper.size);\n    callInfo.GetReturnValue().Set(v8Res);"
}

func (m *CustomType) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	// Here we avoid string copy, so we haven't to call free.
	return ""
}

func (m *CustomType) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	typeName := m.typeName
	ctx.AddNamespace("encoding/json")

	res := "    b" + paramName + " := C.GoBytes(unsafe.Pointer(" + paramName + ".p), " + paramName + ".n)\n"
	res += "    var v" + paramName + " " + typeName + "\n"
	res += "    if err := json.Unmarshal(b" + paramName + ", &v" + paramName + "); err !=nil {\n"
	res += "        res.errorMessage = C.CString(err.Error())\n"
	res += "        return\n"
	res += "    }"

	return res, "v" + paramName
}

func (m *CustomType) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	ctx.AddNamespace("encoding/json")

	return `    asBytes, err := json.Marshal(goRes)

	if err != nil {
		res.errorMessage = C.CString(err.Error())
	} else {
		res.value = unsafe.Pointer(&asBytes[0])
		res.size = C.int(len(asBytes))
	}`
}
