c_proto_compile:
	cd ./main/proto_tool/smn_protocpl && go install 
	
c_itf2proto:
	cd ./main/proto_tool/proto_itf/smn_itf2proto  && go install
	
c_itf2rpc_go:
	cd ./main/proto_tool/proto_lang/smn_itf2rpc_go	&& go install 

c_goitf2lang:
	cd ./main/proto_tool/smn_goitf2lang && go install

goitf2lang: c_goitf2lang
	smn_goitf2lang -lang=cpp -i="./test/rpc_itfs/" -o="./cpp_itf"
	
itf2proto: c_itf2proto
	smn_itf2proto -i "./test/rpc_itfs/" -o ./datas/proto/
	
install: c_proto_compile c_itf2proto c_itf2rpc_go smgit smlines c_goitf2lang 
	echo "finish"

itf2rpc:c_itf2rpc_go
	smn_itf2rpc_go -i "./test/rpc_itfs/" -s -c -o "./rpc_nitf/" -module="github.com/ProtossGenius/SureMoonNet"
	
proto_compile: c_proto_compile
	smn_protocpl -i ./datas/proto/ -o ./pb/ -gm "github.com/ProtossGenius/SureMoonNet" -lang=go

smgit:
	cd ./main/smn_tool/smgit && go install

smlines:
	cd ./main/smlines && go install

clean:
	rm -f datas/proto/rip_rpc_itf.proto
	rm -f datas/proto/rip_ano_rpc_itf.proto
	rm -f datas/proto/smn_dict.proto
	rm -f bin/*.exe
	rm -rf ./javapb
	rm -rf ./cpppb
	rm -rf ./cpp_itf
	rm -rf ./java_itf
	rm -rf ./rpc_nitf
	rm -rf ./pbr
	rm -rf ./pb/rip_rpc_itf
	rm -rf ./pb/rip_ano_rpc_itf
	rm -rf ./pb/smn_dict
	rm -rf datas/proto/temp
	
test: itf2proto proto_compile itf2rpc
	go run ./test/smn_net_rpc/test.go

qrun:
	go run ./test/smn_net_rpc/test.go

nothing:
