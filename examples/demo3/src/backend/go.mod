module github.com/InterDigitalInc/AdvantEDGE/example/demo3/src

go 1.12

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/roymx/viper v1.3.3-0.20190416163942-b9a223fc58a3
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client => ../../../../go-packages/meep-ams-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../../../go-packages/meep-app-support-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../../../go-packages/meep-service-mgmt-client
)
