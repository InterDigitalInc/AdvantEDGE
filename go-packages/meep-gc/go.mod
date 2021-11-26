module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gc

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store v0.0.0
	github.com/influxdata/influxdb1-client v0.0.0-20191209144304-8bf82d3c094d
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.17.0 // indirect
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store => ../../go-packages/meep-sandbox-store
)
