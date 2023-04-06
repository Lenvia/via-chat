package conf

/*参数说明
app.port // 应用端口
app.upload_file_path // 图片上传的临时文件夹目录，绝对路径！
app.cookie_key // 生成加密session
app.serve_type // 默认请使用GoServe
mysql.dsn // mysql 连接地址dsn
app.debug_mod // 开发模式建议设置为`true` 避免修改静态资源需要重启服务
*/

var AppJsonConfig = []byte(`
{
  "app": {
    "port": "8322",
    "upload_file_path": "/Users/yy/GithubProjects/via-chat/tmp_images/",
    "cookie_key": "4238uihfieh49r3453kjdfg",
    "serve_type": "GoServe",
    "debug_mod": "true"
  },
  "mysql": {
    "dsn": "root:y724290941@tcp(172.17.0.2:3306)/room?charset=utf8mb4&parseTime=True&loc=Local"
  }
}
`) // 如果要在本地运行，需要把172.17.0.2（docker中的mysql的IP）改成127.0.0.1
