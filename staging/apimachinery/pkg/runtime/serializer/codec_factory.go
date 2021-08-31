package serializer

import "s8s/staging/apimachinery/pkg/runtime"

type serializerType struct {
	StreamContentType string
}

type CodecFactory struct {
	scheme *runtime.Scheme
}

type CodecFactoryOptions struct {
	Pretty bool
}

func NewCodecFactory(scheme *runtime.Scheme) CodecFactory {
	return CodecFactory{
		scheme: scheme,
	}
}