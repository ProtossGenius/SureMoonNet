c_proto_read:
	cd ./main/proto_tool/proto_read/smn_pr_go &&  go install

c_proto_compile:
	cd ./main/proto_tool/smn_protocpl && go install 
	
c_itf2proto:
	cd ./main/proto_tool/proto_itf/smn_itf2proto  && go install
	
c_itf2rpc_go:
	cd ./main/proto_tool/proto_lang/smn_itf2rpc_go	&& go install 

itf2proto: c_itf2proto
	smn_itf2proto -i "./test/rpc_itf/" -o ./datas/proto/
	
install: c_proto_read c_proto_compile c_itf2proto c_itf2rpc_go smgit
	echo "finish"

itf2rpc:c_itf2rpc_go
	smn_itf2rpc_go -i "./test/rpc_itf/" -s -c -o "./rpc_nitf/" -gopath=$(GOPATH)/src
	
proto_compile: c_proto_compile
	smn_protocpl -i ./datas/proto/ -o ./pb/ -ep "github.com/ProtossGenius/SureMoonNet"
	
go_protoread: c_proto_read
	smn_pr_go -proto "./datas/proto/" -pkgh "pb/" -o "./pbr/read.go" -gopath=$(GOPATH)/src -ext="/github.com/ProtossGenius/SureMoonNet"

smgit:
	cd ./main/smn_tool/smgit && go install

getlines:
	go run ./main/get-project-lines/get-pro-lines.go

importpkg:
	go get -u  github.com/json-iterator/go
	go get -u  github.com/robertkrimen/otto
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u github.com/xtaci/kcp-go

clean:
	rm -f datas/proto/rip_rpc_itf.proto
	rm -f datas/proto/smn_dict.proto
	rm -f bin/*.exe
	rm -rf ./rpc_nitf
	rm -rf ./pbr
	rm -rf ./pb/rip_rpc_itf ./pb/smn_dict
test: clean itf2proto proto_compile go_protoread itf2rpc
	go run ./test/smn_net_rpc/test.go
