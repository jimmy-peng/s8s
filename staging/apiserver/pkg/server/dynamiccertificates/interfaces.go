package dynamiccertificates

// CertKeyContentProvider provides a certificate and matching private key.
type CertKeyContentProvider interface {
	//Notifier

	// Name is just an identifier.
	Name() string
	// CurrentCertKeyContent provides cert and key byte content.
	CurrentCertKeyContent() ([]byte, []byte)
}
