package runtime

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Proxy provides a web server proxy
type Proxy interface {
	// Run bootstraps the web service proxy
	Run(config *Config) error
	io.Closer
}

type proxy struct {
	listener net.Listener
	proxy    *httputil.ReverseProxy
	builder  Builder
	runner   Runner
	to       *url.URL
}

// NewProxy constructs a new Proxy
func NewProxy(builder Builder, runner Runner) Proxy {
	return &proxy{
		builder: builder,
		runner:  runner,
	}
}

func (p *proxy) Run(config *Config) error {
	// create our reverse proxy
	url, err := url.Parse(config.ProxyTo)
	if err != nil {
		return err
	}
	p.proxy = httputil.NewSingleHostReverseProxy(url)
	p.to = url

	server := http.Server{Handler: http.HandlerFunc(p.defaultHandler)}

	if config.CertFile != "" && config.KeyFile != "" {
		cer, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			return err
		}

		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}

		p.listener, err = tls.Listen("tcp", fmt.Sprintf("%s:%d", config.Laddr, config.Port), server.TLSConfig)
		if err != nil {
			return err
		}
	} else {
		p.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", config.Laddr, config.Port))
		if err != nil {
			return err
		}
	}

	go server.Serve(p.listener)

	return nil
}

func (p *proxy) Close() error {
	return p.listener.Close()
}

func (p *proxy) defaultHandler(res http.ResponseWriter, req *http.Request) {
	errors := p.builder.Errors()
	if len(errors) > 0 {
		t := template.Must(template.New("errors").Parse(tplError))
		safe := template.HTMLEscapeString(errors)
		safe = strings.Replace(safe, "\n", "<br>", -1)
		if err := t.Execute(res, template.HTML(safe)); err != nil {
			res.Write([]byte(errors))
		}
	} else {
		p.runner.Run()
		if strings.ToLower(req.Header.Get("Upgrade")) == "websocket" || strings.ToLower(req.Header.Get("Accept")) == "text/event-stream" {
			proxyWebsocket(res, req, p.to)
		} else {
			p.proxy.ServeHTTP(res, req)
		}
	}
}

func proxyWebsocket(w http.ResponseWriter, r *http.Request, host *url.URL) {
	d, err := net.Dial("tcp", host.Host)
	if err != nil {
		http.Error(w, "Error contacting backend server.", 500)
		fmt.Errorf("Error dialing websocket backend %s: %v", host, err)
		return
	}
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Not a hijacker?", 500)
		return
	}
	nc, _, err := hj.Hijack()
	if err != nil {
		fmt.Errorf("Hijack error: %v", err)
		return
	}
	defer nc.Close()
	defer d.Close()

	err = r.Write(d)
	if err != nil {
		fmt.Errorf("Error copying request to target: %v", err)
		return
	}

	errc := make(chan error, 2)
	cp := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go cp(d, nc)
	go cp(nc, d)
	<-errc
}

var tplError = `
<!DOCTYPE HTML>
<html>
  <head>
    <title>Build Error</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no" />
    <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.1.0/css/all.css" />
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" />
  </head>
  <body>
    <nav class="navbar navbar-inverse navbar-static-top">
      <div class="container-fluid">
        <div class="navbar-header">
          <a class="navbar-brand" href="#"> <i class="fas fa-sync-alt"></i> reload</a>
        </div>
      </div>
    </nav>
    <br/>
    <div class="container">
      <div class="jumbotron">
        <h1><i class="fa fa-exclamation-circle text-danger"></i> Build Failed</h1>
        <p class="bg-danger">
          {{ . }}
        </p>
      </div>
    </div>
  </body>
</html>
`
