module httpsvrtrace

go 1.18

require (
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pp861/go-stdlib v1.0.1
	github.com/pp861/httptrace v1.0.0
)

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.9.0 // indirect
)

replace github.com/pp861/go-stdlib => ../../../go-stdlib