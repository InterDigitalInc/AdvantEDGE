module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-virt-engine

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog v0.0.0
	github.com/gorilla/handlers v1.4.0 // indirect
	github.com/gorilla/mux v1.7.1 // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../go-packages/meep-ctrl-engine-model

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog => ../../go-packages/meep-watchdog
