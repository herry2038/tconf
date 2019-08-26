module github.com/toventang/tconf

go 1.12

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.44.4-0.20190822213831-0f42522d468a
	cloud.google.com/go/bigquery => github.com/googleapis/google-cloud-go/bigquery v0.44.4-0.20190822213831-0f42522d468a
	cloud.google.com/go/datastore => github.com/googleapis/google-cloud-go/datastore v0.44.4-0.20190822213831-0f42522d468a
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190731235908-ec7cb31e5a56
	golang.org/x/image => github.com/golang/image v0.0.0-20190823064033-3a9bac650e44
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190814143026-e8b3e6111d02
	golang.org/x/mod => github.com/golang/mod v0.1.0
	golang.org/x/net => github.com/golang/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190816200558-6889da9d5479
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
	//gopkg.in/jcmturner/gokrb5.v7 => github.com/jcmturner/gokrb5 v7.2.3+incompatible
	//google.golang.org/api => github.com/googleapis/google-api-go-client v0.0.0-20190821000710-329ecc3c9c34
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.9.0
	google.golang.org/appengine => github.com/golang/appengine v1.6.2-0.20190801181406-fb139bde60fa
	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20190819201941-24fa4b261c55
	google.golang.org/grpc => github.com/grpc/grpc-go v1.2.1-0.20190822205623-0574097d6793
	sigs.k8s.io/yaml v1.1.0 => github.com/kubernetes-sigs/yaml v1.1.0
)

require (
	github.com/spf13/viper v1.4.0
	go.etcd.io/etcd v0.0.0-20190823073701-67d0c21bb04c
)
