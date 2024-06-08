// Copyright © 2024 Grain. All rights reserved.
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
//go:build arm

package jsonx

import (
	"encoding/json"
	"log"
)

func Marshal(data interface{}) []byte {
	marshal, err := json.Marshal(data)
	if err != nil {
		log.Print(err)
		return nil
	}
	return marshal
}

func Unmarshal(data []byte, v interface{}) error {
	if len(data) == 0 {
		return errors.New("value is empty")
	}

	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}
