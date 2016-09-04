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

## Other Contributor

Logcool 从gogstash而来，在此感谢@tsaikd

## Licensing

Logcool is licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
