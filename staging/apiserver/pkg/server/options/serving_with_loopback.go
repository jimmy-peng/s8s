package options

import (
	"s8s/staging/apiserver/pkg/server"
)

type SecureServingOptionsWithLoopback struct {
	*SecureServingOptions
}

// ApplyTo fills up serving information in the server configuration.
func (s *SecureServingOptionsWithLoopback) ApplyTo(secureServingInfo **server.SecureServingInfo /*, loopbackClientConfig **rest.Config*/) error {
	if s == nil || s.SecureServingOptions == nil || secureServingInfo == nil {
		return nil
	}

	if err := s.SecureServingOptions.ApplyTo(secureServingInfo); err != nil {
		return err
	}

/*
		if *secureServingInfo == nil || loopbackClientConfig == nil {
			return nil
		}

		// create self-signed cert+key with the fake server.LoopbackClientServerNameOverride and
		// let the server return it when the loopback client connects.
		certPem, keyPem, err := certutil.GenerateSelfSignedCertKey(server.LoopbackClientServerNameOverride, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
		}
		certProvider, err := dynamiccertificates.NewStaticSNICertKeyContent("self-signed loopback", certPem, keyPem, server.LoopbackClientServerNameOverride)
		if err != nil {
			return fmt.Errorf("failed to generate self-signed certificate for loopback connection: %v", err)
		}

		// Write to the front of SNICerts so that this overrides any other certs with the same name
		(*secureServingInfo).SNICerts = append([]dynamiccertificates.SNICertKeyContentProvider{certProvider}, (*secureServingInfo).SNICerts...)

		secureLoopbackClientConfig, err := (*secureServingInfo).NewLoopbackClientConfig(uuid.New().String(), certPem)
		switch {
		// if we failed and there's no fallback loopback client config, we need to fail
		case err != nil && *loopbackClientConfig == nil:
			(*secureServingInfo).SNICerts = (*secureServingInfo).SNICerts[1:]
			return err

		// if we failed, but we already have a fallback loopback client config (usually insecure), allow it
		case err != nil && *loopbackClientConfig != nil:

		default:
			*loopbackClientConfig = secureLoopbackClientConfig
		}*/

	return nil
}
