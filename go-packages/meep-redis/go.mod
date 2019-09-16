module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/KromDaniel/jonson v0.0.0-20180630143114-d2f9c3c389db // indirect
	github.com/KromDaniel/rejonson v0.0.0-20180822072824-00b5bcf2b351
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
