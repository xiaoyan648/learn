.PHONY: types
types:
	cd types && protoc --proto_path=. --go_out=paths=source_relative:. ./*.proto 

.PHONY: grpc
grpc:
	protoc --proto_path=. --go_out=paths=source_relative:. ./*.proto 
