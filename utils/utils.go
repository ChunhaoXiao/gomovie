package utils

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func StreamVideo(filename string, c *gin.Context) {

	filePath := filepath.Join("./uploads", filename)

	file, err := os.Open(filePath)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer file.Close()

	fmt.Println("first place")

	fileInfo, err := file.Stat()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("second place place")

	fileSize := fileInfo.Size()
	rangeHeader := c.GetHeader("Range")

	if rangeHeader == "" {
		c.Header("Content-Type", "video/mp4")
		c.Header("Content-Length", strconv.FormatInt(fileSize, 10))
		c.File(filePath)
		return
	}

	fmt.Println("third place")

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

func DownloadVideo(c *gin.Context) {
	id := c.Param("id")
	movieId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
	fmt.Println(movieId)

}

func MakeMovieThumb(file string, duration int16) []string {
	if !strings.Contains(file, "uploads/") {
		file = "uploads/" + file
	}

	var wg sync.WaitGroup
	fileListChan := make(chan string, 7)
	timePosition := duration / 7
	var i int16
	//wg.Add(7)
	for i = 1; i < 8; i++ {
		wg.Add(1)
		seconds := timePosition * i
		duration := time.Duration(seconds) * time.Second
		hours := int(duration.Hours())
		pos_minutes := int(duration.Minutes()) % 60
		pos_seconds := int(duration.Seconds()) % 60
		hmsString := fmt.Sprintf("%02d:%02d:%02d", hours, pos_minutes, pos_seconds)
		// go func(wg *sync.WaitGroup) {
		// 	defer wg.Done()
		// 	createThunmb(file, hmsString, wg)
		// }(&wg)
		go createThunmb(file, hmsString, &wg, fileListChan)
		//wg.Go(func(){createThunmb(file, hmsString, &wg)})
	}
	wg.Wait() // Wait for all worker goroutines to complete
	close(fileListChan)
	// go func() {
	// 	wg.Wait()           // Wait for all worker goroutines to complete
	// 	close(fileListChan) // Close the channel to signal no more results will be sent
	// }()
	//wg.Wait()
	var res []string
	for item := range fileListChan {
		fmt.Println("item###########", item)
		res = append(res, item)
	}

	return res

}

func MakeSingleMovieThumb(file string) string {
	if !strings.Contains(file, "uploads/") {
		file = "uploads/" + file
	}
	var wg sync.WaitGroup
	fileChan := make(chan string, 1)
	wg.Add(1)
	go createThunmb(file, "00:00:10", &wg, fileChan)
	wg.Wait() // Wait for all worker goroutines to complete
	close(fileChan)
	return <-fileChan
}

func createThunmb(file string, timepos string, wg *sync.WaitGroup, fileList chan string) {
	fmt.Println("timePos############", timepos)
	defer wg.Done()
	width := 500
	height := 500
	pictureFileName := uuid.New().String() + ".jpg"
	outputImage := "./thumb/" + pictureFileName

	cmd := exec.Command("ffmpeg", "-ss", timepos, "-t", "1", "-i", file,
		"-frames:v", "1",
		"-s", fmt.Sprintf("%dx%d", width, height),
		outputImage)
	cmd.Run()
	fileList <- pictureFileName
	fmt.Println("cmd completed..............")

}

func GetVfalidationMsg(err validator.ValidationErrors) []string {
	res := []string{}
	msgMap := map[string]string{
		"Username":             "用户名",
		"Password":             "密码",
		"OldPassword":          "原密码",
		"RePassword":           "确认密码",
		"PasswordConfirmation": "确认密码",
		"Email":                "邮箱",
		"required":             "不能为空",
		"min":                  "不能小于",
		"max":                  "不能大于",
		"email":                "格式不正确",
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

func Paginate(r *http.Request, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		q := r.URL.Query()
		page, _ := strconv.Atoi(q.Get("page"))
		if page <= 0 {
			page = 1
		}

		//pageSize, _ := strconv.Atoi(q.Get("page_size"))
		switch {
		case pageSize > 30:
			pageSize = 30
		case pageSize <= 0:
			pageSize = 30
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func GeneratePageStr(perpage int, currentPage int, totalCount int, displayPages int) []int {
	var pageList []int
	totalPages := totalCount / perpage

	start := 1
	end := displayPages

	if totalPages <= displayPages {
		for i := 1; i <= totalPages; i++ {
			pageList = append(pageList, i)
		}
	} else {

		if currentPage > displayPages/2 {
			start = currentPage - displayPages/2
			fmt.Println("start page::::", start)
			if currentPage+displayPages/2 > totalPages {
				end = totalPages
				start = start - (currentPage + displayPages/2 - totalPages)
			} else {
				end = currentPage + displayPages/2
			}
		}
		for i := start; i <= end; i++ {
			pageList = append(pageList, i)
		}

	}
	return pageList
}
