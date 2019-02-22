package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidTitle(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Valid Title
Content
`)).Title, "Valid Title")
}

func TestInvalidTitle(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
#Invalid Title
Content
`)).Title, "")
}

func TestTags(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1 #tag2
Content
`)).Tags, []string{"tag1", "tag2"})
}

func TestEmptyTags(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1 # #
Content
`)).Tags, []string{"tag1"})
}

func TestTagsWithGarbage(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1 # some string
Content
`)).Tags, []string{"tag1"})
}

func TestMultilineTags(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1 #tag2
#tag3
Content
`)).Tags, []string{"tag1", "tag2", "tag3"})
}

func TestDuplicateTags(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1 #tag2
#tag2 #tag1
Content
`)).Tags, []string{"tag1", "tag2"})
}

func TestTagsInContent(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag1
#tag2 https://example.tld/path/to/file#slides=1,2,3 #tag3 some string #tag4
`)).Tags, []string{"tag1", "tag2", "tag3", "tag4"})
}

func TestMultipleHashesAreNotTags(t *testing.T) {
	assert.Equal(t, NewNote(strings.NewReader(`
#tag1 ## ### #### #tag2 ##### #tag3
`)).Tags, []string{"tag1", "tag2", "tag3"})
}

func TestLinks(t *testing.T) {
	assert.ElementsMatch(t, NewNote(strings.NewReader(`
# Title
#tag
What's difference between [[Thrift]] and [[GRPC]]
Or between [[GRPC]] and [[Thrift]] and [[SOAP]]
Empty Links [[]] should be ignored
`)).Links, []Link{
		{
			Title: "Thrift",
			Count: 2,
		},
		{
			Title: "GRPC",
			Count: 2,
		},
		{
			Title: "SOAP",
			Count: 1,
		},
	})
}

func TestMarshal(t *testing.T) {
	note, _ := NewNote(strings.NewReader(`
	# Title
	#tag1 #tag2
	What's difference between [[Thrift]] and [[GRPC]]
	Or between [[GRPC]] and [[Thrift]] and [[SOAP]]
	`)).Marshal()

	expected := `{"title":"","tags":["tag1","tag2"],"links":[{"title":"Thrift","count":2},{"title":"GRPC","count":2},{"title":"SOAP","count":1}]}`

	ok, _ := AreEqualJSON(expected, note)
	assert.True(t, ok)
}

func AreEqualJSON(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}
