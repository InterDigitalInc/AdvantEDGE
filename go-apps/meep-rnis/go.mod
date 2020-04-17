module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-rnis

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store v0.0.0-20200306214341-11d08c83c4d6
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-notification-client v0.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.3
	golang.org/x/sys v0.0.0-20190412213103-97732733099d // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../go-packages/meep-ctrl-engine-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store => ../../go-packages/meep-metric-store
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-rnis-notification-client => ../../go-packages/meep-rnis-notification-client
)
