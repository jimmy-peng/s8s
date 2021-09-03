package dynamiccertificates

import (
	 "bytes"
)

// certKeyContent holds the content for the cert and key
type certKeyContent struct {
	cert []byte
	key  []byte
}

func (c *certKeyContent) Equal(rhs *certKeyContent) bool {
	if c == nil || rhs == nil {
		return c == rhs
	}

	return bytes.Equal(c.key, rhs.key) && bytes.Equal(c.cert, rhs.cert)
}
