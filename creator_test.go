package gotten_test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

type (
	SampleService struct {
		GetPosts      func(params *GetPostsParams) (gotten.Response, error)      `path:"/post/{year}/{month}/{day}"`
		AddPost       func(params *AddPostParams) (gotten.Response, error)       `method:"POST" path:"/post/{year}/{month}/{day}"`
		AddPostByForm func(params *AddPostByFormParams) (gotten.Response, error) `method:"POST" path:"/post"`
		UploadAvatar  func(params *UploadAvatarParams) (gotten.Response, error)  `method:"POST" path:"/avatar"`
	}

	GetPostsParams struct {
		Year  int `type:"path"`
		Month int `type:"path"`
		Day   int `type:"path"`
		Page  int `type:"query" default:"1"`
		Limit int `type:"query" default:"15"`
	}

	AddPostParams struct {
		Year  int       `type:"path"`
		Month int       `type:"path"`
		Day   int       `type:"path"`
		Post  *TestPost `type:"json"`
	}

	AddPostByFormParams struct {
		Year  int       `type:"form"`
		Month int       `type:"form"`
		Day   int       `type:"form"`
		Post  *TestPost `type:"json"`
	}

	AvatarDescription struct {
		Creator   string
		CreatedAt time.Time
	}

	UploadAvatarParams struct {
		Uid         int                `type:"part"`
		Username    string             `type:"part"`
		Avatar      gotten.FilePath    `type:"part"`
		Description *AvatarDescription `type:"json"`
	}

	UploadedData struct {
		Hash      string
		Uid       int
		Username  string
		Filename  string
		FileSize  int64
		Creator   string
		CreatedAt time.Time
	}
)

func TestCreator_Impl(t *testing.T) {
	creator, err := gotten.NewBuilder().
		SetBaseUrl("https://mock.io").
		SetClient(mockClient).
		Build()

	assert.Nil(t, err)
	service := new(SampleService)
	assert.Nil(t, creator.Impl(service))
	assert.NotNil(t, service.GetPosts)
	resp, err := service.GetPosts(&GetPostsParams{2018, 10, 1, 1, 10})
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	results := make([]TestPost, 0)
	assert.Nil(t, resp.Unmarshal(&results))

	data, err := json.Marshal(&results)
	assert.Nil(t, err)
	assert.Equal(t, "[{\"author\":\"Hexilee\",\"title\":\"Start!\",\"content\":\"Hello world!\"}]", string(data))

	resp, err = service.AddPost(&AddPostParams{
		Year:  2018,
		Month: 10,
		Day:   1,
		Post:  &TestPost{"Hexilee", "AddPost Test", "Success!"},
	})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	var addedResult AddedData
	assert.Nil(t, resp.Unmarshal(&addedResult))
	assert.True(t, addedResult.Success)
	assert.Equal(t, 2018, addedResult.Year)
	assert.Equal(t, 10, addedResult.Month)
	assert.Equal(t, 1, addedResult.Day)
	assert.Equal(t, 2, addedResult.Order)

	resp, err = service.AddPostByForm(&AddPostByFormParams{
		Year:  2018,
		Month: 10,
		Day:   1,
		Post:  &TestPost{"Hexilee", "AddPostByForm Test", "Success!"},
	})

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	assert.Nil(t, resp.Unmarshal(&addedResult))
	assert.True(t, addedResult.Success)
	assert.Equal(t, 2018, addedResult.Year)
	assert.Equal(t, 10, addedResult.Month)
	assert.Equal(t, 1, addedResult.Day)
	assert.Equal(t, 3, addedResult.Order)

	now := time.Now()
	resp, err = service.UploadAvatar(&UploadAvatarParams{
		Uid:      1,
		Username: "Hexilee",
		Avatar:   "testAssets/Concurrency-in-Go.pdf",
		Description: &AvatarDescription{
			Creator:   "Hexilee",
			CreatedAt: now,
		},
	})

	var uploadedData UploadedData

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode())
	assert.Equal(t, headers.MIMEApplicationJSONCharsetUTF8, resp.ContentType())
	assert.Nil(t, resp.Unmarshal(&uploadedData))
	assert.Equal(t, 1, uploadedData.Uid)
	assert.Equal(t, "Hexilee", uploadedData.Username)
	assert.Equal(t, "Concurrency-in-Go.pdf", uploadedData.Filename)
	assert.Equal(t, "Hexilee", uploadedData.Creator)
	assert.True(t, now.Equal(uploadedData.CreatedAt))

	file, err := os.Open("testAssets/Concurrency-in-Go.pdf")
	assert.Nil(t, err)
	h := md5.New()
	n, err := io.Copy(h, file)
	assert.Nil(t, err)
	assert.Equal(t, n, uploadedData.FileSize)
	assert.Equal(t, fmt.Sprintf("%x", h.Sum(nil)), uploadedData.Hash)
}
