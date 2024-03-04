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

import "fmt"

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

//region progpAPI.JsFunction

type TypeJsFunction struct {
}

func (m *TypeJsFunction) CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string {
	return paramName
}

func (m *TypeJsFunction) CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeJsFunction) V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string {
	return "V8CALLARG_EXPECT_FUNCTION"
}

func (m *TypeJsFunction) ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeJsFunction) ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string {
	return ""
}

func (m *TypeJsFunction) CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string {
	return "C.ProgpV8FunctionPtr"
}

func (m *TypeJsFunction) CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string) {
	return "", "newV8Function(res.isAsync, " + paramName + ", res.currentEvent)"
}

func (m *TypeJsFunction) GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string {
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
	ctx.AddNamespace("github.com/progpjs/progpAPI/v2")
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

//region >>> For function caller

//region string

func (m *TypeString) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(
		"    argArray[%d] = v8::String::NewFromUtf8(v8Iso, p%d_val, v8::NewStringType::kNormal, (int)p%d_size).ToLocalChecked();\n", paramId, paramId, paramId)

}

func (m *TypeString) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(", const char* p%d_val, size_t p%d_size", paramId, paramId)
}

func (m *TypeString) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                (*C.char)(unsafe.Pointer(unsafe.StringData(p%d))), C.size_t(len(p%d)),", paramId, paramId)
}

func (m *TypeString) FcGoToCppConvCache(paramId int) string {
	return ""
}

//endregion

//region bool

func (m *TypeBool) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(
		"    argArray[%d] = BOOL_TO_V8VALUE(p%d);\n", paramId, paramId)
}

func (m *TypeBool) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(",int p%d", paramId)
}

func (m *TypeBool) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                asCBool(p%d),", paramId)
}

func (m *TypeBool) FcGoToCppConvCache(_ int) string {
	return ""
}

//endregion

//region float64

func (m *TypeFloat64) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(
		"    argArray[%d] = DOUBLE_TO_V8VALUE(p%d);\n", paramId, paramId)
}

func (m *TypeFloat64) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(", double p%d", paramId)
}

func (m *TypeFloat64) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                (C.double)(p%d),", paramId)
}

func (m *TypeFloat64) FcGoToCppConvCache(_ int) string {
	return ""
}

//endregion

//region []uint / ArrayBuffer

func (m *TypeUIntArray) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(`
	auto p%d_bs = std::shared_ptr(v8::ArrayBuffer::NewBackingStore(v8Iso, p%d_size));
	memcpy(p%d_bs->Data(), p%d_buffer, p%d_size);
	argArray[%d] = v8::ArrayBuffer::New(v8Iso, p%d_bs);

`, paramId, paramId, paramId, paramId, paramId, paramId, paramId)
}

func (m *TypeUIntArray) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(", const char* p%d_buffer, size_t p%d_size", paramId, paramId)
}

func (m *TypeUIntArray) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                (*C.char)(unsafe.Pointer(&p%d[0])), C.size_t(len(p%d)),", paramId, paramId)
}

func (m *TypeUIntArray) FcGoToCppConvCache(_ int) string {
	return ""
}

//endregion

//region progpAPI.StringBuffer

func (m *TypeStringBuffer) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(
		"    argArray[%d] = v8::String::NewFromUtf8(v8Iso, (char*)p%d_val, v8::NewStringType::kNormal, (int)p%d_size).ToLocalChecked();\n", paramId, paramId, paramId)

}

func (m *TypeStringBuffer) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(", const char* p%d_val, size_t p%d_size", paramId, paramId)
}

func (m *TypeStringBuffer) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                (*C.char)(unsafe.Pointer(&p%d[0])), C.size_t(len(p%d)),", paramId, paramId)
}

func (m *TypeStringBuffer) FcGoToCppConvCache(_ int) string {
	return ""
}

//endregion

//region *progpAPI.SharedResource

func (m *TypeSharedResource) FcCppToV8Encoder(paramId int) string {
	return fmt.Sprintf(
		"    argArray[%d] = DOUBLE_TO_V8VALUE(p%d);\n", paramId, paramId)
}

func (m *TypeSharedResource) FcCppFunctionHeader(paramId int) string {
	return fmt.Sprintf(", double p%d", paramId)
}

func (m *TypeSharedResource) FcGoToCppCallParam(paramId int) string {
	return fmt.Sprintf("\n                (C.double)(p%d.GetId()),", paramId)
}

func (m *TypeSharedResource) FcGoToCppConvCache(_ int) string {
	return ""
}

//endregion

//endregion
