dep:
	go mod download
dep-sync:
	go mod tidy
proto-build:
	docker run --rm -v `pwd`/external:/defs namely/protoc-all -d proto/BPS -l go -o ./