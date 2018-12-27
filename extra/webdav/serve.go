package webdav

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/nyaxt/otaru/cli"
	"github.com/nyaxt/otaru/logger"
)

var mylog = logger.Registry().Category("fe-webdav")
var accesslog = logger.Registry().Category("http-webdav")

func Serve(cfg *cli.CliConfig, closeC <-chan struct{}) error {
	wcfg := cfg.Webdav

	handler := &Handler{cfg}

	if wcfg.ListenAddr == "" {
		return errors.New("Webdav server listen addr must be configured.")
	}

	lis, err := net.Listen("tcp", wcfg.ListenAddr)
	if err != nil {
		return fmt.Errorf("Failed to listen \"%s\": %v", wcfg.ListenAddr, err)
	}

	closed := false
	if closeC != nil {
		go func() {
			<-closeC
			closed = true
			lis.Close()
		}()
	}

	cert, err := tls.LoadX509KeyPair(wcfg.CertFile, wcfg.KeyFile)
	if err != nil {
		return fmt.Errorf("Failed to load webdav {cert,key} pair: %v", err)
	}

	// Note: This doesn't enable h2. Reconsider this if there is a webdav client w/ h2 support.
	httpsrv := http.Server{
		Addr:    wcfg.ListenAddr,
		Handler: logger.HttpHandler(accesslog, logger.Info, handler),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"http/1.1"},
		},
	}
	tlis := tls.NewListener(lis, httpsrv.TLSConfig)

	if err := httpsrv.Serve(tlis); err != nil {
		if closed {
			// Suppress "use of closed network connection" error if we intentionally closed the listener.
			return nil
		}
		return err
	}
	return nil
}
