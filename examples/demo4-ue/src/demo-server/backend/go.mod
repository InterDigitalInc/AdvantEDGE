module github.com/AdvantEDGE/examples/demo4-ue/src/demo-server/backend

go 1.13

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0-20211214133749-f203f7ab4f1c
	github.com/antihax/optional v1.0.0
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/spf13/viper v1.12.0
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../../../../go-packages/meep-app-support-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-dai-client => ../../../../../go-packages/meep-dai-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../../../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger => ../../../../../go-packages/meep-http-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-client => ../../../../../go-packages/meep-loc-serv-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../../../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics => ../../../../../go-packages/meep-metrics
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../../../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../../../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../../../../go-packages/meep-service-mgmt-client
)
