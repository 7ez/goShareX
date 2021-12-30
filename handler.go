package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

var allowableFileTypes = []string{".png", ".jpg", ".jpeg", ".gif", ".mp4", ".webm", ".mov", ".avi"}

const characters = "ko8LXtAKFU039ZDeGnOl2yENC7fTxVQPcshJ65dmBMapzSquHj4gbRrwi1WIYv"

func GenFileName() string {
	b := make([]byte, 16)
	for i := range b {
		b[i] = characters[rand.Intn(16)]
	}
	return string(b)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func canBeUploaded(fileExt string) bool {
	for _, ext := range allowableFileTypes {
		if ext == fileExt {
			return true
		}
	}
	return false
}

func getFile(c *gin.Context) {
	file := fmt.Sprintf("./files/%s", c.Param("file"))
	if fileExists(file) {
		c.File(file)
		return
	}

	c.String(http.StatusNotFound, "File not found")
	return
}

func uploadFile(c *gin.Context) {
	if c.GetHeader("k") != Config.UploadKey {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "file upload err: invalid key",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  fmt.Sprintf("get file err: %s", err.Error()),
		})
		return
	}

	filename := filepath.Base(file.Filename)
	fileext := filepath.Ext(filename)

	if !canBeUploaded(fileext) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "file upload err: unsupported filetype",
		})
		return
	}

	filename = fmt.Sprintf("%s%s", GenFileName(), fileext)

	dst := fmt.Sprintf("./files/%s", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  fmt.Sprintf("file upload err: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"url":    fmt.Sprintf("%s/i/%s", Config.Domain, filename),
	})
	return
}

func genShareXConf(c *gin.Context) {
	c.Header("Content-Disposition", "attachment; filename=config.sxcu")
	c.Header("Content-Type", "application/json; charset=utf-8")

	c.JSON(http.StatusOK, gin.H{
		"Version":         "13.5.0",
		"Name":            Config.AppName,
		"DestinationType": "ImageUploader",
		"RequestMethod":   "POST",
		"RequestURL":      fmt.Sprintf("%s/i/upload", Config.Domain),
		"Body":            "MultipartFormData",
		"Headers": gin.H{
			"k": Config.UploadKey,
		},
		"FileFormName": "file",
		"URL":          "$json:url$",
		"ErrorMessage": "$json:error$",
	})
}
