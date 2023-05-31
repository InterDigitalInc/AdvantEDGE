module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm

go 1.16

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client v0.0.0-00010101000000-000000000000
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr v0.0.0-00010101000000-000000000000
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/prometheus/client_golang v1.9.0
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/bwm => ../../go-apps/meep-tm/server/bwm
	github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tm/server/mts => ../../go-apps/meep-tm/server/mts
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client => ../../go-packages/meep-app-support-client
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
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-ctrl-client => ../../go-packages/meep-sandbox-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client => ../../go-packages/meep-service-mgmt-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-subscriptions => ../../go-packages/meep-subscriptions
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-swagger-api-mgr => ../../go-packages/meep-swagger-api-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-websocket => ../../go-packages/meep-websocket
)
