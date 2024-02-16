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

//region void

type TypeVoid struct {
}

func (m *TypeVoid) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeVoid) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeVoid) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnVoid"
}

func (m *TypeVoid) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().SetUndefined();"
}

func (m *TypeVoid) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeVoid) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeVoid) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", ""
}

func (m *TypeVoid) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return ""
}

//endregion

//region bool

type TypeBool struct {
}

func (m *TypeBool) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeBool) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_BOOL"
}

func (m *TypeBool) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnInt"
}

func (m *TypeBool) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_BOOL(res));"
}

func (m *TypeBool) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.int"
}

func (m *TypeBool) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeBool) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "    var " + paramName + "_asBool = true\n    if C.int(" + paramName + ") == 0 {\n        " + paramName + "_asBool = false\n    }", paramName + "_asBool"
}

func (m *TypeBool) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    if goRes {\n        res.value = C.int(1)\n    } else {\n        res.value = C.int(0)\n}"
}

//endregion

//region int

type TypeInt struct {
}

func (m *TypeInt) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeInt) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_DOUBLE"
}

func (m *TypeInt) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnLong"
}

func (m *TypeInt) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_DOUBLE(res));"
}

func (m *TypeInt) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.double"
}

func (m *TypeInt) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeInt) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "int(" + paramName + ")"
}

func (m *TypeInt) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    res.value = C.long(goRes)"
}

//endregion

//region float32

type TypeFloat32 struct {
}

func (m *TypeFloat32) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeFloat32) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_DOUBLE"
}

func (m *TypeFloat32) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnDouble"
}

func (m *TypeFloat32) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_DOUBLE(res));"
}

func (m *TypeFloat32) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.double"
}

func (m *TypeFloat32) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeFloat32) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "float32(" + paramName + ")"
}

func (m *TypeFloat32) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    res.value = C.double(goRes)"
}

//endregion

//region float64

type TypeFloat64 struct {
}

func (m *TypeFloat64) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeFloat64) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_DOUBLE"
}

func (m *TypeFloat64) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnDouble"
}

func (m *TypeFloat64) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_DOUBLE(res));"
}

func (m *TypeFloat64) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.double"
}

func (m *TypeFloat64) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeFloat64) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "float64(" + paramName + ")"
}

func (m *TypeFloat64) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    res.value = C.double(goRes)"
}

//endregion

//region string

type TypeString struct {
}

func (m *TypeString) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return "&" + paramName
}

func (m *TypeString) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_GOSTRING"
}

func (m *TypeString) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "*C.s_progp_goStringOut"
}

func (m *TypeString) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnString"
}

func (m *TypeString) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_GOSTRING(res));"
}

func (m *TypeString) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeString) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "C.GoStringN(" + paramName + ".p, " + paramName + ".n)"
}

func (m *TypeString) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	ctx.AddNamespace("unsafe")
	return "    res.value = unsafe.Pointer(&goRes)"
}

//endregion

//region progpAPI.StringBuffer

type TypeStringBuffer struct {
}

func (m *TypeStringBuffer) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return "&" + paramName
}

func (m *TypeStringBuffer) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_GOSTRING"
}

func (m *TypeStringBuffer) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "*C.s_progp_goStringOut"
}

func (m *TypeStringBuffer) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnString"
}

func (m *TypeStringBuffer) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_GOSTRING(res));"
}

func (m *TypeStringBuffer) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeStringBuffer) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	ctx.AddNamespace("unsafe")
	return "", "C.GoBytes(unsafe.Pointer(" + paramName + ".p), " + paramName + ".n)"
}

func (m *TypeStringBuffer) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	ctx.AddNamespace("unsafe")
	return "    resString:= unsafe.String(unsafe.SliceData(goRes), len(goRes))\n    res.value = unsafe.Pointer(&resString)"
}

//endregion

//region []uint / ArrayBuffer

type TypeUIntArray struct {
}

func (m *TypeUIntArray) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeUIntArray) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeUIntArray) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_ARRAYBUFFER_DATA"
}

func (m *TypeUIntArray) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnArrayBuffer"
}

func (m *TypeUIntArray) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "v8::Local<v8::ArrayBuffer> buffer = v8::ArrayBuffer::New(v8Iso, resWrapper.size);\n    memcpy(buffer->GetBackingStore()->Data(), res, resWrapper.size);\n    callInfo.GetReturnValue().Set(buffer);"
}

func (m *TypeUIntArray) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.ProgpV8BufferPtr"
}

func (m *TypeUIntArray) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "C.GoBytes(" + paramName + ".buffer, " + paramName + ".length)"
}

func (m *TypeUIntArray) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	ctx.AddNamespace("unsafe")
	return "    res.value = unsafe.Pointer(&goRes[0])\n    res.size = C.int(len(goRes))"
}

//endregion

//region unsafe.Pointer

type TypeUnsafePointer struct {
}

func (m *TypeUnsafePointer) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeUnsafePointer) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeUnsafePointer) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_EXTERNAL"
}

func (m *TypeUnsafePointer) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnExternal"
}

func (m *TypeUnsafePointer) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "v8::Local<v8::External> ext = v8::External::New(v8Iso, (void*)res);\n    callInfo.GetReturnValue().Set(ext);"
}

func (m *TypeUnsafePointer) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	ctx.AddNamespace("unsafe")
	return "unsafe.Pointer"
}

func (m *TypeUnsafePointer) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", paramName
}

func (m *TypeUnsafePointer) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    res.value = goRes"
}

//endregion

//region progpAPI.ScriptFunction

type TypeV8Function struct {
}

func (m *TypeV8Function) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeV8Function) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeV8Function) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_FUNCTION"
}

func (m *TypeV8Function) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeV8Function) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeV8Function) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.ProgpV8FunctionPtr"
}

func (m *TypeV8Function) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "newV8Function(res.isAsync, " + paramName + ", res.currentEvent)"
}

func (m *TypeV8Function) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return ""
}

//endregion

// region *progpAPI.SharedResource

type TypeSharedResource struct {
}

func (m *TypeSharedResource) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeSharedResource) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResource) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_DOUBLE"
}

func (m *TypeSharedResource) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return "ProgpFunctionReturnLong"
}

func (m *TypeSharedResource) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return "callInfo.GetReturnValue().Set(V8VALUE_FROM_DOUBLE(res));"
}

func (m *TypeSharedResource) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.double"
}

func (m *TypeSharedResource) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	ctx.AddNamespace("github.com/progpjs/progpAPI")
	return "", "resolveSharedResourceFromDouble(res.currentEvent.id, " + paramName + ")"
}

func (m *TypeSharedResource) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return "    res.value = C.long(goRes.GetId())"
}

//endregion

//region *progpAPI.TypeSharedResourceContainer

type TypeSharedResourceContainer struct {
}

func (m *TypeSharedResourceContainer) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeSharedResourceContainer) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "getSharedResourceContainerFromUIntPtr(res.currentEvent.id)"
}

func (m *TypeSharedResourceContainer) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
	return ""
}

//endregion
