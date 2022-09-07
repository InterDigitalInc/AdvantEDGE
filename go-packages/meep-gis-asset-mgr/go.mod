module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-asset-mgr

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0-20211214133749-f203f7ab4f1c
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/lib/pq v1.5.2
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
