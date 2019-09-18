module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-metrics-engine

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/google/go-cmp v0.3.1 // indirect
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/olivere/elastic v6.2.21+incompatible
	github.com/pkg/errors v0.8.1 // indirect
	golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
