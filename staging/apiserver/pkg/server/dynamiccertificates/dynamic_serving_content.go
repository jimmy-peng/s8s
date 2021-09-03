package dynamiccertificates

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"sync/atomic"
)

// DynamicCertKeyPairContent provides a CertKeyContentProvider that can dynamically react to new file content
type DynamicCertKeyPairContent struct {
	name string

	// certFile is the name of the certificate file to read.
	certFile string
	// keyFile is the name of the key file to read.
	keyFile string

	// certKeyPair is a certKeyContent that contains the last read, non-zero length content of the key and cert
	certKeyPair atomic.Value

	//listeners []Listener

	// queue only ever has one item, but it has nice error handling backoff/retry semantics
	//queue workqueue.RateLimitingInterface
}

//var _ CertKeyContentProvider = &DynamicCertKeyPairContent{}
//var _ ControllerRunner = &DynamicCertKeyPairContent{}

// NewDynamicServingContentFromFiles returns a dynamic CertKeyContentProvider based on a cert and key filename
func NewDynamicServingContentFromFiles(purpose, certFile, keyFile string) (*DynamicCertKeyPairContent, error) {
	if len(certFile) == 0 || len(keyFile) == 0 {
		return nil, fmt.Errorf("missing filename for serving cert")
	}
	name := fmt.Sprintf("%s::%s::%s", purpose, certFile, keyFile)

	ret := &DynamicCertKeyPairContent{
		name:     name,
		certFile: certFile,
		keyFile:  keyFile,
		//queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), fmt.Sprintf("DynamicCABundle-%s", purpose)),
	}
	if err := ret.loadCertKeyPair(); err != nil {
		return nil, err
	}

	return ret, nil
}

// loadCertKeyPair determines the next set of content for the file.
func (c *DynamicCertKeyPairContent) loadCertKeyPair() error {
	cert, err := ioutil.ReadFile(c.certFile)
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(c.keyFile)
	if err != nil {
		return err
	}
	if len(cert) == 0 || len(key) == 0 {
		return fmt.Errorf("missing content for serving cert %q", c.Name())
	}

	// Ensure that the key matches the cert and both are valid
	_, err = tls.X509KeyPair(cert, key)
	if err != nil {
		return err
	}

	newCertKey := &certKeyContent{
		cert: cert,
		key:  key,
	}

	// check to see if we have a change. If the values are the same, do nothing.
	existing, ok := c.certKeyPair.Load().(*certKeyContent)
	if ok && existing != nil && existing.Equal(newCertKey) {
		return nil
	}

	c.certKeyPair.Store(newCertKey)
	/*
	klog.V(2).InfoS("Loaded a new cert/key pair", "name", c.Name())

	for _, listener := range c.listeners {
		listener.Enqueue()
	}
	*/

	return nil
}

// Name is just an identifier
func (c *DynamicCertKeyPairContent) Name() string {
	return c.name
}

// CurrentCertKeyContent provides cert and key byte content
func (c *DynamicCertKeyPairContent) CurrentCertKeyContent() ([]byte, []byte) {
	certKeyContent := c.certKeyPair.Load().(*certKeyContent)
	return certKeyContent.cert, certKeyContent.key
}