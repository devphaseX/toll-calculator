obu:
	@go build -o ./bin/obu ./obu/
	@./bin/obu
.PHONY:obu

receiver:
	@go build -o ./bin/receiver ./data_receiver/
	@./bin/receiver
.PHONY:receiver

calculator:
	@go build -o ./bin/calculator ./distance_calculator/
	@./bin/calculator
.PHONY:calculator


aggregator:
	@go build -o ./bin/aggregator ./aggregator/
	@./bin/aggregator
.PHONY:aggregator


proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   types/gtypes.proto
.PHONY:proto
