c_proto_read:
	go build -o ./bin/proto_read.exe ./src/com.suremoon.net/main/proto_tool/proto_read/MsgReader.go

c_proto_compile:
	go build -o ./bin/proto_compile.exe ./src/com.suremoon.net/main/proto_tool/proto_compile/proto_compile.go
	
c_itf2proto:
	go build -o ./bin/itf2proto.exe ./src/com.suremoon.net/main/proto_tool/proto_itf/to_proto/ItfToProto.go
	
c_itf2rpc:
	go build -o ./bin/itf2rpc.exe ./src/com.suremoon.net/main/proto_tool/proto_lange/go/itf_to_rpc.go
	
itf2proto: c_itf2proto
	"./bin/itf2proto.exe" -i "./src/rpc_itf/" -o ./datas/proto/
	
itf2rpc:c_itf2rpc
	"./bin/itf2rpc.exe" -i "./src/rpc_itf/" -s -c -o "./src/rpc_nitf/"
	
proto_compile: c_proto_compile
	"./bin/proto_compile.exe" -i ./datas/proto/ -o ./src/pb/ -p "./bin/protoc.exe"
	
go_protoread: itf2proto c_proto_read proto_compile
	"./bin/proto_read.exe" -proto "./datas/proto/" -pkgh "pb/" -o "./src/pbr/read.go"

getlines:
	go build -o ./bin/getlines.exe ./src/com.suremoon.net/main/get-project-lines/get-pro-lines.go 
	./bin/getlines.exe

importpkg:
	go get -u  github.com/json-iterator/go
	go get -u  github.com/robertkrimen/otto
