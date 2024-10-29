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
- Define a proto file likes proto/bookstore/v1/bookstore.proto
- Compile it:

```
protoc ^
--go_out=gen/ ^
--go_opt=paths=source_relative ^
--go-gorm2_out=gen/ ^
--go-gorm2_opt=paths=source_relative ^
--go-gorm2_opt=engine=postgres ^
--proto_path=proto/ ^
proto/bookstore/v1/bookstore.proto
```

## Notes on buf.build

As of buf v2, buf does not support setting different output directories for generated code for multiple modules, for example:

`proto/gorm/v1/gorm.proto` outputs to

```
proto/gorm/v1/gorm.pb.go
```

`proto/bookstore/v1/bookstore.proto` outputs to

```
gen/bookstore/v1/bookstore.pb.go
gen/bookstore/v1/bookstore_gorm.pb.go
```

Request for improvement to official developers will not be accepted in 2024-10. See also

- https://app.slack.com/client/TS9DC8PPX/CRZ680FUH
- https://gist.github.com/kakuiho/0a8064b2626e9e352aab6b9afbba025d
