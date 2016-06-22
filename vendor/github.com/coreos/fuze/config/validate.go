// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"reflect"
)

type ErrKeysUnrecognized []string

func (e ErrKeysUnrecognized) Error() string {
	return fmt.Sprintf("unrecognized keys: %v", []string(e))
}

func assertKeysValid(value interface{}, refType reflect.Type) ErrKeysUnrecognized {
	var err ErrKeysUnrecognized

	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}
	switch value.(type) {
	case map[interface{}]interface{}:
		ks := value.(map[interface{}]interface{})
	keys:
		for key := range ks {
			for i := 0; i < refType.NumField(); i++ {
				sf := refType.Field(i)
				tv := sf.Tag.Get("yaml")
				if tv == key {
					if serr := assertKeysValid(ks[key], sf.Type); serr != nil {
						err = append(err, serr...)
					}
					continue keys
				}
			}

			err = append(err, fmt.Sprintf("%v", key))
		}
	case []interface{}:
		ks := value.([]interface{})
		for i := range ks {
			if serr := assertKeysValid(ks[i], refType.Elem()); serr != nil {
				err = append(err, serr...)
			}
		}
	default:
	}

	return err
}
