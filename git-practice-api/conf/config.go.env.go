package conf

var AppJsonConfig = []byte(`
{
  "app": {
    "port": "8322",
    "upload_file_path": "G:\\golang_project\\chatroom\\go-gin-chat\\tempt_img",
    "cookie_key": "4238uihfieh49r3453kjdfg",
    "serve_type": "GoServe",
    "debug_mod": "true"
  },
  "mysql": {
    "dsn": "root:123456@tcp(127.0.0.1:3306)/go_gin_chat?charset=utf8mb4&parseTime=True&loc=Local"
  }
}
`)
