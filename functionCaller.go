package progpAPI

func GetFunctionCaller(functionTemplate string) any {
	if gSelectedScriptEngine == nil {
		return nil
	}

	return gSelectedScriptEngine.GetFunctionCaller(functionTemplate)
}
