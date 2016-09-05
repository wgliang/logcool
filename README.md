# Logcool [![Version Status](https://img.shields.io/badge/release-v0.1.0-orange.svg)](https://github.com/wgliang/logcool/releases/tag/v0.1.0)

[![Build Status](https://travis-ci.org/wgliang/logcool.svg?branch=master)](https://travis-ci.org/wgliang/logcool)
[![GoDoc](https://godoc.org/github.com/wgliang/logcool?status.svg)](https://godoc.org/github.com/wgliang/logcool)
[![Join the chat at https://gitter.im/logcool/Lobby](https://badges.gitter.im/logcool/Lobby.svg)](https://gitter.im/logcool/Lobby?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/wgliang/logcool)](https://goreportcard.com/report/github.com/wgliang/logcool)
[![License](https://img.shields.io/badge/LICENSE-Apache2.0-ff69b4.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)


Logcool is an open source project to collect, filter ,transfer and response any log or event-flow data as a lightweight system.[中文](./docs/README_ZH.md)

![Logcool](./logcool.jpg)

Logcool's design learn from Heka and Logstash and it's implementation was inspired by gogstash. What's more, the logcool's goal is to be a completely independent project and not much rely on other non-standard libiaries.

Because it is difficult to fully meet the needs of different services, this repository provides basic plugins, such as encryption and decryption of data, compression and decompression of data, data format conversion, support files, command line, http, or the output of any system or redis, influxDB, MySQL database and so on. Importantly, you can easily develop a plugin according to your needs, and easily use it.

And can use logcool any way you can.

## Getting started

Logcool can collect all-types los or event-flow data, and support some input/output types.Besides,you can  new your's plugs if you need it. To get started, [check out the installation instructions in the documentation](https://godoc.org/github.com/wgliang/logcool).

## Using Example

A easy stdin2stdout example. 
![Logcool](./logcool.gif)

## Plugins

Some plugins that had finished and will developed in the future.

### input
- [file](https://github.com/wgliang/logcool/tree/master/input/file)
- [stdin](https://github.com/wgliang/logcool/tree/master/input/stdin)
- [http](https://github.com/wgliang/logcool/tree/master/input/stdin)

### filter
- [zeus](https://github.com/wgliang/logcool/tree/master/filter/zeus)
- metrics

### codec
- aes
- zip
- json

### output
- [stdout](https://github.com/wgliang/logcool/tree/master/output/stdout)
- [redis](https://github.com/wgliang/logcool/tree/master/output/redis)
- influxdb
- mysql
- pg

## Versions

[versions](https://github.com/wgliang/logcool/docs/README_ZH.md)

## Other Contributor

Logcool learn from gogstash much. Thank you for your contribution, and I also learn a lot from your project. @tsaikd

## Licensing

Logcool is licensed under the Apache License, Version 2.0. See LICENSE for the full license text.