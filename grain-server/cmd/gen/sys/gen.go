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

package sys

import (
	sysModel "github.com/go-grain/grain/model/system"
	"gorm.io/gen"
)

func Gen() {
	xGen(
		"internal/repo/system/query",
		sysModel.SysRole{},
		sysModel.SysUser{},
		sysModel.Upload{},
		sysModel.CasbinRule{},
		sysModel.SysApi{},
		sysModel.SysMenu{},
		sysModel.Upload{},
		sysModel.Project{},
		sysModel.Models{},
		sysModel.Fields{},
		sysModel.Organize{},
	)
}

func xGen(outPath string, model ...interface{}) {
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		model...,
	)
	g.Execute()
}
