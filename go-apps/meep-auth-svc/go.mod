module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-auth-svc

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store v0.0.0 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-users v0.0.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/lkysow/go-gitlab v0.7.1
	github.com/prometheus/client_golang v1.9.0
	github.com/roymx/viper v1.3.3-0.20190416163942-b9a223fc58a3
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-model => ../../go-packages/meep-data-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics => ../../go-packages/meep-metrics
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq => ../../go-packages/meep-mq
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-platform-ctrl-client => ../../go-packages/meep-platform-ctrl-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sandbox-store => ../../go-packages/meep-sandbox-store
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-sessions => ../../go-packages/meep-sessions
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-users => ../../go-packages/meep-users
)
