c_proto_read:
	go build -o ./bin/proto_read.exe ./main/proto_tool/proto_read/MsgReader.go

c_proto_compile:
	go build -o ./bin/proto_compile.exe ./main/proto_tool/proto_compile/proto_compile.go
	
c_itf2proto:
	go build -o ./bin/itf2proto.exe ./main/proto_tool/proto_itf/to_proto/ItfToProto.go
	
c_itf2rpc:
	go build -o ./bin/itf2rpc.exe ./main/proto_tool/proto_lang/go/itf_to_rpc.go
	
itf2proto: c_itf2proto
	"./bin/itf2proto.exe" -i "./test/rpc_itf/" -o ./datas/proto/
	
itf2rpc:c_itf2rpc
	"./bin/itf2rpc.exe" -i "./test/rpc_itf/" -s -c -o "./rpc_nitf/" -gopath=$(GOPATH)/src
	
proto_compile: c_proto_compile
	"./bin/proto_compile.exe" -i ./datas/proto/ -o ./pb/ -ep "github.com/ProtossGenius/SureMoonNet"
	
go_protoread: c_proto_read
	"./bin/proto_read.exe" -proto "./datas/proto/" -pkgh "pb/" -o "./pbr/read.go" -gopath=$(GOPATH)/src

getlines:
	go run ./main/get-project-lines/get-pro-lines.go

importpkg:
	go get -u  github.com/json-iterator/go
	go get -u  github.com/robertkrimen/otto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/xtaci/kcp-go

clean:
	rm -f bin/*.exe
	rm -rf ./pb
	rm -rf ./rpc_nitf
	rm -rf ./pbr
test: itf2proto proto_compile go_protoread itf2rpc
	go run ./test/smn_net_rpc/test.go
