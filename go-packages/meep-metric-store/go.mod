module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metric-store

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/influxdata/influxdb1-client v0.0.0-20190809212627-fc22c7df067e
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
