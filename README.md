# Akamai DataStream CLI

[![goreportcard](https://goreportcard.com/badge/github.com/coralogix/akamai-datastream-cli)](https://goreportcard.com/report/github.com/coralogix/akamai-datastream-cli)
[![godoc](https://img.shields.io/badge/godoc-reference-brightgreen.svg?style=flat)](https://godoc.org/github.com/coralogix/akamai-datastream-cli)
[![license](https://img.shields.io/github/license/coralogix/akamai-datastream-cli.svg)](https://raw.githubusercontent.com/coralogix/akamai-datastream-cli/master/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/coralogix/akamai-datastream-cli.svg)](https://github.com/coralogix/akamai-datastream-cli/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/coralogix/akamai-datastream-cli.svg)](https://github.com/coralogix/akamai-datastream-cli/pulls)
[![GitHub contributors](https://img.shields.io/github/contributors/coralogix/akamai-datastream-cli.svg)](https://github.com/coralogix/akamai-datastream-cli/graphs/contributors)

`Akamai DataStream CLI` is the utility to collect logs from `Akamai DataStream API` and send it [Coralogix](https://coralogix.com/).

## Prerequisites

Before beginning you must have installed:

* `docker` (`^18.09`)
* `docker-compose` (`^1.27.3`)
* `cmake`

## Usage

Complete next steps:

1. Clone/download this repository
2. Copy `examples/akamai.env.example` to `configs/akamai.env` and fill with `Akamai` credentials and datastream settings ([Akamai DataStream API](https://developer.akamai.com/api/web_performance/datastream/v1.html))
3. Copy `examples/coralogix.env.example` to `configs/coralogix.env` and fill with `Coralogix` credentials

Then deploy logging agent to `Docker`:

```bash
$ make deploy
```

## License

This project is licensed under the Apache-2.0 License.