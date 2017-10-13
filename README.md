# Token back-end

## Bootstrap

### Installing Go
First you need to install [gvm](https://github.com/moovweb/gvm).

Project use go v1.8.
To install needed package
```bash
gvm install go1.8
```

After that use v1.8 by default
```bash
gvm use go1.8 --default
```
#### Troubleshooting
Usually `gvm` will throw error about installation of v1.8.
To fix that you need to install v1.4 first.

```bash
gvm install go1.4 -B
gvm use go1.4
```

### Installing protobuf

First Install the standard C++ implementation of protocol buffers (protoc) from [here](https://developers.google.com/protocol-buffers/)

Then install packages
```bash
go get github.com/micro/protobuf/{proto,protoc-gen-go}
```

To generate protobuf go files
```bash
make proto-generate
```

### Installing Glide

Glide is a package manager for Go

To install
```bash
curl https://glide.sh/get | sh
```

After that install project dependencies
```bash
cd src
glide install
cd ..
```

Then create package set
```bash
gvm pkgset create token
gvm pkgset use token
```

You need to link the project to `GOPATH`

```bash
rm -rf ~/.gvm/pkgsets/go1.8/token/src
ln -s $PWD/src ~/.gvm/pkgsets/go1.8/token/
```

### Installing linter

Linter is the code analyzer that will give a hints about how code can be improved.

```bash
go get -u gopkg.in/alecthomas/gometalinter.v1
```

Then you need to install linters
```bash
gometalinter.v1 --install
```

That's all here.

### Docker
First of all install the Docker.

Create an internal containers network
```bash
docker network create token_natasha
```

Then init the environment
```bash
make docker-init
```

### Database Migrations
To run database migrations
```bash
make migrate
```

## Development

### GVM

You need to use the appropriate `gvm` package set for development
```bash
gvm pkgset use token
```

### Dependencies

To add new dependencies please use `glide`.
```bash
cd src
glide get *package name*
glide up
cd src ..
```

### Lint

To run lint

```bash
make lint
```

### Tests

To run tests

```bash
make test
```

### Docker containers

To check the running docker containers statuses
```bash
docker ps
```

To start container
```bash
make docker-start-*container name*
```

Available container names
* redis
* kafka
* postgres

Before starting the project please add this line to your `/etc/hosts`
```
token-kafka 127.0.0.1
```

### New module

* Module should be created in `src` folder.
* Add module's folder name to `ignore` array within `glide.yaml` file