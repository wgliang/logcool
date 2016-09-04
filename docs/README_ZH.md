# Logcool [![Version Status](https://img.shields.io/badge/release-v0.1.0-orange.svg)](https://github.com/wgliang/logcool/releases/tag/v0.1.0)

[![Build Status](https://travis-ci.org/wgliang/logcool.svg?branch=master)](https://travis-ci.org/wgliang/logcool)
[![GoDoc](https://godoc.org/github.com/wgliang/logcool?status.svg)](https://godoc.org/github.com/wgliang/logcool)
[![Join the chat at https://gitter.im/logcool/Lobby](https://badges.gitter.im/logcool/Lobby.svg)](https://gitter.im/logcool/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/wgliang/logcool)](https://goreportcard.com/report/github.com/wgliang/logcool)


Logcool是一个开源的集日志和事件流数据收集，过滤，传输及响应的轻量级数据采集系统。

![Logcool](../logcool.jpg)

Logcool的设计灵感来自Heka和Logstash，它的实现受到了gogstash的启发。重要的是它的目标在于解决前面系统的一些弊端，例如系统庞大或者不能容易的在业务环境中部署，在这方面gogsyash实现的已经非常好了，不过我不看好的是依赖过多非标准库，在结构上也不够简单清晰。这是为什么我将gogstash重构并重新设计一些逻辑的原因。

Logcool 目前还处于“婴儿期”，所以未来会有很大改变，不仅在设计上，在代码上也会有重大改变，所以非常不建议目前直接用于生产环境中。

## Getting started

Logcool 可以收集各类型的日志和事件数据，并且支持输入／输出以及过滤采用插件的形式注入，所以你喜欢的话可以根据自己的业务改写自己的插件，这是很简单的. To get started, [check out the installation instructions in the documentation](https://godoc.org/github.com/wgliang/logcool).

## Using Example

![Logcool](../logcool.gif)

## Other Contributor

Logcool 从gogstash而来，在此感谢@tsaikd

## Licensing

Logcool is licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
