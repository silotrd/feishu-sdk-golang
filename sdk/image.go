package sdk

import (
	"bytes"
	"fmt"
	"github.com/galaxy-book/feishu-sdk-golang/core/consts"
	http2 "github.com/galaxy-book/feishu-sdk-golang/core/util/http"
	"github.com/galaxy-book/feishu-sdk-golang/core/util/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"unsafe"
)

func (t Tenant) NewFileUploadRequest(uri string, params map[string]string, paramName, path string) (imgKey string ,err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return
	}
	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	writer.WriteField("image_type", "message")
	err = writer.Close()
	if err != nil {
		return
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", "Bearer "+t.TenantAccessToken)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	str := (*string)(unsafe.Pointer(&respBytes))
//	fmt.Println(*str)
	return *str, nil
}

//func (t Tenant) GetImage(imageKey string) ([]byte, error) {
func (t Tenant) GetImage(imageKey string, isApp bool) (io.ReadCloser, error) {
	queryParams := map[string]interface{}{}
	queryParams["image_key"] = imageKey
	if isApp {
		queryParams["operator"] = "app"
	}
	request, err := http.NewRequest("GET", consts.ApiGetImage+http2.ConvertToQueryParams(queryParams), nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+t.TenantAccessToken)
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//respBytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}

	return resp.Body, nil
}
