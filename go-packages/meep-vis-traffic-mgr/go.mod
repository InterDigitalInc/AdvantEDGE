module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-vis-traffic-mgr

go 1.16

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/eclipse/paho.mqtt.golang v1.4.2
	github.com/lib/pq v1.10.7
	github.com/roymx/viper v1.3.3-0.20190416163942-b9a223fc58a3

)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
