package apiserver

import (
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/nyaxt/otaru/apiserver"
	"github.com/nyaxt/otaru/cli"
	"github.com/nyaxt/otaru/logger"
)

type apiproxy struct {
	cfg *cli.CliConfig
}

var reMatch = regexp.MustCompile(`^/proxy/(\w+)(/.*)$`)

func copyHeader(d, s http.Header) {
	for k, v := range s {
		d[k] = v
	}
}

func (ap *apiproxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ms := reMatch.FindStringSubmatch(path)
	if len(ms) != 3 {
		http.Error(w, "Invalid proxy URL.", http.StatusBadRequest)
		return
	}
	// fmt.Printf("path %v match %v", path, ms)

	hname, tgtpath := ms[1], ms[2]

	ep, tc, err := cli.ConnectionInfo(ap.cfg, hname)
	if err != nil {
		http.Error(w, "Failed to construct connection info.", http.StatusInternalServerError)
		return
	}

	cli := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tc,
		},
	}
	url := &url.URL{
		Scheme:   "https",
		Host:     ep,
		Path:     tgtpath,
		RawQuery: r.URL.RawQuery,
	}
	logger.Debugf(mylog, "URL: %v", url)

	// FIXME: filter tgtpath
	// FIXME: add forwarded-for

	preq := &http.Request{
		Method: r.Method,
		Header: make(http.Header),
		URL:    url,
		Body:   r.Body,
	}
	copyHeader(preq.Header, r.Header)

	presp, err := cli.Do(preq)
	if err != nil {
		http.Error(w, "Failed to issue request.", http.StatusInternalServerError)
		return
	}
	defer presp.Body.Close()

	// logger.Debugf(mylog, "resp st: %d", presp.StatusCode)

	copyHeader(w.Header(), presp.Header)
	w.WriteHeader(presp.StatusCode)
	io.Copy(w, presp.Body)
}

func InstallProxyHandler(cfg *cli.CliConfig) apiserver.Option {
	return apiserver.AddMuxHook(func(mux *http.ServeMux) {
		mux.Handle("/proxy/", &apiproxy{cfg})
	})
}