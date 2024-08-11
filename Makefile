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
