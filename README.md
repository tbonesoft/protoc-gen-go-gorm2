# protoc-gen-go-gorm2

A modern protobuf compiler plugin designed to generate GORM models for simple object persistence tasks.

[![Go](https://github.com/tbonesoft/protoc-gen-go-gorm2/actions/workflows/go.yml/badge.svg)](https://github.com/tbonesoft/protoc-gen-go-gorm2/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tbonesoft/protoc-gen-go-gorm2)](https://goreportcard.com/report/github.com/tbonesoft/protoc-gen-go-gorm2)
[![LICENSE](https://img.shields.io/github/license/tbonesoft/protoc-gen-go-gorm2.svg?style=flat-square)](https://github.com/tbonesoft/protoc-gen-go-gorm2/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/tbonesoft/protoc-gen-go-gorm2?status.svg)](https://godoc.org/github.com/tbonesoft/protoc-gen-go-gorm2)

Features

- [x] Support modern toolchain: Protocol Buffers v3;
- [x] Go field decorators/tags can be defined by option in the .proto file for GORM;
- [x] No force conventions;
- [x] Converters between PB version and ORM version of the objects are included.

Build requirement: Go 1.20.

Usage

- [Get the Protocol Buffers compiler](https://github.com/protocolbuffers/protobuf?tab=readme-ov-file#protobuf-compiler-installation);
- Get the golang GORM code generator: `go install github.com/tbonesoft/protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2@latest`
- Define a proto file likes example/bookstore/proto/bookstore/v1/bookstore.proto
- Compile it:

```shell
protoc ^
--go_out=gen/ ^
--go_opt=paths=source_relative ^
--go-gorm2_out=gen/ ^
--go-gorm2_opt=paths=source_relative ^
--go-gorm2_opt=engine=postgres ^
--proto_path=proto/ ^
example/bookstore/proto/bookstore/v1/bookstore.proto
```

TODO

- support all PostgreSQL [data types](https://www.postgresql.org/docs/17/datatype.html)
- more examples
- 100% unit tests

## Examples

Generate code for Go backend

```shell
go install -v github.com/tbonesoft/protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2
buf generate --template buf.gen.bookstore.yaml
```

Generate code for frontend

```shell
npx buf generate --template buf.gen.bookstore.frontend.yaml
```

## Notes on buf.build

As of buf v2, buf does not support setting different output directories for generated code for multiple modules, for example:

`proto/gorm/v1/gorm.proto` outputs to

```shell
proto/gorm/v1/gorm.pb.go
```

`proto/bookstore/v1/bookstore.proto` outputs to

```shell
gen/bookstore/v1/bookstore.pb.go
gen/bookstore/v1/bookstore_gorm.pb.go
```
