package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sku-manage/mixin"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// 修改成可配置
const host = "http://127.0.0.1:8500/v1/inspect/file/"
const PICTUREDIR = "./file/"

var fs = http.FileServer(http.Dir(PICTUREDIR))

func (this *AccountService) get_file(w http.ResponseWriter, r *http.Request) {
	http.StripPrefix("/v1/inspect/file", fs).ServeHTTP(w, r)
}

// 上传单个文件
func (this *AccountService) upload_file_handle(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, head, err := r.FormFile("file")
	if err != nil {
		logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}
	defer file.Close()

	fileNameSli := strings.Split(head.Filename, ".")
	if len(fileNameSli) < 2 {
		logrus.Debugf("[AccountService.upload_file_handle] error file name")
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}

	uuid, err := generalUuid()
	if err != nil {
		logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}
	newName := uuid + "." + fileNameSli[len(fileNameSli)-1]

	//创建文件
	fW, err := os.Create(PICTUREDIR + newName)
	if err != nil {
		logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}
	defer fW.Close()
	//写内容
	_, err = io.Copy(fW, file)
	if err != nil {
		logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
		this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
		return
	}

	fmt.Fprint(w, host+newName)
}

// 上传多个文件
func (this *AccountService) upload_files_handle(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	mp := r.MultipartForm

	if mp == nil {
		log.Println("not MultipartForm.")
		w.Write(([]byte)("不是MultipartForm格式"))
		return
	}

	fileHeaders, findFile := mp.File["file"]
	if !findFile || len(fileHeaders) == 0 {
		log.Println("file count == 0.")
		w.Write(([]byte)("没有上传文件"))
		return
	}
	logrus.Debugf("上传了%d个文件", len(fileHeaders))

	var resp []string
	for _, v := range fileHeaders {
		fileNameSli := strings.Split(v.Filename, ".")
		if len(fileNameSli) < 2 {
			logrus.Debugf("[AccountService.upload_file_handle] error file name")
			this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
			return
		}

		file, err := v.Open()
		if err != nil {
			logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
			this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
			return
		}
		defer file.Close()

		uuid, err := generalUuid()
		if err != nil {
			logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
			this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
			return
		}
		newName := uuid + "." + fileNameSli[len(fileNameSli)-1]

		//创建文件
		fW, err := os.Create(PICTUREDIR + newName)
		if err != nil {
			logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
			this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
			return
		}
		defer fW.Close()
		//写内容
		_, err = io.Copy(fW, file)
		if err != nil {
			logrus.Debugf("[AccountService.upload_file_handle] error %s", err.Error())
			this.ResponseErrCode(w, mixin.ErrorServerUnKnow)
			return
		}

		resp = append(resp, host+newName)
	}

	this.ResponseOK(w, map[string][]string{
		"url": resp,
	})
}

type fileState struct {
	FileName   string `json:"file_name"`
	UploadPath string `json:"upload_path"`
}

//使用nginx上传
func (this *AccountService) file_upload_handle(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)

	logrus.Debugf("[file_upload_handle] %v", r.MultipartForm)

	file := fileState{
		FileName:   r.MultipartForm.Value["file.name"][0],
		UploadPath: r.MultipartForm.Value["file.path"][0],
	}
	path, err := saveFile(file)
	if err != nil {
		log.Println(err.Error())
		fmt.Fprintln(w, err.Error())
		return
	}

	fmt.Fprint(w, path)
}

func saveFile(file fileState) (string, error) {

	fileNameSli := strings.Split(file.FileName, ".")
	if len(fileNameSli) < 2 {
		return "", errors.New("file format error")
	}

	uuid, err := generalUuid()
	newName := uuid + "." + fileNameSli[len(fileNameSli)-1]
	if err != nil {
		return "", err
	}

	newPath := PICTUREDIR + newName
	err = os.Rename(file.UploadPath, newPath)
	if err != nil {
		return "", err
	}
	return host + newName, nil
}

func generalUuid() (string, error) {
	unix32bits := uint32(time.Now().UTC().Unix())

	buff := make([]byte, 12)
	numRead, err := rand.Read(buff)
	if numRead != len(buff) || err != nil {
		return "", err
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x-%x", unix32bits, buff[0:2], buff[2:4], buff[4:6], buff[6:8], buff[8:]), nil
}
