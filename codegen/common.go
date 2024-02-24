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

type IsTypeHandler interface {
	// CppToCgoParamCall allows casting const char* to char*.
	// Strings coming from v8 are const char*, which allow to remember that they must not be deleted.
	// But CGo requires char* string, so we must cast the parameter before calling.
	//
	// ===> Inside "codeBinding.cpp":
	//
	//		void v8Function_g_progpAdmin_f_pmyFunction(
	//		   ProgpV8ContextPtr ctxPtr,
	//	       const v8::FunctionCallbackInfo<v8::Value> &callInfo,
	//	       const v8::Local<v8::Context> &v8Ctx,
	//	       v8::Isolate *v8Iso) {
	//
	//	   V8CALLARG_EXPECT_ARGCOUNT(1);
	//	   V8CALLARG_EXPECT_CSTRING(p0, 0);
	//
	//		ProgpFunctionReturnVoid resWrapper{};
	//		progpCgoBinding__g_progpAdmin_f_progpSendSignal(&resWrapper, (char*)p0);
	//																	 [^HERE^^^^^^^^^]
	//
	// .
	CppToCgoParamCall(paramName string, ctx *ProgpV8CodeGenerator) string

	// V8ToCppDecoder decode the V8 call parameter from V8 to C++.
	//
	// ===> Inside "codeBinding.cpp":
	//
	//		void v8Function_g_progpDefault_f_myFunction(
	//
	//			ProgpV8ContextPtr ctxPtr,
	//	       	const v8::FunctionCallbackInfo<v8::Value> &callInfo,
	//	       	const v8::Local<v8::Context> &v8Ctx,
	//	       	v8::Isolate *v8Iso) {
	//	   			V8CALLARG_EXPECT_ARGCOUNT(1);					<--- HERE
	//	   			V8CALLARG_EXPECT_FUNCTION(p0, 0);				<--- HERE
	//
	// .
	V8ToCppDecoder(ctx *ProgpV8CodeGenerator) string

	// ReturnTypeWrapper is for the Go and the C++ part.
	// It's allowing our Go function to return a value to the C++ part
	// by avoiding some strange CGo behaviors where directly sending
	// the value with a "return" randomly breaks things.
	//
	// ===> Inside "codeBinding.cpp":
	//
	//		void v8Function_g_progpDefault_f_myFunction(
	//				ProgpV8ContextPtr ctxPtr,
	//	       		const v8::FunctionCallbackInfo<v8::Value> &callInfo,
	//	       		const v8::Local<v8::Context> &v8Ctx,
	//	       		v8::Isolate *v8Iso) {
	//
	//					ProgpFunctionReturnInt resWrapper{};					<--- HERE
	//					progpCgoBinding__g_progpDefault_f_myFunction(&resWrapper);
	//
	// ===> Inside "codeBinding.go"
	//
	//	func progpCgoBinding__g_test1_f_myFunction(res *C.ProgpFunctionReturnInt) {
	//												  	  [^HERE^^^^^^^^^]
	//
	// .
	ReturnTypeWrapper(ctx *ProgpV8CodeGenerator) string

	// ReturnTypeEncoder allows converting a C++ value to a V8 value.
	// It's used in order to encode the return type.
	//
	// ===> Inside "codeBinding.cpp":
	//
	//		progpCgoBinding__g_progpDefault_f_returnInteger(&resWrapper);
	//		// ...
	//	   	auto res = resWrapper.value;
	//	   	callInfo.GetReturnValue().Set(V8VALUE_FROM_INT64(res));
	//									  [^HERE^^^^^^^^^]
	//
	// .
	ReturnTypeEncoder(ctx *ProgpV8CodeGenerator) string

	// CgoFunctionParamType allowing to know the parameters types of the CGo function.
	//
	// ===> Inside "codeBinding.go":
	//
	//	func progpCgoBinding__g_progpDefault_f_myFunctionName(res *C.ProgpFunctionReturnVoid, p0 C.int) {
	//		                                                                                     [^ HERE]
	//
	// .
	CgoFunctionParamType(ctx *ProgpV8CodeGenerator) string

	// CppArgResourcesFreeing free the resources generated by getV8ToCppDecoder
	CppArgResourcesFreeing(paramName string, ctx *ProgpV8CodeGenerator) string

	// CgoToGoDecoding allows decoding a CGo value (pseudo C from Go) and get a pure Go value.
	// It's return two parts. The first one allows "long decoding" and the second "inline decoding".
	//
	// ===> Inside "codeBinding.go":
	//
	//			func progpCgoBinding__g_progpDefault_f_progpExecuteFile(res *C.ProgpFunctionReturnVoid, p1 *C.char, p2 C.int) {
	//	   		var p2_asBool = true				<==
	//	   		if C.int(p2) == 0 {					<== Here : long decoding sample (fist value returned)
	//	       		p2_asBool = false				<==
	//	   		}									<==
	//
	//				JSProgpExecuteFile(resolveV8Context(res.ctx, int(res.contextId)), C.GoString(p1), p2_asBool)
	//	. 												                              [^ HERE: inline decoding (second value returned)]
	//	}
	//
	// .
	CgoToGoDecoding(paramName string, ctx *ProgpV8CodeGenerator) (string, string)

	// GoValueToCgoValue allows converting the return type from pure Go to CGo (pseudo C type).
	// It's doing the inverse of cgoToGoDecoding.
	//
	// ===> Inside "codeBinding.go":
	//
	//	func progpCgoBinding__g_progpDefault_f_returnBool(res *C.ProgpFunctionReturnInt) {
	//			goRes := progpModSample.JsReturnBool()
	//	   		if goRes {									<== HERE
	//	       		res.value = C.int(1)					<== HERE
	//	   		} else {									<== HERE
	//	       		res.value = C.int(0)					<== HERE
	//			}											<== HERE
	//	}
	//
	// .
	GoValueToCgoValue(ctx *ProgpV8CodeGenerator) string
}

type IsFunctionCallerSupportedType interface {
	FcCppToV8Encoder(paramId int) string
	FcCppFunctionHeader(paramId int) string
	FcGoToCppCallParam(paramId int) string
	FcGoToCppConv(paramId int) string
}
