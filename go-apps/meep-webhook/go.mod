module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-webhook

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model v0.0.0
	github.com/ghodss/yaml v1.0.0
	k8s.io/api v0.0.0-20190430012547-97d6bb8ea5f4
	k8s.io/apimachinery v0.0.0-20190430211124-5bae42371a56
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ctrl-engine-model => ../../go-packages/meep-ctrl-engine-model

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-model => ../../go-packages/meep-model

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
