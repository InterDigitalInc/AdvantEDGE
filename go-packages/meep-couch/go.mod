module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/flimzy/kivik v1.8.1
	github.com/go-kivik/couchdb v1.8.1
	github.com/go-kivik/kivik v1.8.1
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
