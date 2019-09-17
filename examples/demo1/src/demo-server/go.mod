module github.com/InterDigitalInc/AdvantEDGE/demoserver

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/locservapi v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/mgmanagerapi v0.0.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/locservapi => ../../../../go-packages/meep-loc-serv-client
	github.com/InterDigitalInc/AdvantEDGE/mgmanagerapi => ../../../../go-packages/meep-mg-manager-client
)
