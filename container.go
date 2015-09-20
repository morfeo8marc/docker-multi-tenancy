package main
import (
	"regexp"
	"net/http"
"encoding/json"
	"bytes"
	"bufio"
	"io/ioutil"
	"encoding/binary"
	"net/url"
	"log"
)

// APIPort is a type that represents a port mapping returned by the Docker API
type APIPort struct {
	PrivatePort int64  `json:"PrivatePort,omitempty" yaml:"PrivatePort,omitempty"`
	PublicPort  int64  `json:"PublicPort,omitempty" yaml:"PublicPort,omitempty"`
	Type        string `json:"Type,omitempty" yaml:"Type,omitempty"`
	IP          string `json:"IP,omitempty" yaml:"IP,omitempty"`
}

// APIContainers represents each container in the list returned by
// ListContainers.
type APIContainers struct {
	ID         string            `json:"Id" yaml:"Id"`
	Image      string            `json:"Image,omitempty" yaml:"Image,omitempty"`
	Command    string            `json:"Command,omitempty" yaml:"Command,omitempty"`
	Created    int64             `json:"Created,omitempty" yaml:"Created,omitempty"`
	Status     string            `json:"Status,omitempty" yaml:"Status,omitempty"`
	Ports      []APIPort         `json:"Ports,omitempty" yaml:"Ports,omitempty"`
	SizeRw     int64             `json:"SizeRw,omitempty" yaml:"SizeRw,omitempty"`
	SizeRootFs int64             `json:"SizeRootFs,omitempty" yaml:"SizeRootFs,omitempty"`
	Names      []string          `json:"Names,omitempty" yaml:"Names,omitempty"`
	Labels     map[string]string `json:"Labels,omitempty" yaml:"Labels, omitempty"`
}

type containerTransformer struct{}

func NewContainerListTransformer() *Transformer {
	t := &Transformer{}

	t.Regexp = regexp.MustCompile("/v.*/containers/json")

	t.MultiTenancyTransformer = &containerTransformer{}

	return t
}

func (c *containerTransformer) transformRequest(r *http.Request) {
	log.Println("NewContainerListTransformer.transformRequest")

	items := url.Values(map[string][]string{})

	s := make([]string, 1, 1)

	s[0] = "{\"labels\":[\"hello:world\"]}"
	items["filters"] = s

	r.URL.RawQuery = items.Encode()
}

func (c *containerTransformer) transformResponse(r *http.Response) {

	var cS []APIContainers
	var cSFilter []APIContainers

	if err := json.NewDecoder(r.Body).Decode(&cS); err != nil {
		return
	}

	cSFilter = make([]APIContainers, 0, len(cS))

	for _, c := range cS {

		if v, ok := c.Labels["hello"]; ok {

			if v == "world" {
				cSFilter = append(cSFilter, c)
			}
		}

	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	// Now take the struct and encode it
	if err := json.NewEncoder(w).Encode(&cSFilter); err != nil {
		return
	}

	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b.Bytes()))

	// Set size of modified body
	r.ContentLength = int64(binary.Size(b))
}
