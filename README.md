<div align=center>
<img src="http://grain.gitbili.com/uploads/systemFile/2023/8-26/grain-logo-v2.png" alt=""/>
</div>
<div align=center>
<img src="https://img.shields.io/badge/Go-v1.20-blue" alt=""/>
<img src="https://img.shields.io/badge/Gin-v1.9.1-lightBlue" alt=""/>
<img src="https://img.shields.io/badge/Gorm-v1.25.2-red" alt=""/>
<img src="https://img.shields.io/badge/Gen-v0.3.23-lightgred" alt=""/>
</div>


#### Grain是什么:
Grain是一个基于Gin + Gorm&Gen + Vue + ArcoDesign开发的前后端分离的可开箱即用的中后台管理基础脚手架

#### Grain能做什么:
- Grain致力于提供最简单的开箱即用的脚手架基础平台,帮助用户快速搭建中后台管理系统
- Grain使用 JWT 进行身份验证。使用 Casbin 实现基于角色的访问控制，控制资源权限。前端菜单根据角色权限动态显示对应的权限菜单。
- Grain能够根据输入的模型字段自动生成前端 CRUD 管理页面和后端 CRUD,API 接口

#### Grain基础通用功能：
- 代码助手:
- 用户管理:
- API管理:
- 菜单管理:
- 角色管理:

#### 使用介绍
写完后在开始写文档...
##### 工程目录结构
    ├─base          #放一些基础的东西
    ├─cmd           #程序入口
    │  └─gen        #存放Gorm Gen build文件
    ├─config        #配置文件
    ├─core          #程序核心启动层
    ├─docs          #swagger文档
    ├─log           #日志记录层
    ├─middleware    #中间件
    ├─model         #数据模型层
    ├─repo          #repo层
    ├─internal      #存放不希望外部调用的东西
    │  ├─handler    #api层
    │  ├─repo       #repo层
    │  ├─router     #路由层
    │  └─service    #业务逻辑处理层
    ├─stencil       #模版
    └─utils         #工具包

#### 讨论群
备注:Grain
<div align=start>
<img src="http://grain.gitbili.com/uploads/systemFile/2023/7-22/wx.png" width=200" height="200" />
</div>

#### 特别的非感谢
- [ArcoDesignPro](https://arco.design/)(感谢提供了开箱即用的ArcoPro)
- [casbin](https://github.com/casbin/casbin)(感谢提供资源对象鉴权能力)
- [gin](https://github.com/gin-gonic/gin)（感谢提供出色的Web和API开发框架）
- [go-redis](https://github.com/redis/go-redis)(感谢提供进行Redis数据库交互能力)
- [viper](https://github.com/spf13/viper)(感谢提供友好的配置解决方案)
- [gin-swagger](https://github.com/swaggo/gin-swagger)(感谢提供在gin框架中生成Swagger文档的库)
- [gorm&Gen](https://gorm.io)(感谢提供ORM和Gen代码生成能力)
- [gva](https://www.gin-vue-admin.com/)（非常感谢三水哥的gva项目，它为我提供了许多有价值的借鉴和思路(^*狗头*^)）