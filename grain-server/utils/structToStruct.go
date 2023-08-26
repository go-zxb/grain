// Copyright © 2023 Grain. All rights reserved.
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

package utils

import (
	xjson "github.com/go-grain/go-utils/json"
)

// StructToStruct
// 简单粗暴的实现一个结构体数据转到另一个结构体,不知道会不会有啥Bug
// 序列化s1到数组字节类型,在放序列化到s2
func StructToStruct(s1 any, s2 any) error {
	bytes := xjson.Marshal(&s1)
	err := xjson.Unmarshal(bytes, &s2)
	if err != nil {
		return err
	}
	return nil
}
