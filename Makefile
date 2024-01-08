dep:
	go mod download
dep-sync:
	go mod tidy
proto-build:
	docker run --rm -v %cd%/external:/defs namely/protoc-all -d dto\proto\BPS -l go