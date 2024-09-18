package controller

import (
	"fmt"
	"git-practice-api/go-gin-chat/services/img_upload_connector"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

func ImgKrUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filepath := viper.GetString(`app.upload_file_path`)
	//指定文件或目录的元信息（例如大小、修改时间、权限等）
	if _, err := os.Stat(filepath); err != nil {
		if !os.IsExist(err) {
			os.MkdirAll(filepath, os.ModePerm)
		}
	}

	filename := filepath + file.Filename

	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	krUpload := img_upload_connector.ImgCreate().Upload(filename)

	// 删除临时图片
	os.Remove(filename)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": map[string]interface{}{
			"url": krUpload,
		},
	})
}
