module github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-watchdog

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/KromDaniel/rejonson v0.0.0-20180822072824-00b5bcf2b351
	github.com/go-redis/redis v6.15.2+incompatible
)

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger

replace github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
