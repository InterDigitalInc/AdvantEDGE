module github.com/InterDigitalInc/AdvantEDGE/example/demo4/src/onboarded-demo

go 1.13

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/antihax/optional v1.0.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.3
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../../../go-packages/meep-app-support-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-mgr => ../../../../go-packages/meep-dai-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../../../go-packages/meep-service-mgmt-client
)
