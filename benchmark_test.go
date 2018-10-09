package gotten_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
)

type (
	ChangeParam struct {
		NewToken   string `type:"form"`
		OldToken   string `type:"form" require:"true"`
		SecureId   string `type:"form" require:"true"`
		TokenSec   string `type:"form"`
		OldSec     string `type:"form"`
		Expiration int    `type:"form" require:"true"`
	}

	BenchService struct {
		Change func(param *ChangeParam) (*http.Request, error) `method:"POST" path:"change_item"`
		Upload func(param *UploadParam) (*http.Request, error) `method:"POST" path:"add_item"`
	}

	GitService struct {
		UpdateDeployKey func(param *UpdateParam) (*http.Request, error) `method:"PUT" path:"/projects/{id}/deploy_keys/{key_id}"`
	}

	UploadParam struct {
		PhpSession string          `type:"part" key:"PHP_SESSION_UPLOAD_PROGRESS" default:"qscbox"`
		Filecount  int             `type:"part" default:"1"`
		File       gotten.FilePath `type:"part" require:"true"`
		Callback   string          `type:"part" default:"handleUploadCallback"`
		IsIe9      int             `type:"part" default:"0"`
	}

	UpdateParam struct {
		Id    string `type:"path"`
		KeyId string `type:"path"`
		Key   *Key   `type:"json" require:"true"`
	}

	Key struct {
		Title   string `json:"title"`
		Key     string `json:"key,omitempty"`
		CanPush bool   `json:"can_push"`
	}
)

var (
	service     = new(BenchService)
	gitService  = new(GitService)
	changeParam = &ChangeParam{
		OldToken:   "test",
		SecureId:   "8nx1391907c5971n9112321d9y",
		Expiration: 86400,
	}

	uploadParam = &UploadParam{
		File: gotten.FilePath("testAssets/avatar.jpg"),
	}

	updateParam = &UpdateParam{
		Id:    "12",
		KeyId: "1234",
		Key: &Key{
			Title:   "Push Key",
			CanPush: true,
		},
	}
)

func init() {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://box.zjuqsc.com/item/").
		Build()
	if err != nil {
		panic(err)
	}
	err = creator.Impl(service)
	if err != nil {
		panic(err)
	}

	gitCreator, err := gotten.NewBuilder().
		SetBaseUrl("https://git.zjuqsc.com/").
		Build()
	if err != nil {
		panic(err)
	}
	err = gitCreator.Impl(gitService)
	if err != nil {
		panic(err)
	}
}

func buildFormRequestTraditionally(param *ChangeParam) (req *http.Request, err error) {
	form := make(url.Values)
	form.Add("new_token", param.NewToken)
	form.Add("old_token", param.OldToken)
	form.Add("secure_id", param.SecureId)
	form.Add("token_sec", param.TokenSec)
	form.Add("old_sec", param.OldSec)
	form.Add("new_token", param.NewToken)
	body := bytes.NewBufferString(form.Encode())
	req, err = http.NewRequest("POST", "https://box.zjuqsc.com/item/change_item", body)
	if err == nil {
		req.Header.Set(headers.HeaderContentType, headers.MIMEApplicationForm)
	}
	return
}

func buildMultipartRequestTraditionally(param *UploadParam) (req *http.Request, err error) {
	var partWriter io.Writer
	body := bytes.NewBufferString("")
	writer := multipart.NewWriter(body)
	if param.PhpSession == gotten.ZeroStr {
		param.PhpSession = "qscbox"
	}
	if param.Filecount == gotten.ZeroInt {
		param.Filecount = 1
	}
	if param.Callback == gotten.ZeroStr {
		param.Callback = "handleUploadCallback"
	}

	writer.WriteField("PHP_SESSION_UPLOAD_PROGRESS", param.PhpSession)
	writer.WriteField("filecount", strconv.Itoa(param.Filecount))
	writer.WriteField("callback", param.Callback)
	writer.WriteField("is_ie9", strconv.Itoa(param.IsIe9))
	var file *os.File
	file, err = os.Open(string(param.File))
	if err == nil {
		if partWriter, err = writer.CreateFormFile("file", "avatar.jpg"); err == nil {
			_, err = io.Copy(partWriter, file)
		}
	}
	file.Close()
	writer.Close()
	if err == nil {
		req, err = http.NewRequest("POST", "https://box.zjuqsc.com/item/add_item", body)
		if err == nil {
			req.Header.Set(headers.HeaderContentType, headers.MIMEMultipartForm)
		}
	}
	return
}

func buildJSONRequestTraditionally(param *UpdateParam) (req *http.Request, err error) {
	if param.Id == gotten.ZeroStr || param.KeyId == gotten.ZeroStr || param.Key == nil {
		err = errors.New("param is invalid")
	}

	if err == nil {
		target := fmt.Sprintf("https://git.zjuqsc.com/projects/%s/deploy_keys/%s", url.QueryEscape(param.Id), url.QueryEscape(param.KeyId))
		var data []byte
		data, err = json.Marshal(param.Key)
		if err == nil {
			body := bytes.NewBuffer(data)
			req, err = http.NewRequest("PUT", target, body)
			if err == nil {
				req.Header.Set(headers.HeaderContentType, headers.MIMEApplicationJSONCharsetUTF8)
			}
		}
	}
	return
}

func BenchmarkCreateFormTraditionally(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildFormRequestTraditionally(changeParam)
	}
}

func BenchmarkCreateFormByGotten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.Change(changeParam)
	}
}

func BenchmarkCreateMultipartTraditionally(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildMultipartRequestTraditionally(uploadParam)
	}
}

func BenchmarkCreateMultipartByGotten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		service.Upload(uploadParam)
	}
}

func BenchmarkCreateJSONReqTraditionally(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buildJSONRequestTraditionally(updateParam)
	}
}

func BenchmarkCreateJSONReqByGotten(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gitService.UpdateDeployKey(updateParam)
	}
}
