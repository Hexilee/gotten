package gotten

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
)

const (
	ComplexPath = "/post/{year}/{m}/{d}/{hash_id}"
)

type (
	Day struct {
		day int
	}

	SimpleParams struct {
		Year   int          `type:"path" default:"2018"`
		Month  string       `type:"path" key:"m"`
		Day    fmt.Stringer `type:"path" key:"d"`
		HashId string       `type:"path"`
		Page   int          `type:"query" default:"1"`
		Num    int          `type:"query" require:"true"`
	}
)

func (day Day) String() string {
	return strconv.Itoa(day.day)
}

func TestPathKeyList(t *testing.T) {
	// test addKey
	keyList := make(PathKeyList)
	assert.True(t, keyList.addKey("1"))
	assert.True(t, keyList.addKey("2"))
	assert.False(t, keyList.addKey("1"))

	// test deleteKey
	assert.False(t, keyList.deleteKey("3"))
	assert.True(t, keyList.deleteKey("1"))
	assert.False(t, keyList.deleteKey("1"))

	// test empty
	for _, testCase := range []struct {
		listOne []string
		listTwo []string
		result  bool
	}{
		{[]string{"1", "2"}, nil, false},
		{[]string{"1", "2"}, []string{}, false},
		{[]string{"1", "2"}, []string{"1"}, false},
		{[]string{"1", "2"}, []string{"1", "2"}, true},
	} {
		keyList := make(PathKeyList)
		for _, key := range testCase.listOne {
			keyList.addKey(key)
		}

		for _, key := range testCase.listTwo {
			keyList.deleteKey(key)
		}

		assert.Equal(t, testCase.result, keyList.empty())
	}

}

func TestPathKeyRegexp(t *testing.T) {
	assert.NotNil(t, pathKeyRegexp)
	// find
	for _, testCase := range []struct {
		Src    string
		Result []string
	}{
		{`/user/{_}`, []string{`{_}`}},
		{`/user/{Gid}`, []string{`{Gid}`}},
		{`/user/{gid}`, []string{`{gid}`}},
		{`/user/{group_id}`, []string{`{group_id}`}},
		{`/user/{group1}`, []string{`{group1}`}},
		{`/user/{gid}/{uid}`, []string{`{gid}`, `{uid}`}},
		{`/user/{group-id}`, []string{}},
		{`/user/{}`, []string{}},
		{`/user/{0gid}`, []string{}},
	} {
		result := pathKeyRegexp.FindAllString(testCase.Src, -1)
		for i := range result {
			assert.Equal(t, testCase.Result[i], result[i])
		}
	}
}

func TestFieldExportable(t *testing.T) {
	for _, testCase := range []struct {
		name       string
		exportable bool
	}{
		{"Name", true},
		{"N", true},
		{"nAME", false},
		{"name", false},
		{"n", false},
	} {
		assert.Equal(t, testCase.exportable, fieldExportable(testCase.name))
	}
}

func TestVarsParser(t *testing.T) {
	parser, err := newVarsParser(ComplexPath)
	assert.Nil(t, err)
	assert.Nil(t, parser.parse(reflect.TypeOf(new(SimpleParams))))

	for _, testCase := range []struct {
		params *SimpleParams
		path   string
		query  string
	}{
		{&SimpleParams{
			Month:  "1",
			Day:    Day{1},
			HashId: "1",
			Num:    10,
		}, `/post/2018/1/1/1`, `num=10&page=1`},
	} {
		ctr := parser.Builder()
		assert.Nil(t, ctr.setValues(reflect.ValueOf(testCase.params)))
		result, err := ctr.getUrl()
		assert.Nil(t, err)
		assert.Equal(t, testCase.path, result.Path)
		assert.Equal(t, testCase.query, result.RawQuery)
	}
}
