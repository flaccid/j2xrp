package main

import (
	"bytes"
	"github.com/clbanning/mxj"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type transport struct {
	http.RoundTripper
}

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {
		Trace = log.New(traceHandle,
			"TRACE: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Info = log.New(infoHandle,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Warning = log.New(warningHandle,
			"WARNING: ",
			log.Ldate|log.Ltime|log.Lshortfile)

		Error = log.New(errorHandle,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)
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
		Error.Println(err)
	}
	Info.Println(string(responseDump))

	return resp, nil
}

var _ http.RoundTripper = &transport{}

func convJ2X(json []byte) []byte {
  m, err := mxj.NewMapJson(json)
  if err != nil {
    Error.Println("error mapping json: ", err)
  }

  xml, err := m.Xml()
  if err != nil {
    Error.Println("error converting xml: ", err)
  }

  return xml
}

func main() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	// create a reverse proxy to rightscale
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host:   "wstunnel10-1.rightscale.com:443",
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
			Error.Println(err)
		}
		Info.Println(string(requestDump))

		realDirector(req)
	}

	http.Handle("/", proxy)
	Info.Println("starting http server")
	log.Fatal(http.ListenAndServe(":9090", proxy))
}
