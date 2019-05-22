package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("服务器启动失败", err.Error())
		return
	}

}
func upload(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32 << 20)
	//接收客户端传来的文件 uploadfile 与客户端保持一致
	file, handler, err := request.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	//上传的文件保存在ppp路径下
	ext := path.Ext(handler.Filename) //获取文件后缀
	fileNewName := string(time.Now().Format("20060102150405")) + strconv.Itoa(time.Now().Nanosecond()) + ext

	f, err := os.OpenFile("./datas/fileweb/"+fileNewName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	fmt.Fprintln(writer, "upload ok!"+fileNewName)
}

func index(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(tpl))
}

const tpl = `<html>
<head>
<title>上传文件</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
<input type="file" name="uploadfile">
<input type="hidden" name="token" value="{...{.}...}">
<input type="submit" value="upload">
</form>
</body>
</html>
`
