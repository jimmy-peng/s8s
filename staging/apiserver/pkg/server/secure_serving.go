package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

// tlsConfig produces the tls.Config to serve with.
func (s *SecureServingInfo) tlsConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		// Can't use SSLv3 because of POODLE and BEAST
		// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
		// Can't use TLSv1.1 because of RC4 cipher usage
		MinVersion: tls.VersionTLS12,
		// enable HTTP2 for go's 1.7 HTTP Server
		NextProtos: []string{"h2", "http/1.1"},
	}

	// these are static aspects of the tls.Config
	if s.DisableHTTP2 {
		//klog.Info("Forcing use of http/1.1 only")
		tlsConfig.NextProtos = []string{"http/1.1"}
	}
	/*
		if s.MinTLSVersion > 0 {
			tlsConfig.MinVersion = s.MinTLSVersion
		}
		if len(s.CipherSuites) > 0 {
			tlsConfig.CipherSuites = s.CipherSuites
			insecureCiphers := flag.InsecureTLSCiphers()
			for i := 0; i < len(s.CipherSuites); i++ {
				for cipherName, cipherID := range insecureCiphers {
					if s.CipherSuites[i] == cipherID {
						klog.Warningf("Use of insecure cipher '%s' detected.", cipherName)
					}
				}
			}
		}

		if s.ClientCA != nil {
			// Populate PeerCertificates in requests, but don't reject connections without certificates
			// This allows certificates to be validated by authenticators, while still allowing other auth types
			tlsConfig.ClientAuth = tls.RequestClientCert
		}


		if s.ClientCA != nil || s.Cert != nil || len(s.SNICerts) > 0 {
			dynamicCertificateController := dynamiccertificates.NewDynamicServingCertificateController(
				tlsConfig,
				s.ClientCA,
				s.Cert,
				s.SNICerts,
				nil, // TODO see how to plumb an event recorder down in here. For now this results in simply klog messages.
			)

			if s.ClientCA != nil {
				s.ClientCA.AddListener(dynamicCertificateController)
			}
			if s.Cert != nil {
				s.Cert.AddListener(dynamicCertificateController)
			}

			// start controllers if possible
			if controller, ok := s.ClientCA.(dynamiccertificates.ControllerRunner); ok {
				// runonce to try to prime data.  If this fails, it's ok because we fail closed.
				// Files are required to be populated already, so this is for convenience.
				if err := controller.RunOnce(); err != nil {
					klog.Warningf("Initial population of client CA failed: %v", err)
				}

				go controller.Run(1, stopCh)
			}
			if controller, ok := s.Cert.(dynamiccertificates.ControllerRunner); ok {
				// runonce to try to prime data.  If this fails, it's ok because we fail closed.
				// Files are required to be populated already, so this is for convenience.
				if err := controller.RunOnce(); err != nil {
					klog.Warningf("Initial population of default serving certificate failed: %v", err)
				}

				go controller.Run(1, stopCh)
			}
			for _, sniCert := range s.SNICerts {
				sniCert.AddListener(dynamicCertificateController)
				if controller, ok := sniCert.(dynamiccertificates.ControllerRunner); ok {
					// runonce to try to prime data.  If this fails, it's ok because we fail closed.
					// Files are required to be populated already, so this is for convenience.
					if err := controller.RunOnce(); err != nil {
						klog.Warningf("Initial population of SNI serving certificate failed: %v", err)
					}

					go controller.Run(1, stopCh)
				}
			}

			// runonce to try to prime data.  If this fails, it's ok because we fail closed.
			// Files are required to be populated already, so this is for convenience.
			if err := dynamicCertificateController.RunOnce(); err != nil {
				klog.Warningf("Initial population of dynamic certificates failed: %v", err)
			}
			go dynamicCertificateController.Run(1, stopCh)

			tlsConfig.GetConfigForClient = dynamicCertificateController.GetConfigForClient
		}*/

	return tlsConfig, nil
}

// ServeWithListenerStopped runs the secure http server. It fails only if certificates cannot be loaded or the initial listen call fails.
// The actual server loop (stoppable by closing stopCh) runs in a go routine, i.e. ServeWithListenerStopped does not block.
// It returns a stoppedCh that is closed when all non-hijacked active requests have been processed.
// It returns a listenerStoppedCh that is closed when the underlying http Server has stopped listening.
// TODO: do a follow up PR to remove this function, change 'Serve' to return listenerStoppedCh
//  and update all components that call 'Serve'
func (s *SecureServingInfo) ServeWithListenerStopped(handler http.Handler) error {
	if s.Listener == nil {
		return fmt.Errorf("listener must not be nil")
	}

	tlsConfig, err := s.tlsConfig()
	if err != nil {
		return err
	}

	secureServer := &http.Server{
		Addr:           s.Listener.Addr().String(),
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,

		IdleTimeout:       90 * time.Second, // matches http.DefaultTransport keep-alive timeout
		ReadHeaderTimeout: 32 * time.Second, // just shy of requestTimeoutUpperBound
	}

	// At least 99% of serialized resources in surveyed clusters were smaller than 256kb.
	// This should be big enough to accommodate most API POST requests in a single frame,
	// and small enough to allow a per connection buffer of this size multiplied by `MaxConcurrentStreams`.
	const resourceBody99Percentile = 256 * 1024

	http2Options := &http2.Server{
		IdleTimeout: 90 * time.Second, // matches http.DefaultTransport keep-alive timeout
	}

	// shrink the per-stream buffer and max framesize from the 1MB default while still accommodating most API POST requests in a single frame
	http2Options.MaxUploadBufferPerStream = resourceBody99Percentile
	http2Options.MaxReadFrameSize = resourceBody99Percentile

	// use the overridden concurrent streams setting or make the default of 250 explicit so we can size MaxUploadBufferPerConnection appropriately
	if s.HTTP2MaxStreamsPerConnection > 0 {
		http2Options.MaxConcurrentStreams = uint32(s.HTTP2MaxStreamsPerConnection)
	} else {
		http2Options.MaxConcurrentStreams = 250
	}

	// increase the connection buffer size from the 1MB default to handle the specified number of concurrent streams
	http2Options.MaxUploadBufferPerConnection = http2Options.MaxUploadBufferPerStream * int32(http2Options.MaxConcurrentStreams)

	if !s.DisableHTTP2 {
		// apply settings to the server
		if err := http2.ConfigureServer(secureServer, http2Options); err != nil {
			return fmt.Errorf("error configuring http2: %v", err)
		}
	}

	// use tlsHandshakeErrorWriter to handle messages of tls handshake error
	/*
		tlsErrorWriter := &tlsHandshakeErrorWriter{os.Stderr}
		tlsErrorLogger := log.New(tlsErrorWriter, "", 0)
		secureServer.ErrorLog = tlsErrorLogger
	*/

	//klog.Infof("Serving securely on %s", secureServer.Addr)
	return RunServer(secureServer, s.Listener)
}

// TODO: make private when insecure serving is gone from the kube-apiserver
func RunServer(
	server *http.Server,
	ln net.Listener) error {
	if ln == nil {
		return fmt.Errorf("listener must not be nil")
	}

	//go func() {

		var listener net.Listener
		listener = tcpKeepAliveListener{ln}
		/*
		if server.TLSConfig != nil {
			listener = tls.NewListener(listener, server.TLSConfig)
		}
		*/

		//server.Serve(listener)
		server.ServeTLS(listener, "kubernetes.pem", "kubernetes-key.pem")
		/*
			msg := fmt.Sprintf("Stopped listening on %s", ln.Addr().String())
			select {
			case <-stopCh:
				klog.Info(msg)
			default:
				panic(fmt.Sprintf("%s due to error: %v", msg, err))
			}
		*/
	//}()

	return nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Copied from Go 1.7.2 net/http/server.go
type tcpKeepAliveListener struct {
	net.Listener
}
