package wrapper

import (
	"bytes"
	"fmt"
	"image-upload/safety"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func BiliUpload(filename string) ([]byte, error) {

	client := &http.Client{}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("biz", "new_dyn")
	bodyWriter.WriteField("category", "daily")
	bodyWriter.WriteField("csrf", safety.Csrf)

	fileWriter, err := bodyWriter.CreateFormFile("file_up", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
	}

	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		fmt.Println("error copy ")
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, _ := http.NewRequest("POST", "https://api.bilibili.com/x/dynamic/feed/draw/upload_bfs", bodyBuf)

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", "api.bilibili.com")
	req.Header.Set("Origin", "https://t.bilibili.com")
	req.Header.Set("Referer", "https://t.bilibili.com/")

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/113.0")

	req.Header.Set("Cookie", safety.Cookie)

	resp, err := client.Do(req)

	fmt.Println(bodyBuf)
	if err != nil {
		fmt.Println("error copy ")
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)

	return respBody, err
}
