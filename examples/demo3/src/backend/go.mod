module github.com/InterDigitalInc/AdvantEDGE/example/demo3/src

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/spf13/viper v1.8.1
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client => ../../../../go-packages/meep-ams-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../../../go-packages/meep-app-support-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../../../go-packages/meep-service-mgmt-client
)
