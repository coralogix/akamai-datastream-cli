# Akamai DataStream CLI

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