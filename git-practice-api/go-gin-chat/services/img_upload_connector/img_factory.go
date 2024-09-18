package img_upload_connector

import (
	"git-practice-api/go-gin-chat/services"
	"git-practice-api/go-gin-chat/services/img_freeimage"
)

// 定义 serve 的映射关系
var serveMap = map[string]services.ImgUploadInterface{
	"fi": &img_freeimage.ImgFreeImageService{},
}

func ImgCreate() services.ImgUploadInterface {
	return serveMap["fi"]
}
