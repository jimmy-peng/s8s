package dynamiccertificates

// DynamicFileSNIContent provides a SNICertKeyContentProvider that can dynamically react to new file content
type DynamicFileSNIContent struct {
	*DynamicCertKeyPairContent
	sniNames []string
}

//var _ SNICertKeyContentProvider = &DynamicFileSNIContent{}
//var _ ControllerRunner = &DynamicFileSNIContent{}

// NewDynamicSNIContentFromFiles returns a dynamic SNICertKeyContentProvider based on a cert and key filename and explicit names
func NewDynamicSNIContentFromFiles(purpose, certFile, keyFile string, sniNames ...string) (*DynamicFileSNIContent, error) {
	servingContent, err := NewDynamicServingContentFromFiles(purpose, certFile, keyFile)
	if err != nil {
		return nil, err
	}

	ret := &DynamicFileSNIContent{
		DynamicCertKeyPairContent: servingContent,
		sniNames:                  sniNames,
	}
	if err := ret.loadCertKeyPair(); err != nil {
		return nil, err
	}

	return ret, nil
}

// SNINames returns explicitly set SNI names for the certificate. These are not dynamic.
func (c *DynamicFileSNIContent) SNINames() []string {
	return c.sniNames
}