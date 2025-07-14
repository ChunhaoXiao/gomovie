package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func StreamVideo(filename string, c *gin.Context) {

	filePath := filepath.Join("./uploads", filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fileSize := fileInfo.Size()
	rangeHeader := c.GetHeader("Range")

	if rangeHeader == "" {
		c.Header("Content-Type", "video/mp4")
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		c.File(filePath)
		return
	}

	parts := strings.Split(rangeHeader, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	rangeValues := strings.Split(parts[1], "-")
	startByte, err := strconv.ParseInt(rangeValues[0], 10, 64)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var endByte int64
	if len(rangeValues) > 1 && rangeValues[1] != "" {
		endByte, err = strconv.ParseInt(rangeValues[1], 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	} else {
		endByte = fileSize - 1
	}

	contentLength := endByte - startByte + 1
	c.Header("Content-Range", "bytes "+strconv.FormatInt(startByte, 10)+"-"+strconv.FormatInt(endByte, 10)+"/"+strconv.FormatInt(fileSize, 10))
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Header("Content-Type", "video/mp4")
	c.Header("Accept-Ranges", "bytes")
	c.Status(http.StatusPartialContent)

	_, err = file.Seek(startByte, 0)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	buffer := make([]byte, contentLength)
	_, err = file.Read(buffer)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Writer.Write(buffer)

}

func MakeMovieThumb(file string) string {
	time := "00:00:03" // Example: 5 seconds
	width := 200
	height := 300
	if !strings.Contains(file, "uploads/") {
		file = "uploads/" + file
	}

	pictureFileName := uuid.New().String() + ".jpg"
	outputImage := "thumb/" + pictureFileName
	fmt.Println("outputPage::::", outputImage)
	var stderr bytes.Buffer

	cmd := exec.Command("ffmpeg", "-i", file,
		"-ss", time,
		"-frames:v", "1",
		"-s", fmt.Sprintf("%dx%d", width, height),
		outputImage)
	cmd.Stderr = &stderr
	err := cmd.Run()
	fmt.Println(fmt.Sprint(err)+":::::::", stderr.String())
	if err != nil {
		//panic("could not generate frame")
		return ""
	}
	return pictureFileName
}

func GetVfalidationMsg(err validator.ValidationErrors) []string {
	res := []string{}
	msgMap := map[string]string{
		"Username":   "用户名",
		"Password":   "密码",
		"RePassword": "确认密码",
		"Email":      "邮箱",
		"required":   "不能为空",
		"min":        "不能小于",
		"max":        "不能大于",
		"email":      "格式不正确",
	}

	for _, e := range err {
		msg := msgMap[e.Field()] + msgMap[e.Tag()] + e.Param()
		fmt.Println("field is:", e.Field())
		fmt.Println("tag is:", e.Tag())
		// fmt.Println("params###", e.Param())
		res = append(res, msg)
	}

	return res
}
