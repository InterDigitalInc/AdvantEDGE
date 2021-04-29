module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-pdu-session-store

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model => ../../go-packages/meep-data-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
)
