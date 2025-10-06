.PHONY: pbgen pbgen-mm pbgen-game pbgen-web

PROTO_DIR := ../proto
OUT_DIR := ../pb

pbgen: pbgen-mm pbgen-game pbgen-web

pbgen-mm:
	protoc --proto_path=$(PROTO_DIR) \
		matchmaker/v1/matchmaker.proto \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--ts_out=$(OUT_DIR)

pbgen-game:
	protoc --proto_path=$(PROTO_DIR) \
		game/v1/game.proto \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--ts_out=$(OUT_DIR)

# pbgen-web:
# 	protoc --proto_path=$(PROTO_DIR) \
# 		matchmaker/v1/matchmaker.proto \
# 		game/v1/game.proto \
# 		--js_out=import_style=commonjs,binary:$(OUT_DIR) \
# 		--grpc-web_out=import_style=typescript,mode=grpcwebtext:$(OUT_DIR)
