# nano [![Build Status][7]][8] [![Coverage Status][9]][10] [![GoDoc][1]][2] [![Docs][11]][12] [![Go Report Card][3]][4] [![MIT licensed][5]][6]

---

[1]: https://godoc.org/github.com/topfreegames/nano?status.svg
[2]: https://godoc.org/github.com/topfreegames/nano
[3]: https://goreportcard.com/badge/github.com/topfreegames/nano
[4]: https://goreportcard.com/report/github.com/topfreegames/nano
[5]: https://img.shields.io/badge/license-MIT-blue.svg
[6]: LICENSE
[7]: https://github.com/topfreegames/nano/actions/workflows/tests.yaml/badge.svg
[8]: https://github.com/topfreegames/nano/actions/workflows/tests.yaml
[9]: https://coveralls.io/repos/github/topfreegames/nano/badge.svg?branch=master
[10]: https://coveralls.io/github/topfreegames/nano?branch=master
[11]: https://readthedocs.org/projects/nano/badge/?version=latest
[12]: https://nano.readthedocs.io/en/latest/?badge=latest

Nano is an simple, fast and lightweight game server framework with clustering support and client libraries for iOS, Android, Unity and others through the [C SDK](https://github.com/topfreegames/libpitaya).
It provides a basic development framework for distributed multiplayer games and server-side applications.

## Getting Started

### Prerequisites

* [Go](https://golang.org/) >= 1.16
* [etcd](https://github.com/coreos/etcd) (optional, used for service discovery)
* [nats](https://github.com/nats-io/nats.go) (optional, used for sending and receiving rpc)
* [docker](https://www.docker.com) (optional, used for running etcd and nats dependencies on containers)

### Installing

clone the repo

```bash
git clone https://github.com/nut-game/nano.git
```

setup nano dependencies

```bash
make setup
```

### Hacking nano

Here's one example of running nano:

Start etcd (This command requires docker-compose and will run an etcd container locally. An etcd may be run without docker if preferred.)

```
cd ./examples/testing && docker compose up -d etcd
```
run the connector frontend server from cluster_grpc example
```
make run-cluster-grpc-example-connector
```
run the room backend server from the cluster_grpc example
```
make run-cluster-grpc-example-room
```

Now there should be 2 nano servers running, a frontend connector and a backend room. To send requests, use a REPL client for nano [pitaya-cli](https://github.com/topfreegames/pitaya/tree/main/pitaya-cli).

```bash
$ nano-cli
Nano REPL Client
>>> connect localhost:3250
connected!
>>> request room.room.entry
>>> sv-> {"code":0,"result":"ok"}
```

## Running the tests

```bash
make test
```

This command will run both unit and e2e tests.


## Authors

* **TFG Co** - Initial work

## License

[MIT License](./LICENSE)

## Acknowledgements

* [nano](https://github.com/lonnng/nano) authors for building the framework pitaya is based on.
* [pomelo](https://github.com/NetEase/pomelo) authors for the inspiration on the distributed design and protocol

## Security

If you have found a security vulnerability, please email security@tfgco.com

## Resources

- Other nano-related projects
  + [nano-admin](https://github.com/topfreegames/pitaya-admin)
  + [nano-cli](https://github.com/nut-game/nano-cli)

- Documents
  + [API Reference](https://godoc.org/github.com/nut-game/nano)

- Demo
  + [nano cluster mode example](./examples/demo/cluster)
