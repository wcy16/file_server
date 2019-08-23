package middleware

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type FileInfo struct {
	FileName    string
	DisplayName string
	Expire      int64
	Message     string
}

var key []byte

func init() {
	file, err := os.Open("key")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	key = scanner.Bytes()
}

func CheckUploadPermission() gin.HandlerFunc {

	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()

	return func(c *gin.Context) {
		//params := c.Request.URL.Query()
		//msg := params.Get("token")
		msg := c.Request.Header.Get("token")

		token, _ := base64.StdEncoding.DecodeString(msg)

		if len(token) < nonceSize {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		nonce, ciphertext := token[:nonceSize], token[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var fileInfo FileInfo
		if err := json.Unmarshal(plaintext, &fileInfo); err == nil {
			if fileInfo.FileName == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, "empty file name")
				return
			}
			if fileInfo.Message != "upload" { // message should be "upload"
				c.AbortWithStatusJSON(http.StatusBadRequest, "error task")
				return
			}
			// validate time
			now := time.Now().Unix()
			if now > fileInfo.Expire {
				c.AbortWithStatusJSON(http.StatusBadRequest, "expired")
				return
			}

			if c.Keys == nil {
				c.Keys = make(map[string]interface{})
			}
			c.Keys["filename"] = fileInfo.FileName
			c.Keys["displayname"] = fileInfo.DisplayName

		} else {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Next()
	}
}

func CheckDownloadPermission() gin.HandlerFunc {

	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()

	return func(c *gin.Context) {
		params := c.Request.URL.Query()
		msg := params.Get("token")

		token, _ := base64.StdEncoding.DecodeString(msg)

		if len(token) < nonceSize {
			c.Redirect(http.StatusMovedPermanently, "/error")
			c.Abort()
			//c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		nonce, ciphertext := token[:nonceSize], token[nonceSize:]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			c.Redirect(http.StatusMovedPermanently, "/error")
			c.Abort()
			//c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		var fileInfo FileInfo
		if err := json.Unmarshal(plaintext, &fileInfo); err == nil {
			if fileInfo.FileName == "" {
				c.Redirect(http.StatusMovedPermanently, "/error")
				c.Abort()
				//c.AbortWithStatusJSON(http.StatusBadRequest, "empty file name")
				return
			}
			if fileInfo.Message != "download" { // message should be "download"
				c.Redirect(http.StatusMovedPermanently, "/error")
				c.Abort()
				//c.AbortWithStatusJSON(http.StatusBadRequest, "error task")
				return
			}
			// validate time
			now := time.Now().Unix()
			if now > fileInfo.Expire {
				c.Redirect(http.StatusMovedPermanently, "/error")
				c.Abort()
				//c.AbortWithStatusJSON(http.StatusBadRequest, "expired")
				return
			}

			if c.Keys == nil {
				c.Keys = make(map[string]interface{})
			}
			c.Keys["filename"] = fileInfo.FileName
			c.Keys["displayname"] = fileInfo.DisplayName

		} else {
			c.Redirect(http.StatusMovedPermanently, "/error")
			c.Abort()
			//c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.Next()
	}
}
