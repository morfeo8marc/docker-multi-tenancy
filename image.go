package main

import (
	"net/http"
	"regexp"
"encoding/json"
	"bytes"
	"bufio"
	"io/ioutil"
	"encoding/binary"
	"log"
)

// APIImages represent an image returned in the ListImages call.
type APIImages struct {
	ID          string            `json:"Id" yaml:"Id"`
	RepoTags    []string          `json:"RepoTags,omitempty" yaml:"RepoTags,omitempty"`
	Created     int64             `json:"Created,omitempty" yaml:"Created,omitempty"`
	Size        int64             `json:"Size,omitempty" yaml:"Size,omitempty"`
	VirtualSize int64             `json:"VirtualSize,omitempty" yaml:"VirtualSize,omitempty"`
	ParentID    string            `json:"ParentId,omitempty" yaml:"ParentId,omitempty"`
	RepoDigests []string          `json:"RepoDigests,omitempty" yaml:"RepoDigests,omitempty"`
	Labels      map[string]string `json:"Labels,omitempty" yaml:"Labels,omitempty"`
}

type imageTransformer struct{}

func NewImageListTransformer() *Transformer {
	t := &Transformer{}

	t.Regexp = regexp.MustCompile("/v.*/images/json")

	t.MultiTenancyTransformer = &imageTransformer{}

	return t
}

func (c *imageTransformer) transformRequest(r *http.Request) {
	log.Println("NewImageListTransformer.transformRequest")
}

func (c *imageTransformer) transformResponse(r *http.Response) {

	var images []APIImages

	if err := json.NewDecoder(r.Body).Decode(&images); err != nil {
		return
	}

	for _, im :=  range images{

		if im.Labels == nil {
			im.Labels = make(map[string]string)
		}

		im.Labels["hola"] = "world"
	}


	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	// Now take the struct and encode it
	if err := json.NewEncoder(w).Encode(&images); err != nil {
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b.Bytes()))

	// Set size of modified body
	r.ContentLength = int64(binary.Size(b))
}
