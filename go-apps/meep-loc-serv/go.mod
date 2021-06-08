module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-loc-serv

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache v0.0.0
        github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.8.0
	github.com/prometheus/client_golang v1.9.0
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model => ../../go-packages/meep-data-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-cache => ../../go-packages/meep-gis-cache
        github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-gis-engine-client => ../../go-packages/meep-gis-engine-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger => ../../go-packages/meep-http-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics => ../../go-packages/meep-metrics
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq => ../../go-packages/meep-mq
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
)
