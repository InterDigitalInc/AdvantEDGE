module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-websocket

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/gorilla/websocket v1.4.2
	github.com/stretchr/testify v1.5.1 // indirect
	golang.org/x/sys v0.0.0-20210423082822-04245dca01da // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
