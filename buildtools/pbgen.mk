.PHONY: pbgen pbgen-mm pbgen-game pbgen-web pbgen-gb

PROTO_DIR := ../proto
OUT_DIR := ../pb

pbgen: pbgen-mm pbgen-game pbgen-gb

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

pbgen-gb:
	protoc --proto_path=$(PROTO_DIR) \
		gamebox/v1/gamebox.proto \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative
