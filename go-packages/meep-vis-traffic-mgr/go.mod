module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-vis-traffic-mgr

go 1.16

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/lib/pq v1.5.2
	github.com/spf13/viper v1.3.2
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
