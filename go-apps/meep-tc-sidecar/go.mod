module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-tc-sidecar

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis v0.0.0
	github.com/coreos/go-iptables v0.6.0
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344
	gopkg.in/inf.v0 v0.9.1 // indirect
	k8s.io/api v0.0.0-20181204000039-89a74a8d264d // indirect
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93 // indirect
	k8s.io/klog v0.0.0-20181108234604-8139d8cb77af // indirect
	k8s.io/kubernetes v1.13.4
	k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-data-key-mgr => ../../go-packages/meep-data-key-mgr
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger => ../../go-packages/meep-logger
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics => ../../go-packages/meep-metrics
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-mq => ../../go-packages/meep-mq
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-redis => ../../go-packages/meep-redis
)
