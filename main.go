package main

import (
	"flag"
	"fmt"
	"image-upload/dir"
	"image-upload/safety"
	"image-upload/wrapper"
	"io"
	"os"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
	"golang.org/x/exp/slog"
)

type UploadResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		ImageURL    string `json:"image_url"`
		ImageWidth  int    `json:"image_width"`
		ImageHeight int    `json:"image_height"`
	} `json:"data"`
}

func init() {
	logFile, err := os.OpenFile(safety.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		fmt.Printf("create file error %s \n\n", err.Error())
		panic(err)
	}

	opts := slog.HandlerOptions{
		AddSource: true,
	}

	jsonHandler := opts.NewJSONHandler(io.MultiWriter(logFile)).WithAttrs(
		[]slog.Attr{slog.String("app-version", "v0.0.1")},
	)

	slog.SetDefault(slog.New(jsonHandler))
}

func main() {

	var wg sync.WaitGroup

	allFile := make([]string, 0)

	var file string
	var fileDirPath string

	flag.StringVar(&file, "f", "", "上传的图片地址")
	flag.StringVar(&fileDirPath, "d", "", "上传的图片目录")

	flag.Parse()

	if file == "" && fileDirPath == "" {
		panic("是否遗漏 -f flag")
	}

	if file != "" {
		allFile = append(allFile, file)
	}

	if fileDirPath != "" {
		allFile = dir.GetDirFile(fileDirPath)
	}

	for _, fileName := range allFile {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()

			respBody, err := wrapper.BiliUpload(fileName)
			if err != nil {
				panic(err)
			}

			var bili UploadResult

			sonic.Unmarshal(respBody, &bili)

			slog.Info("upload",
				slog.String("source", fileName),
				slog.Group("bili",
					slog.String("link", bili.Data.ImageURL),
					slog.Int("width", bili.Data.ImageWidth),
					slog.Int("height", bili.Data.ImageHeight),
				),
			)

			color.Red("%s = %s \n", fileName, bili.Data.ImageURL)
		}(fileName)

		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			respBody, err := wrapper.OssUpload(fileName)
			if err != nil {
				panic(err)
			}

			slog.Info("upload",
				slog.String("source", fileName),
				slog.Group("oss",
					slog.String("link", fmt.Sprintf("https://cdn.yryz2.com/%s", respBody)),
				),
			)

			color.Blue("%s = %s \n", fileName, fmt.Sprintf("https://cdn.yryz2.com/%s", respBody))
		}(fileName)
	}

	wg.Wait()
}
