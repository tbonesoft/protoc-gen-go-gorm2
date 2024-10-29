# protoc-gen-go-gorm2

A modern protobuf compiler plugin designed to generate GORM models for simple object persistence tasks.

Features

- [x] Support modern toolchain: Go v1.22 or newer, Protocol Buffers v3;
- [x] Go field decorators/tags can be defined by option in the .proto file for GORM;
- [x] No force conventions;
- [ ] Converters between PB version and ORM version of the objects are included.

Usage

- [Get the Protocol Buffers compiler](https://github.com/protocolbuffers/protobuf?tab=readme-ov-file#protobuf-compiler-installation);
- Get the golang GORM code generator: `go install github.com/tbonesoft/protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2@latest`
- Define proto file likes examples/bookstore/proto/bookstore/v1/bookstore.proto
- Compile it:

```
protoc ^
--go_out=gen/ ^
--go_opt=paths=source_relative ^
--go-gorm2_out=gen/ ^
--go-gorm2_opt=paths=source_relative ^
--go-gorm2_opt=engine=postgres ^
--proto_path=examples/bookstore/proto/ ^
--proto_path=proto/ ^
examples/bookstore/proto/bookstore/v1/bookstore.proto
```