#include<iostream>
#include <boost/asio.hpp>

#include "smncpp/socket_mtd.h"
#include "smncpp/base_asio_socket.h"
#include "pb/rip_rpc_itf.pb.h"
#include "pb/smn_base.pb.h"
#include "pb/smn_dict.pb.h"
#include "smn_itf/rpc_itf.Login.h"

using namespace std;
using namespace boost::asio;

class Login :public rpc_itf::Login{
typedef boost::asio::ip::tcp::socket socket;
public:
	Login(socket& s):_c(s){}
private:
	smnet::SMConn _c;

public:
	rip_rpc_itf::Login_DoLogin_Ret DoLogin(const std::string& user, const std::string& pswd, int64_t code);
	std::vector<int64_t> Test1(const std::vector<std::string>& a, const std::vector<int64_t>& b, const std::vector<uint64_t>& c, const std::vector<uint64_t>& d, const std::vector<int32_t>& e);
	bool Test2(const std::string& key, const smnet::Conn& c);
};


rip_rpc_itf::Login_DoLogin_Ret Login::DoLogin(const std::string& user, const std::string& pswd, int64_t code){
	smn_base::Call call;	
	rip_rpc_itf::Login_DoLogin_Prm prm;
	prm.set_user(user);
	prm.set_pswd(pswd);
	prm.set_code(code);
	call.set_dict(smn_dict::rip_rpc_itf_Login_DoLogin_Prm);
	call.set_msg(prm.SerializeAsString());
	auto result = smnet::writeString(this->_c, call.SerializeAsString());

	if (result != smnet::ConnStatusSucc){
		throw this->_c.lastError();
	}

	smnet::Bytes retBuff;
	smnet::readLenBytes(this->_c, retBuff);
	smn_base::Ret ret;
	rip_rpc_itf::Login_DoLogin_Ret lret;	
	ret.ParseFromArray(retBuff.arr, retBuff.size());
	lret.ParseFromString(ret.msg());
	return lret;
	
}
std::vector<int64_t> Login::Test1(const std::vector<std::string>& a, const std::vector<int64_t>& b, const std::vector<uint64_t>& c, const std::vector<uint64_t>& d, const std::vector<int32_t>& e){
	rip_rpc_itf::Login_Test1_Prm prm;
	prm.add_a();
	std::vector<int64_t> ret;
	return ret;	
}

bool Login::Test2(const std::string& key, const smnet::Conn& c){
	return false;
}


int main(){
	io_service ioc;
	ip::tcp::socket socket(ioc);
	ip::tcp::endpoint ep(ip::address_v4::from_string("127.0.0.1"), 7000);
	boost::system::error_code ec;
	socket.connect(ep, ec);
	if (ec){
		cout << boost::system::system_error(ec).what() <<endl;
		return -1;
	}
	Login login(socket);
	auto result = login.DoLogin("hello_user", "hello_pwsd", -123);
	
	cout <<result.DebugString() <<endl;
}
