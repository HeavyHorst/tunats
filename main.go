package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

type Request struct {
	URL        string
	Method     string
	Header     http.Header
	Form       url.Values
	RemoteAddr string
	Body       []byte
}

type Response struct {
	Header     http.Header
	StatusCode int
	Body       []byte
	Error      string
}

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
}
var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

func requestHandler(nc *nats.Conn, subject, to string) {
	nc.Subscribe(subject, func(m *nats.Msg) {
		go func() {
			var req Request
			var resp Response
			msgpack.Unmarshal(m.Data, &req)

			request, err := http.NewRequest(req.Method, to+req.URL, bytes.NewReader(req.Body))
			if err != nil {
				resp.Error = err.Error()
			} else {
				request.Header = req.Header
				request.Form = req.Form
				request.RemoteAddr = req.RemoteAddr

				response, err := netClient.Do(request)
				if err != nil {
					resp.Error = err.Error()
				} else {
					bd, _ := ioutil.ReadAll(response.Body)
					defer response.Body.Close()
					resp.Body = bd
					resp.Header = response.Header
					resp.StatusCode = response.StatusCode
				}
			}

			jdat, _ := msgpack.Marshal(resp)
			m.Respond(jdat)
		}()
	})
}

func httpServer(nc *nats.Conn, subject string, port int) {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		bd, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()

		r := Request{
			URL:        req.URL.String(),
			Method:     req.Method,
			Header:     req.Header,
			Form:       req.Form,
			RemoteAddr: req.RemoteAddr,
			Body:       bd,
		}

		jdat, _ := msgpack.Marshal(r)
		m, err := nc.Request(subject, jdat, 5*time.Second)
		if err != nil {
			fmt.Println(err)
			return
		}

		var resp Response
		msgpack.Unmarshal(m.Data, &resp)
		for k, v := range resp.Header {
			for _, hv := range v {
				w.Header().Add(k, hv)
			}
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(resp.Body)
	})
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}

func main() {
	httpPort := flag.Int("port", 8090, "")
	serve := flag.Bool("serve", false, "run the http server")
	subject := flag.String("sub", "", "the nats subject to listen on/send to")
	forward := flag.Bool("forward", false, "forward nats messages to an http server")
	insecureSkipVerify := flag.Bool("insecureSkipVerify", false, "")
	remoteURL := flag.String("to", "", "the remote http server")
	nurl := flag.String("nurl", nats.DefaultURL, "the nats cluster url")
	ncreds := flag.String("creds", "", "the path to the nats credentials")
	nname := flag.String("nats_name", "", "nats connection name")
	flag.Parse()

	if *insecureSkipVerify {
		netTransport.TLSClientConfig.InsecureSkipVerify = true
	}

	nc, err := nats.Connect(*nurl, nats.UserCredentials(*ncreds), nats.Name(*nname))
	if err != nil {
		panic(err)
	}

	if *forward {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		requestHandler(nc, *subject, *remoteURL)
		<-sigs
	}

	if *serve {
		httpServer(nc, *subject, *httpPort)
	}
}
