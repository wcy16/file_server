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
	// return file name
	c.String(http.StatusOK, fmt.Sprint(file.Filename))
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
	// check filename
	if !strings.HasPrefix(filepath.Clean(targetPath), path) {
		c.Redirect(http.StatusTemporaryRedirect, "/not_found")
		return
	}

	c.Writer.Header().Set("content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", displayName))
	c.File(targetPath)
}
