module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-mg-manager

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-app-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model v0.0.0
	github.com/KromDaniel/jonson v0.0.0-20180630143114-d2f9c3c389db // indirect
	github.com/KromDaniel/rejonson v0.0.0-20180822072824-00b5bcf2b351
	github.com/RyanCarrier/dijkstra v0.0.0-20180928224145-3fe1cac16289
	github.com/RyanCarrier/dijkstra-1 v0.0.0-20170512020943-0e5801a26345 // indirect
	github.com/albertorestifo/dijkstra v0.0.0-20160910063646-aba76f725f72 // indirect
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/mattomatic/dijkstra v0.0.0-20130617153013-6f6d134eb237 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/sirupsen/logrus v1.4.1
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../../../../ctrl-engine/model/go/src/meep-ctrl-engine-model
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-app-client => ../../../../../mg-manager/client/go/src/meep-mg-app-api-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mg-manager-model => ../../../../../mg-manager/model/go/src/meep-mg-manager-model
)
