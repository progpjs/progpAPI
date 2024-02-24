package codegen

import (
	"github.com/progpjs/progpAPI/v2"
	"strings"
)

func GetFunctionCaller(functionTemplate any) any {
	// >>> Extract function signature
	rf := progpAPI.RegisteredFunction{GoFunctionRef: functionTemplate}
	res, err := progpAPI.ParseGoFunction(&rf)
	if err != nil {
		panic(err)
	}

	signature := strings.Join(res.ParamTypes, ",")
	signature = "(" + signature + "):" + res.ReturnType

	// >>> Returns the existing function

	println("signature:", signature)

	// Here it's the compiled version referencing the real version.
	existingCaller := gFunctionCallerMap[signature]
	if existingCaller != nil {
		return existingCaller
	}

	// >>> Add to the function which need to be created

	if res.ReturnErrorOffset != -1 {
		panic("Returning an error value isn't supported")
	}

	if (len(res.ParamTypes) == 0) || (res.ParamTypes[0] != "progpAPI.JsFunction") {
		panic("The first parameter must be of type progpAPI.JsFunction")
	}

	for i, pType := range res.ParamTypes {
		if i == 0 {
			continue
		}

		if !isAllowedFunctionType(pType) {
			panic("Parameter not supported: " + pType)
		}
	}

	gFunctionCallerToBuild[signature] = &functionCallerToBuild{
		paramTypes: res.ParamTypes,
		returnType: res.ReturnType,
	}

	gHasFunctionCallerToBuild = true

	return nil
}

func getAllFunctionCallerToBuild() map[string]*functionCallerToBuild {
	if gHasFunctionCallerToBuild {
		return gFunctionCallerToBuild
	}

	return nil
}

func isAllowedFunctionType(paramType string) bool {
	if paramType == "string" {
		return true
	}

	return false
}

type functionCallerToBuild struct {
	paramTypes []string
	returnType string
}

var gHasFunctionCallerToBuild = false
var gFunctionCallerMap = make(map[string]any)
var gFunctionCallerToBuild = make(map[string]*functionCallerToBuild)
