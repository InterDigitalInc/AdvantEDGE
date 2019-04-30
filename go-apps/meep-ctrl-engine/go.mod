module github.com/InterDigitalInc/AdvantEDGE/go-apps/meep-ctrl-engine

go 1.12

require (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-tc-engine-client v0.0.0
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-virt-engine-client v0.0.0
	github.com/KromDaniel/jonson v0.0.0-20180630143114-d2f9c3c389db // indirect
	github.com/KromDaniel/rejonson v0.0.0-20180822072824-00b5bcf2b351
	github.com/flimzy/diff v0.1.5 // indirect
	github.com/flimzy/kivik v1.8.1
	github.com/flimzy/testy v0.1.16 // indirect
	github.com/go-kivik/couchdb v1.8.1
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190411002643-bd77b112433e // indirect
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.1
	github.com/gregjones/httpcache v0.0.0-20190212212710-3befbb6ad0cc // indirect
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/json-iterator/go v1.1.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/spf13/pflag v1.0.3 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
	k8s.io/api v0.0.0-20181204000039-89a74a8d264d // indirect
	k8s.io/apimachinery v0.0.0-20181127025237-2b1284ed4c93 // indirect
	k8s.io/client-go v10.0.0+incompatible // indirect
	k8s.io/klog v0.0.0-20181108234604-8139d8cb77af // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

replace (
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-tc-engine-client => ../../go-packages/meep-tc-engine-client
	github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-virt-engine-client => ../../go-packages/meep-virt-engine-client
)
