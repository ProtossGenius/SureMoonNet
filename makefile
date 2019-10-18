c_proto_read:
	go build -o ./bin/proto_read.exe ./src/com.suremoon.net/main/proto_tool/proto_read/MsgReader.go

c_proto_compile:
	go build -o ./bin/proto_compile.exe ./src/com.suremoon.net/main/proto_tool/proto_compile/proto_compile.go
	
c_itf2proto:
	go build -o ./bin/itf2proto.exe ./src/com.suremoon.net/main/proto_tool/proto_itf/to_proto/ItfToProto.go
	
c_itf2rpc:
	go build -o ./bin/itf2rpc.exe ./src/com.suremoon.net/main/proto_tool/proto_lang/go/itf_to_rpc.go
	
itf2proto: c_itf2proto
	"./bin/itf2proto.exe" -i "./src/rpc_itf/" -o ./datas/proto/
	
itf2rpc:c_itf2rpc
	"./bin/itf2rpc.exe" -i "./src/rpc_itf/" -s -c -o "./src/rpc_nitf/"
	
proto_compile: c_proto_compile
	"./bin/proto_compile.exe" -i ./datas/proto/ -o ./src/pb/
	
go_protoread: c_proto_read
	"./bin/proto_read.exe" -proto "./datas/proto/" -pkgh "pb/" -o "./src/pbr/read.go"

getlines:
	go build -o ./bin/getlines.exe ./src/com.suremoon.net/main/get-project-lines/get-pro-lines.go 
	./bin/getlines.exe

importpkg:
	go get -u  github.com/json-iterator/go
	go get -u  github.com/robertkrimen/otto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/xtaci/kcp-go

clean:
	rm -f bin/*.exe
	rm -rf src/pb
	rm -rf src/rpc_nitf
	rm -rf src/pbr
test: itf2proto go_protoread proto_compile itf2rpc
	go run ./src/com.suremoon.net/test/smn_net_rpc/test.go
