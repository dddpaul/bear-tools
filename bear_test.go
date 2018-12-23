package main

import (
	"github.com/magiconair/properties/assert"
	"strings"
	"testing"
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
	assert.Equal(t, NewNote(strings.NewReader(`
# Title
#tag
What's difference between [[Thrift]] and [[GRPC]]
Or between [[GRPC]] and [[Thrift]] and [[SOAP]]
Empty Links [[]] should be ignored
`)).Links, map[string]int{"Thrift": 2, "GRPC": 2, "SOAP": 1})
}

func TestMarshal(t *testing.T) {
	json, _ := NewNote(strings.NewReader(`
	# Title
	#tag1 #tag2
	What's difference between [[Thrift]] and [[GRPC]]
	Or between [[GRPC]] and [[Thrift]] and [[SOAP]]
	`)).Marshal()
	assert.Equal(t, json, `{"title":"","tags":["tag1","tag2"],"links":{"GRPC":2,"SOAP":1,"Thrift":2}}`)
}
