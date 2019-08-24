package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

const path = "data/"

func prefix(name string) string {
	return path + name
}

func Upload(c *gin.Context) {
	name := c.Keys["displayname"].(string)
	file, err := c.FormFile(name)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(c.Keys["filename"].(string))
	if err := c.SaveUploadedFile(file, prefix(filename)); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}
}

func Download(c *gin.Context) {
	filename := c.Param("filename")
	displayName := c.Keys["displayname"]

	if filename != c.Keys["filename"] {
		c.Redirect(http.StatusTemporaryRedirect, "/not_found")
		return
	}

	targetPath := prefix(filename)
	//log.Println(targetPath)
	//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), path) {
		c.Redirect(http.StatusTemporaryRedirect, "/not_found")
		return
	}

	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))
	c.File(targetPath)
}
