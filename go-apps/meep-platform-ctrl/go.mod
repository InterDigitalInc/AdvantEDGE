module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-platform-ctrl

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog v0.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.3
	github.com/rs/xid v1.2.1 // indirect
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-couch => ../../go-packages/meep-couch
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model => ../../go-packages/meep-data-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq => ../../go-packages/meep-mq
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store => ../../go-packages/meep-sandbox-store
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog => ../../go-packages/meep-watchdog
)
