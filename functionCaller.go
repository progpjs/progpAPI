package progpAPI

func GetFunctionCaller(functionTemplate string) any {
	if gSelectedScriptEngine == nil {
		return nil
	}

	return gSelectedScriptEngine.GetFunctionCaller(functionTemplate)
}

func DynamicFunctionCaller(params ...any) (any, error) {
	// TODO
	return nil, nil
}
