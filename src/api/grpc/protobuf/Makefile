# https://github.com/gogo/protobuf/blob/master/test/example/Makefile
#
regenerate:
#	(protoc -I=. -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf -I=$(GOPATH)/src --gogo_out=. cbtumblebug.proto)
# protoc -I . -I $(GOPATH)/src/github.com/gogo/protobuf/protobuf -I $(GOPATH)/src --gofast_out=plugins=grpc,paths=source_relative,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types:. *.proto

	# gogo proto compile
	protoc \
			./gogoproto/gogo.proto \
			-I . \
			-I $(GOPATH)/src/github.com/gogo/protobuf/protobuf \
			-I $(GOPATH)/src \
			--gofast_out=plugins=grpc,paths=source_relative,\
	Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:\
	.

	# cbtumblebug proto compile
	protoc \
			./cbtumblebug/cbtumblebug.proto \
			-I . \
			-I $(GOPATH)/src/github.com/gogo/protobuf/protobuf \
			-I $(GOPATH)/src/github.com/cloud-barista/cb-tumblebug/src/api/grpc/protobuf \
			--gofast_out=plugins=grpc:\
	.	