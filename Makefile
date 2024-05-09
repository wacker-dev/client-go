gen:
	@curl -o proto/wacker.proto -s https://raw.githubusercontent.com/iawia002/wacker/main/wacker/proto/wacker.proto
	@protoc -I proto --go_out=paths=source_relative:internal --go_opt=Mwacker.proto="github.com/wacker-dev/client-go;internal" \
		--go-grpc_out=paths=source_relative:internal --go-grpc_opt=Mwacker.proto="github.com/wacker-dev/client-go;internal" wacker.proto
