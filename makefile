proto_read:
	go build -o ./bin/proto_read.exe ./src/com.suremoon.net/main/proto_tool/proto_read/MsgReader.go

proto_compile:
	go build -o ./bin/proto_compile.exe ./src/com.suremoon.net/main/proto_tool/proto_compile/proto_compile.go
	
win_proto: proto_compile
	"./bin/proto_compile" -i ./datas/proto/ -o ./src/pb/ -p "./bin/protoc.exe"
	
win_go_protoread: proto_read win_proto
	"./bin/proto_read" -proto "./datas/proto/" -pkgh "pb/" -o "./src/pbr/read.go"

getlines:
	go build -o ./bin/getlines.exe ./src/com.suremoon.net/main/get-project-lines/get-pro-lines.go 
	./bin/getlines.exe

importpkg:
	go get -u  github.com/json-iterator/go
	go get -u  github.com/robertkrimen/otto
