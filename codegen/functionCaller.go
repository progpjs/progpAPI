package codegen

import (
	"github.com/progpjs/progpAPI/v2"
	"reflect"
	"strings"
)

func GetFunctionSignature(reflectFct reflect.Type) string {
	res, err := progpAPI.ParseGoFunctionReflect(reflectFct, "(function caller)")
	if err != nil {
		panic(err)
	}

	signature := strings.Join(res.ParamTypes, ",")
	signature = "(" + signature + "):" + res.ReturnType

	return signature
}

func AddFunctionCallerToGenerate(reflectFct reflect.Type) {
	// >>> Extract function signature
	res, err := progpAPI.ParseGoFunctionReflect(reflectFct, "")
	if err != nil {
		panic(err)
	}

	signature := strings.Join(res.ParamTypes, ",")
	signature = "(" + signature + ")"

	// >>> Add to the function which need to be created

	if (res.ReturnErrorOffset != -1) || (res.ReturnType != "") {
		panic("Returning a value isn't supported")
	}

	// param[0] is the interface type and is automatically added by Go.
	// So we directly test param 1.
	//
	// Strange behaviors "progpAPI.JsFunction" began "progpV8Engine.v8Function"
	//
	if (len(res.ParamTypes) <= 1) || (res.ParamTypeRefs[1].String() != "progpAPI.JsFunction") {
		panic("The first parameter must be of type progpAPI.JsFunction")
	}

	gFunctionCallerToBuildMap[signature] = &functionCallerToBuild{
		paramTypes: res.ParamTypes,
		returnType: res.ReturnType,
	}

	gHasFunctionCallerToBuild = true
}

func getAllFunctionCallerToBuild() map[string]*functionCallerToBuild {
	if gHasFunctionCallerToBuild {
		return gFunctionCallerToBuildMap
	}

	return nil
}

type functionCallerToBuild struct {
	paramTypes []string
	returnType string
}

var gHasFunctionCallerToBuild = false
var gFunctionCallerToBuildMap = make(map[string]*functionCallerToBuild)
