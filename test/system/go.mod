module github.com/InterDigitalInc/AdvantEDGE/test/system

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-wais-client v0.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.4
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../go-packages/meep-app-support-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client => ../../go-packages/meep-gis-engine-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-loc-serv-client => ../../go-packages/meep-loc-serv-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client => ../../go-packages/meep-platform-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-client => ../../go-packages/meep-rnis-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../go-packages/meep-service-mgmt-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-vis-client => ../../go-packages/meep-vis-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-wais-client => ../../go-packages/meep-wais-client
)
