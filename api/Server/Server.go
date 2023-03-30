package Server

import (
	"github.com/gin-gonic/gin"
	"github.com/mhthrh/GigaFileProcess/FileProcess"
	"github.com/mhthrh/GigaFileProcess/Utils/CryptoUtil"
	"github.com/mhthrh/GigaFileProcess/Validation"
	"github.com/mhthrh/GigaFileProcess/entity"
	File "github.com/mhthrh/GigaFileProcess/file"
	"github.com/mhthrh/GigaFileProcess/ftp"
	"net/http"
	"path/filepath"
	"time"
)

const (
	path = "api/files"
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

	client, err := ftp.New("localhost", "ftpuser", "123456", 21)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "ftp calling error")
		ctx.Abort()
		return
	}
	err = client.Download(request.Name, request.Name, path)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "cannot find file")
		ctx.Abort()
		return
	}
	key := CryptoUtil.NewKey()
	key.FilePath = filepath.Join(path, request.Name)
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

	f, err := File.NewFile(path, request.Name)
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
	p, err := FileProcess.New()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "cannot start rabbit/redis")
		ctx.Abort()
		return
	}
	go p.DoProcess(slice)
	ctx.JSON(http.StatusOK, entity.FileResponse{
		ID:          request.ID,
		Status:      0,
		Description: "",
	})
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
