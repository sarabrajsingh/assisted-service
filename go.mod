module github.com/openshift/assisted-service

go 1.13

require (
	github.com/ajeddeloh/go-json v0.0.0-20200220154158-5ae607161559 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/aws/aws-sdk-go v1.32.6
	github.com/cavaliercoder/go-cpio v0.0.0-20180626203310-925f9528c45e
	github.com/containers/image/v5 v5.7.0
	github.com/coreos/ignition/v2 v2.6.0
	github.com/danielerez/go-dns-client v0.0.0-20200630114514-0b60d1703f0b
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/docker v1.13.1 // indirect
	github.com/filanov/stateswitch v0.0.0-20200714113403-51a42a34c604
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-openapi/errors v0.19.6
	github.com/go-openapi/loads v0.19.5
	github.com/go-openapi/runtime v0.19.20
	github.com/go-openapi/spec v0.19.8
	github.com/go-openapi/strfmt v0.19.5
	github.com/go-openapi/swag v0.19.9
	github.com/go-openapi/validate v0.19.10
	github.com/golang-collections/go-datastructures v0.0.0-20150211160725-59788d5eb259
	github.com/golang/mock v1.4.3
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kennygrant/sanitize v1.2.4
	github.com/metal3-io/baremetal-operator v0.0.0-00010101000000-000000000000
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/moby/moby v1.13.1
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/openshift-online/ocm-sdk-go v0.1.130
	github.com/openshift/api v0.0.0-20200326152221-912866ddb162
	github.com/openshift/client-go v0.0.0-20200422192633-6f6c07fc2a70
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.6.0
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.7.0
	github.com/slok/go-http-metrics v0.8.0
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/thoas/go-funk v0.6.0
	github.com/vincent-petithory/dataurl v0.0.0-20191104211930-d1553a71de50
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sync v0.0.0-20201008141435-b3e1573b7520
	golang.org/x/tools v0.0.0-20200914190812-8f9ed77dd8e5 // indirect
	gopkg.in/square/go-jose.v2 v2.3.1
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/sqlite v1.1.3
	gorm.io/gorm v1.20.6
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	k8s.io/api v0.18.5
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20200603063816-c1c6865ac451 // indirect
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/metal3-io/baremetal-operator => github.com/openshift/baremetal-operator v0.0.0-20200715132148-0f91f62a41fe // Use OpenShift fork
	k8s.io/client-go => k8s.io/client-go v0.18.2
)
