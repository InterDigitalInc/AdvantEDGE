module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-asset-mgr

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-postgis v0.0.0-20200703133018-94138d8210a3 // indirect
	github.com/lib/pq v1.5.2
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
