# 斗鱼弹幕本地播放

## 本地安装go编译环境
理论上应该跨平台，我只在mac和linux上测试过

## 依赖包安装
go get github.com/sirupsen/logrus 日志 

go get github.com/faiface/beep    播放声音库 还有点bug


## 阿里云申请语音转换接口的权限
在config.go 中配置 appid和token

## 编译运行
go build

./douyudanmu -room=9999 

