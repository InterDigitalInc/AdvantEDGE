module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store v0.0.0-20200306214341-11d08c83c4d6 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
)
