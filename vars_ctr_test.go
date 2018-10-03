package gotten

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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
