package legacyscheme

import (
	"s8s/staging/apimachinery/pkg/runtime"
	"s8s/staging/apimachinery/pkg/runtime/serializer"
)

var (
	Scheme = runtime.NewScheme()
	Codecs = serializer.NewCodecFactory(Scheme)
	//ParamterCodec = runtime.NewParameterCodec(Scheme)
)