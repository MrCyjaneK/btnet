package proxy

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"git.mrcyjanek.net/mrcyjanek/btnet/btclient"
	"git.mrcyjanek.net/mrcyjanek/btnet/utils"
)

// Hop-by-hop headers. These are removed when sent to the backend.
// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html
var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}

type proxy struct {
}

func (p *proxy) ServeHTTP(wr http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme != "http" {
		http.Error(wr, "[500] Server Error - unsupported protocol scheme", http.StatusInternalServerError)
		return
	}

	client := &http.Client{}

	//http: Request.RequestURI can't be set in client requests.
	//http://golang.org/src/pkg/net/http/client.go
	req.RequestURI = ""

	delHopHeaders(req.Header)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		appendHostToXForwardHeader(req.Header, clientIP)
	}
	magnetok, _ := utils.IsValidMagnetBTnet(req.Host)
	//log.Println(message)
	if !magnetok {
		resp, err := client.Do(req)
		if err != nil {
			http.Error(wr, "[500] Server Error - "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		delHopHeaders(resp.Header)

		copyHeader(wr.Header(), resp.Header)
		wr.WriteHeader(resp.StatusCode)
		io.Copy(wr, resp.Body)
	} else {
		btclient.Serve(wr, req.Host, req.URL.Path)
		//http.Error(wr, "[200] - OK! Please refresh. "+requestPath, http.StatusOK)
	}
}

// Init - Start the proxy
func Init() {
	var addr = flag.String("addr", "127.0.0.1:8080", "The addr of the application.")
	flag.Parse()

	handler := &proxy{}

	log.Println("Starting proxy server on", *addr)
	if err := http.ListenAndServe(*addr, handler); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
