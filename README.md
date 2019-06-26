# 介绍

sku-manage-go是一个使用golang编写程序。实现的功能主要是用于方便网店卖家来管理商品sku，1688上的采购链接。以及生成采购单。

数据库操作使用gorm 

路由使用了mux

# 运行

运行程序之前需要修改一下mysql的配置,并且确保mysql中有一个test的数据库(可以通过mysql连接字符串修改)。

然后使用`go build` 生成可执行文件，运行

# 设置window开机自动启动

因为这个项目需要部署在个人PC上面，并且需要开机启动。通过windows的执行计划，可以很方便的设置开机启动

先将启动命令写成一个批处理程序 start.bat
```
@echo off
d:
cd "D:\code\src\sku-manage"
start /b sku-manage.exe --debug
exit
```

- 右键`我的电脑` 选择`管理`，可以看到`任务计划程序`，点击小箭头可以看到`任务计划程序库` 

- 创建一个任务，在常规中选择`不管用户是否登录都要运行` 并且选择存储密码，使用最高权限运行，触发器选择启动时，操作就输入对应批处理脚本的位置

- 打开`控制面板`、`管理工具`，然后打开`本地安全策略`。在`本地安全策略`窗口中，依次单击`本地策略`、`用户权限分配`，然后单击`作为批处理作业登录`，将此执行此计划任务的用户增加到列表里
