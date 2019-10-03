module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/flimzy/kivik v1.8.1
	github.com/go-kivik/couchdb v1.8.1
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/pkg/errors v0.8.1 // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../go-packages/meep-ctrl-engine-model
