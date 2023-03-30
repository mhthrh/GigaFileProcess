package Server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mhthrh/GigaFileProcess/FileProcess"
	"github.com/mhthrh/GigaFileProcess/Utils/CryptoUtil"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"github.com/mhthrh/GigaFileProcess/entity"
	File "github.com/mhthrh/GigaFileProcess/file"
	"github.com/mhthrh/GigaFileProcess/ftp"
	Redis "github.com/mhthrh/GigaFileProcess/redis"
	"net/http"
	"time"
)

func Run(ctx *gin.Context) {
	var request entity.FileRequest
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot pars request object")
		ctx.Abort()
		return
	}
	if err := Validation.Priority(request.Priority); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	if err := Validation.Count(request.Count); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		ctx.Abort()
		return
	}
	if request.Name == "" {
		ctx.JSON(http.StatusBadRequest, "file name be empty")
		ctx.Abort()
		return
	}

	client, err := ftp.New("", "", "", 0)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "ftp calling error")
		ctx.Abort()
		return
	}
	err = client.Download("", "", "")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot find file")
		ctx.Abort()
		return
	}
	key := CryptoUtil.NewKey()
	key.FilePath = ""
	sha, err := key.Md5Sum()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot create sha")
		ctx.Abort()
		return
	}
	if request.Md5 != sha {
		ctx.JSON(http.StatusBadRequest, "sha is mismatch")
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, entity.FileResponse{
		ID:          request.ID,
		Status:      0,
		Description: "",
	})
	c := Redis.Client{Client: nil}
	byt, _ := json.Marshal(&request)
	_ = c.Set(request.Md5, string(byt))
	f, err := File.NewFile("path", "name")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "cannot find file")
		ctx.Abort()
		return
	}
	slice, err := f.Read()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot split file")
		ctx.Abort()
		return
	}
	err = FileProcess.New()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "cannot start rabbit/redis")
		ctx.Abort()
		return
	}
	go FileProcess.DoProcess(slice)
}
func Version(context *gin.Context) {
	context.JSON(http.StatusOK, "Ver:1.0.0")
}
func NotFound(context *gin.Context) {
	context.JSON(http.StatusOK, struct {
		Time        time.Time `json:"time"`
		Description string    `json:"description"`
	}{
		Time:        time.Now(),
		Description: "Workers are working, coming soon!!!",
	})
}
