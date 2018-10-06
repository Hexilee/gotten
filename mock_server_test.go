package gotten_test

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/Hexilee/gotten"
	"github.com/Hexilee/gotten/headers"
	"github.com/Hexilee/gotten/mock"
	"github.com/go-chi/chi"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	database   *Database
	router     chi.Router
	mockClient gotten.Client
)

func init() {
	database = newDatabase()
	go database.Run()

	database.Add(2018, 10, 1, &TestPost{"Hexilee", "Start!", "Hello world!"})

	router = chi.NewRouter()
	router.Get("/post/{year}/{month}/{day}", getPost)
	router.Post("/post/{year}/{month}/{day}", addPost)
	router.Post("/post", addPostByForm)
	router.Post("/avatar", addAvatar)

	mockBuilder := mock.NewClientBuilder()
	mockBuilder.Register("mock.io", router)
	mockClient = mockBuilder.Build()
}

func getPost(w http.ResponseWriter, r *http.Request) {
	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))
	day, _ := strconv.Atoi(chi.URLParam(r, "day"))
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	posts := database.Get(year, month, day, page, limit)
	result, _ := json.Marshal(&posts)
	w.Header().Set(headers.HeaderContentType, headers.MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func addPost(w http.ResponseWriter, r *http.Request) {
	var result AddedData
	defer func() {
		w.Header().Set(headers.HeaderContentType, headers.MIMEApplicationJSONCharsetUTF8)
		w.WriteHeader(http.StatusCreated)
		respData, _ := json.Marshal(&result)
		w.Write(respData)
	}()

	year, _ := strconv.Atoi(chi.URLParam(r, "year"))
	month, _ := strconv.Atoi(chi.URLParam(r, "month"))
	day, _ := strconv.Atoi(chi.URLParam(r, "day"))

	body := r.Body
	defer body.Close()
	data, err := ioutil.ReadAll(body)
	if err != nil {
		result.Success = false
		return
	}

	var post TestPost
	err = json.Unmarshal(data, &post)
	if err != nil {
		result.Success = false
		return
	}

	result = *database.Add(year, month, day, &post)
}

func addPostByForm(w http.ResponseWriter, r *http.Request) {
	var result AddedData
	defer func() {
		w.Header().Set(headers.HeaderContentType, headers.MIMEApplicationJSONCharsetUTF8)
		w.WriteHeader(http.StatusCreated)
		respData, _ := json.Marshal(&result)
		w.Write(respData)
	}()

	r.ParseForm()
	year, _ := strconv.Atoi(r.PostForm.Get("year"))
	month, _ := strconv.Atoi(r.PostForm.Get("month"))
	day, _ := strconv.Atoi(r.PostForm.Get("day"))
	postStr := r.PostForm.Get("post")

	var post TestPost
	err := json.Unmarshal([]byte(postStr), &post)
	if err != nil {
		result.Success = false
		return
	}

	result = *database.Add(year, month, day, &post)
}

func addAvatar(w http.ResponseWriter, r *http.Request) {
	var result UploadedData
	defer func() {
		w.Header().Set(headers.HeaderContentType, headers.MIMEApplicationJSONCharsetUTF8)
		w.WriteHeader(http.StatusCreated)
		respData, _ := json.Marshal(&result)
		w.Write(respData)
	}()

	defer r.Body.Close()

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("avatar")

	if err != nil {
		fmt.Printf("%#v\n", err.Error())
		return
	}

	result.Filename = handler.Filename
	result.FileSize = handler.Size
	result.Uid, err = strconv.Atoi(r.PostFormValue("uid"))
	if err != nil {
		return
	}

	result.Username = r.PostFormValue("username")
	descData := r.PostFormValue("description")
	var description AvatarDescription

	if json.Unmarshal([]byte(descData), &description) != nil {
		return
	}

	result.Creator = description.Creator
	result.CreatedAt = description.CreatedAt

	h := md5.New()
	io.Copy(h, file)
	result.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

type (
	TestPost struct {
		Author  string `json:"author"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	GetterData struct {
		dataChan                      chan<- []*TestPost
		year, month, day, page, limit int
	}

	AdderData struct {
		dataChan         chan<- *AddedData
		year, month, day int
		post             *TestPost
	}

	AddedData struct {
		Success bool `json:"success"`
		Year    int  `json:"year"`
		Month   int  `json:"month"`
		Day     int  `json:"day"`
		Order   int  `json:"order"`
	}

	Database struct {
		cancel     context.CancelFunc
		getterChan chan *GetterData
		adderChan  chan *AdderData
		data       map[string][]*TestPost
	}
)

func newDatabase() *Database {
	return &Database{
		cancel:     func() {},
		getterChan: make(chan *GetterData, 100),
		adderChan:  make(chan *AdderData, 100),
		data:       make(map[string][]*TestPost),
	}
}

func (database *Database) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	database.cancel = cancel
Loop:
	for {
		select {
		case getterData := <-database.getterChan:
			var resultList []*TestPost
			totalList := database.data[genKey(getterData.year, getterData.month, getterData.day)]
			length := len(totalList)
			offset := (getterData.page - 1) * getterData.limit
			max := getterData.page * getterData.limit
			if length <= offset {
				resultList = make([]*TestPost, 0)
			} else if length <= max {
				resultList = make([]*TestPost, length-offset)
				copy(resultList, totalList[offset:])
			} else {
				resultList = make([]*TestPost, getterData.limit)
				copy(resultList, totalList[offset:max-1])
			}
			getterData.dataChan <- resultList
		case adderData := <-database.adderChan:
			key := genKey(adderData.year, adderData.month, adderData.day)
			if database.data[key] == nil {
				database.data[key] = make([]*TestPost, 0)
			}
			database.data[key] = append(database.data[key], adderData.post)
			adderData.dataChan <- &AddedData{true, adderData.year, adderData.month, adderData.day, len(database.data[key])}
		case <-ctx.Done():
			break Loop
		}
	}
}

func (database *Database) Stop() {
	database.cancel()
}

func (database *Database) Get(year, month, day, page, limit int) []*TestPost {
	dataChan := make(chan []*TestPost, 1)
	database.getterChan <- &GetterData{dataChan, year, month, day, page, limit}
	select {
	case posts := <-dataChan:
		return posts
	}
}

func (database *Database) Add(year, month, day int, post *TestPost) *AddedData {
	dataChan := make(chan *AddedData, 1)
	database.adderChan <- &AdderData{dataChan, year, month, day, post}
	select {
	case added := <-dataChan:
		return added
	}
}

func genKey(year, month, day int) string {
	return fmt.Sprintf("%d-%d-%d", year, month, day)
}
