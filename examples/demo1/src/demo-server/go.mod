module github.com/InterDigitalInc/AdvantEDGE/demoserver

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/locservapi v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/mgmanagerapi v0.0.0
	github.com/antihax/optional v0.0.0-20180407024304-ca021399b1a6 // indirect
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
)

replace github.com/InterDigitalInc/AdvantEDGE/mgmanagerapi => ../../../../go-packages/meep-mg-manager-client

replace github.com/InterDigitalInc/AdvantEDGE/locservapi => ../../../../go-packages/meep-loc-serv-client
