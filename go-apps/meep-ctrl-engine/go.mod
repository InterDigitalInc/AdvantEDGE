module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-ctrl-engine

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-replay-manager v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog v0.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.3
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch => ../../go-packages/meep-couch
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client => ../../go-packages/meep-ctrl-engine-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../go-packages/meep-ctrl-engine-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store => ../../go-packages/meep-metric-store
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-replay-manager => ../../go-packages/meep-replay-manager
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog => ../../go-packages/meep-watchdog
)
