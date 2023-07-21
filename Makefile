
.PHONY gen-grpc:
gen-grpc:
	cd gen/server-grpc; protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative server.proto

.PHONY gen-oapi:
gen-oapi:
	cd gen/server-oapi; oapi-codegen -config codegen.cfg.yml server-oapi.yaml > server-oapi.gen.go

.PHONY gen-ui:
gen-ui:
	cd ui/vue-app; yarn run build

.PHONY dev-run: gen-grpc gen-oapi gen-ui
dev-run:
	go run cmd/app/*.go --log-level 0 server --http-debug
