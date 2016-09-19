# Logcool [![Version Status](https://img.shields.io/badge/release-v0.1.0-orange.svg)](https://github.com/wgliang/logcool/releases/tag/v0.1.0)

[![Build Status](https://travis-ci.org/wgliang/logcool.svg?branch=master)](https://travis-ci.org/wgliang/logcool)
[![GoDoc](https://godoc.org/github.com/wgliang/logcool?status.svg)](https://godoc.org/github.com/wgliang/logcool)
[![Join the chat at https://gitter.im/logcool/Lobby](https://badges.gitter.im/logcool/Lobby.svg)](https://gitter.im/logcool/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/wgliang/logcool)](https://goreportcard.com/report/github.com/wgliang/logcool)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)


Logcool是一个开源的集日志和事件流数据收集，过滤，传输及响应的轻量级数据采集系统。

![Logcool](../logcool.jpg)

Logcool的设计灵感来自Heka和Logstash，它的实现受到了gogstash的启发。重要的是它的目标在于解决前面系统的一些弊端，例如系统庞大或者不能容易的在业务环境中部署，在这方面gogsyash实现的已经非常好了，不过我不看好的是依赖过多非标准库，在结构上也不够简单清晰。这是为什么我将gogstash重构并重新设计一些逻辑的原因。

由于很难完全满足不同场景下的需求，这个库更多的是提供基础插件，例如数据的加密和解密，数据的压缩和解压缩，数据格式的转换，支持文件，命令行或者其他系统的输出格式，支持redis，influxDB和MySQL数据库等等。重要的是，你可以根据自己的需要轻易的开发符合自己需求的插件，并轻松的使用它。

你可以以任何的方式使用logcool。

## Getting started

Logcool 可以收集各类型的日志和事件数据，并且支持输入／输出以及过滤采用插件的形式注入，所以你喜欢的话可以根据自己的业务改写自己的插件，这是很简单的. To get started, [check out the installation instructions in the documentation](https://godoc.org/github.com/wgliang/logcool).

## Using Example

![Logcool](../logcool.gif)

## Plugins

已经完成和未来会增加的插件：

### input
- [file](https://github.com/wgliang/logcool/tree/master/input/file) 数据来源是文件，例如日志文件
- [stdin](https://github.com/wgliang/logcool/tree/master/input/stdin) 从控制台获取数据，这个调试和示例会用到
- [http](https://github.com/wgliang/logcool/tree/master/input/stdin) 从网络获取数据，支持post，get等
- [collectd](https://github.com/wgliang/logcool/tree/master/input/collectd) 监控系统性能数据，例如CPU，内存，网络，硬盘等等

### filter
- [zeus](https://github.com/wgliang/logcool/tree/master/filter/zeus) 简单的打标签过滤器
- [metrics](https://github.com/wgliang/logcool/tree/master/filter/metrics) 打点计数器，可用于告警和dashboard生成
- [grok](https://github.com/wgliang/logcool/tree/master/filter/grok) 正则过滤数据，支持多模式匹配
- [split](https://github.com/wgliang/logcool/tree/master/filter/split) 根据分隔符分割日志或事件信息，生成命令行参数

### output
- [stdout](https://github.com/wgliang/logcool/tree/master/output/stdout) 标准输出到控制台
- [redis](https://github.com/wgliang/logcool/tree/master/output/redis) 将数据打入redis数据库
- influxdb 数据导入influxdb，这个对于时序数据很有用
- [email](https://github.com/wgliang/logcool/tree/master/output/email)通过email发送消息，比如告警和服务异常通知
- [lexec](https://github.com/wgliang/logcool/tree/master/output/lexec) 发送消息执行命令或脚本
- mysql 将数据写入mysql
- pg 将数据写入pg

## Versions

[版本通知](https://github.com/wgliang/logcool/blob/master/docs/VERSION_UPDATE.md)

## Other Contributor

Logcool 从gogstash而来，在此感谢@tsaikd

## Licensing

Logcool is licensed under the Apache License, Version 2.0. See LICENSE for the full license text.

## Welcome to Contribute


也欢迎修改和完善文档。
