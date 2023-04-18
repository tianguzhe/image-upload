package wrapper

import (
	"fmt"
	"image-upload/safety"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func OssUpload(fileName string) (string, error) {
	// 创建OSSClient实例。
	// yourEndpoint填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	client, err := oss.New(safety.Endpoint, safety.AccessKeyID, safety.AccessKeySecret)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 填写存储空间名称，例如examplebucket。
	bucket, err := client.Bucket(safety.BucketName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fileNameSplit := strings.Split(fileName, "/")

	finalAddr := fmt.Sprintf("%s%s", safety.PathPrefix, fileNameSplit[len(fileNameSplit)-1])

	// 依次填写Object的完整路径（例如exampledir/exampleobject.txt）和本地文件的完整路径（例如D:\\localpath\\examplefile.txt）。
	err = bucket.PutObjectFromFile(finalAddr, fileName)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	return finalAddr, err
}
