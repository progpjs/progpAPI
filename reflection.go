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

import (
	"reflect"
	"strings"
)

func ReflectValueToAny(resV reflect.Value) any {
	if !resV.IsValid() {
		return nil
	} else {
		kind := resV.Kind()

		if kind == reflect.String {
			return resV.String()
		} else if kind == reflect.Bool {
			return resV.Bool()
		} else if kind == reflect.UnsafePointer {
			return resV.UnsafePointer()
		} else if resV.CanFloat() {
			return resV.Float()
		} else if resV.CanInt() {
			return resV.Int()
		} else if resV.CanUint() {
			return resV.Uint()
		} else if resV.CanInterface() {
			return resV.Interface()
		}
	}

	return nil
}

func DynamicCallFunction(toCall any, callArgs []reflect.Value) (result []reflect.Value, errorMsg string) {
	doCall := func() {
		defer func() {
			if recoverValue := recover(); recoverValue != nil {
				if errMsg, ok := recoverValue.(string); ok {
					if strings.HasPrefix(errMsg, "reflect:") {
						errorMsg = errMsg
						return
					}
				}

				panic(recoverValue)
			}
		}()

		result = reflect.ValueOf(toCall).Call(callArgs)
	}

	doCall()

	return result, errorMsg
}
