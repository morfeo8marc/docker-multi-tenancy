package main
import (
	"net/http"
	"log"
	"io/ioutil"
	"net/http/httputil"
)

type ProxyConfiguration struct {
	DockerEndpoint string
	Address string
}

type Proxy struct {
	client *Client
	transformers []*Transformer
}


func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var trans *Transformer

	if logRequest {
		if respBytes, err := httputil.DumpRequest(r, true); err == nil {
			log.Printf("Request:\n%v\n", string(respBytes))
		}

	}

	log.Println("Recieved method ", r.Method, " path", r.URL.Path)

	for _, t := range p.transformers{

		m := t.Regexp.FindStringSubmatch(r.URL.Path)

		if m != nil {
			trans = t
			break
		}
	}

	if trans != nil {
		log.Printf("Applying transformation %v to request", trans)
		trans.transformRequest(r)
		if logRequest {
			if respBytes, err := httputil.DumpRequest(r, true); err == nil {
				log.Printf("Modified Request:\n%v\n", string(respBytes))
			}

		}
	}

	resp, err := p.client.do( r.Method,  r.URL.String(), r.Body)

	if logResponse {
		if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
			log.Printf("Response:\n%v\n", string(respBytes))
		}

	}

	if trans != nil && err == nil {
		log.Println("Applying transformation to response")
		trans.transformResponse(resp)

		if logResponse {
			if respBytes, err := httputil.DumpResponse(resp, true); err == nil {
				log.Printf("Modified Response:\n%v\n", string(respBytes))
			}

		}
	}

	if resp != nil {
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		w.Write(content)
	}
}

func StartProxy(pConf *ProxyConfiguration) {

	log.Println("Starting multi-tenancy proxy listening at ", pConf.Address, " and forwarding it to ", pConf.DockerEndpoint)

	c, err := NewClient(pConf.DockerEndpoint)

	if err != nil {
		log.Fatalf("Error creating client connection: %v", err)
		return
	}

	p := &Proxy{client: c,
	            transformers: DockerTransformers}


	err = http.ListenAndServe(pConf.Address, p)

	if err != nil {
		log.Fatal(err)
	}
}
