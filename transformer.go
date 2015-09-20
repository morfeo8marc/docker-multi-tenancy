package main

import (
	"net/http"
	"regexp"
)

type MultiTenancyTransformer interface {
	transformRequest(r *http.Request)
	transformResponse(r *http.Response)
}

type Transformer struct {
	Regexp      *regexp.Regexp
	            MultiTenancyTransformer
}

func (t *Transformer) String() string {
	return "Transformer Regexp" + t.Regexp.String()
}

var DockerTransformers []*Transformer

func init() {
	DockerTransformers = NewDefaultTransformers()
}

func NewDefaultTransformers() []*Transformer {
	ts := make([]*Transformer, 0)

	ts = append(ts, NewImageListTransformer())
	ts = append(ts, NewContainerListTransformer())

	return ts
}

