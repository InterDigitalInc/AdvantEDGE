module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/lib/pq v1.5.2
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
