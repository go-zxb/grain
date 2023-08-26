<div align=center>
<img src="http://grain.gitbili.com/uploads/systemFile/2023/8-26/logo.png"/>
</div>
<div align=center>
<img src="https://img.shields.io/badge/golang-v1.20-blue"/>
<img src="https://img.shields.io/badge/gin-v1.9.1-lightBlue"/>
<img src="https://img.shields.io/badge/gorm-v1.25.2-red"/>
</div>


# Grain 
Grain 是一个基于 Gin + Gorm&Gen + vue + ArcoDesign开发的前后端分离的系统管理基础脚手架
#### 项目特点:
- JWT(身份验证)
- Casbin(使用RBAC控制资源权限)
- 前端动态菜单(根据角色权限获取相应的菜单)
- 代码生成器(提供了前后端 CRUD 代码生成能力,只需输入数据模型 field 字段自动生成前端 CRUD 管理页面 和 后端CRUD Api接口)

#### 基础通用功能
- 用户管理模块
- Api管理模块
- 菜单管理模块
- 角色管理模块

# 使用介绍
写完后在开始写文档...
# 项目工程
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

# 讨论群
备注:Grain 
<div align=center>
<img src="http://grain.gitbili.com/uploads/systemFile/2023/7-22/wx.png" width=300" height="300" />
</div>