package result

// Codes:定义状态
type Codes struct {
	Message          map[uint]string
	Success          uint
	Failure          uint
	PasswordError    uint
	SystemError      uint
	ImgKrUploadError uint
	LoadFileError    uint
	FileWirteError   uint
	OpenFileError    uint
	CopyFileError    uint
	CmdStartError    uint
	DialError        uint
	ReadError        uint
}

// 状态码
var APIcode = &Codes{
	Success:          200,
	Failure:          501,
	PasswordError:    600,
	SystemError:      601,
	ImgKrUploadError: 602,
	LoadFileError:    603,
	FileWirteError:   604,
	OpenFileError:    605,
	CopyFileError:    606,
	CmdStartError:    607,
	DialError:        608,
	ReadError:        609,
}

// 状态信息初始化
func init() {
	APIcode.Message = map[uint]string{
		APIcode.Success:          "成功",
		APIcode.Failure:          "失败",
		APIcode.PasswordError:    "密码错误，请重新输入",
		APIcode.SystemError:      "系统错误",
		APIcode.ImgKrUploadError: "图片上传错误",
		APIcode.LoadFileError:    "上传文件错误",
		APIcode.FileWirteError:   "写文件错误",
		APIcode.OpenFileError:    "打开文件错误",
		APIcode.CopyFileError:    "复制文件失败",
		APIcode.CmdStartError:    "开启cmd失败",
		APIcode.DialError:        "建立websocket连接失败",
		APIcode.ReadError:        "读取失败",
	}
}

// GetMessage 供外部调用
func (c *Codes) GetMessage(code uint) string {
	message, ok := c.Message[code]
	if !ok {
		return ""
	}
	return message
}
