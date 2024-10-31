all:	build-protoc-gen-go-gorm2 build-examples

build-protoc-gen-go-gorm2:
	protoc \
	--go_out=proto/ \
	--go_opt=paths=source_relative \
	--proto_path=proto/ \
	proto/gorm/v1/gorm.proto

build-examples:
	go install -v github.com/tbonesoft/protoc-gen-go-gorm2/cmd/protoc-gen-go-gorm2

	test -d gen || mkdir gen
 
	protoc \
	--go_out=gen/ \
	--go_opt=paths=source_relative \
	--go-gorm2_out=gen/ \
	--go-gorm2_opt=paths=source_relative \
	--go-gorm2_opt=engine=postgres \
	--proto_path=examples/bookstore/proto/ \
	--proto_path=proto/ \
	proto/gorm/v1/gorm.proto \
	examples/bookstore/proto/bookstore/v1/bookstore.proto

push-bsr:
	buf push buf.build/tbonesoft/protoc-gen-go-gorm2
