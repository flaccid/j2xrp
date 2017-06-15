package proxy

import (
	"bytes"
	"github.com/clbanning/mxj"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type transport struct {
	http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	b = bytes.Replace(b, []byte("server"), []byte("schmerver"), -1)
	body := ioutil.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Error(err)
	}
	log.Info(string(responseDump))

	return resp, nil
}

var _ http.RoundTripper = &transport{}

func convJ2X(json []byte) []byte {
	m, err := mxj.NewMapJson(json)
	if err != nil {
		log.Error("error mapping json: ", err)
	}

	xml, err := m.Xml()
	if err != nil {
		log.Error("error converting xml: ", err)
	}

	return xml
}

// Serve : creates a reverse proxy to forward XML requests converted to JSON
func Serve(scheme string, host string, listenPort string) {
	// create a reverse proxy to rightscale
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: scheme,
		Host:   host,
	})

	// take control of the proxy transport and director
	proxy.Transport = &transport{http.DefaultTransport}
	realDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		// modify headers to taste
		req.Header.Set("X-Reverse-Proxy", "j2xrp")
		req.Header.Set("Content-Type", "application/xml")

		// request body JSON->XML happens here
		// we only care about body content if this is a PUT or POST
		if req.Method == "PUT" || req.Method == "POST" {
			// stream the body and convert the expected JSON to XML
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			s := buf.String()
			xml := convJ2X([]byte(s))
			req.Body = ioutil.NopCloser(strings.NewReader(string(xml)))
			req.ContentLength = int64(len(string(xml)))
		}

		// dump the ingress request
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Error(err)
		}
		log.Info(string(requestDump))

		realDirector(req)
	}

	http.Handle("/", proxy)
	log.Info("starting http server")
	log.Info("proxy requests to ", scheme, "://", host)
	log.Info("listening for requests on :" + listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, proxy))
}
