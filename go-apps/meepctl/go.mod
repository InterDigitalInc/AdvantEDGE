module github.com/InterDigitalInc/AdvantEDGE/go-apps/meepctl

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client v0.0.0
	github.com/antihax/optional v1.0.0 // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/roymx/viper v1.3.3-0.20190416163942-b9a223fc58a3
	github.com/spf13/cobra v0.0.0-20190109003409-7547e83b2d85
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-client => ../../go-packages/meep-ctrl-engine-client
