ifeq ($(OS), Windows_NT)
	BIN := nano-cli.exe
	XK6_BIN := k6.exe
	MKFOLDER := if not exist "build" mkdir build
	GREP_CMD := findstr /V
else
	BIN := nano-cli
	XK6_BIN := k6
	MKFOLDER := mkdir -p build
	GREP_CMD := grep -v
endif

TESTABLE_PACKAGES = `go list ./... | $(GREP_CMD) examples | $(GREP_CMD) constants | $(GREP_CMD) mocks | $(GREP_CMD) helpers | $(GREP_CMD) interfaces | $(GREP_CMD) protos | $(GREP_CMD) e2e | $(GREP_CMD) benchmark`

setup: init-submodules
	@go get ./...

build:
	@$(MKFOLDER)
	@go build -o build/$(BIN) .
	@echo "build nano-cli at ./build/$(BIN)"

init-submodules:
	@git submodule update --init --recursive

setup-ci:
	@go install github.com/mattn/goveralls@latest
	@go install github.com/wadey/gocovmerge@latest

setup-protobuf-macos:
	@brew install protobuf
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

run-jaeger-aio:
	@docker compose -f ./examples/testing/docker-compose-jaeger.yml up -d
	@echo "Access jaeger UI @ http://localhost:16686"

run-chat-example:
	@cd examples/testing && docker compose up -d etcd nats && cd ../demo/chat/ && go run main.go

run-demo-cluster-frontend:
	@go run examples/demo/cluster/main.go --port 3250 --type=conn --frontend=true

run-demo-cluster-backend:
	@go run examples/demo/cluster/main.go --port 3251 --type=room --frontend=false

run-cluster-grpc-example-connector:
	@cd examples/demo/cluster_grpc && go run main.go

run-cluster-grpc-example-room:
	@cd examples/demo/cluster_grpc && go run main.go --port 3251 --rpcsvport 3435 --type room --frontend=false

run-cluster-worker-example-room:
	@cd examples/demo/worker && go run main.go --type room --frontend=true

run-cluster-worker-example-metagame:
	@cd examples/demo/worker && go run main.go --type metagame --frontend=false

run-cluster-worker-example-worker:
	@cd examples/demo/worker && go run main.go --type worker --frontend=false

run-custom-metrics-example:
	@cd examples/demo/custom_metrics && go run main.go --port 3250

run-rate-limiting-example:
	@go run examples/demo/rate_limiting/main.go

protos-test-compile:
	@protoc -I benchmark/testdata/ benchmark/testdata/test.proto --go_out=benchmark/testdata --go-grpc_out=benchmark/testdata
	@protoc -I protos/protos/test protos/protos/test/*.proto --go_out=protos/test

protos-compile:
	@protoc -I protos/protos/ protos/protos/*.proto --go_out=protos --go-grpc_out=protos


rm-test-temp-files:
	@rm -f cluster/127.0.0.1* 127.0.0.1*
	@rm -f cluster/localhost* localhost*

ensure-testing-bin:
	@[ -f ./examples/testing/server ] || go build -o ./examples/testing/server ./examples/testing/main.go

ensure-testing-deps:
	@cd ./examples/testing && docker compose up -d

ensure-e2e-deps-grpc:
	@cd ./examples/testing && docker compose up -d etcd

kill-testing-deps:
	@cd ./examples/testing && docker compose down; true

kill-jaeger:
	@docker compose -f ./examples/testing/docker-compose-jaeger.yml down; true

e2e-test: e2e-test-nats e2e-test-grpc

e2e-test-nats: ensure-testing-deps ensure-testing-bin
	@echo "===============RUNNING E2E NATS TESTS==============="
	@go test ./e2e/e2e_test.go -update

e2e-test-grpc: ensure-testing-deps ensure-testing-bin
	@echo "===============RUNNING E2E GRPC TESTS==============="
	@go test ./e2e/e2e_test.go -update -grpc

bench-nats-sv:
	@NANO_METRICS_PROMETHEUS_PORT=9098 ./examples/testing/server -type game -frontend=false > /dev/null 2>&1 & echo $$! > back.PID
	@NANO_METRICS_PROMETHEUS_PORT=9099 ./examples/testing/server -type connector -frontend=true > /dev/null 2>&1 & echo $$! > front.PID

bench-grpc-sv:
	@NANO_METRICS_PROMETHEUS_PORT=9098 ./examples/testing/server -grpc -grpcport=3435 -type game -frontend=false > /dev/null 2>&1 & echo $$! > back.PID
	@NANO_METRICS_PROMETHEUS_PORT=9099 ./examples/testing/server -grpc -grpcport=3436 -type connector -frontend=true > /dev/null 2>&1 & echo $$! > front.PID

benchmark-test-nats: ensure-testing-deps ensure-testing-bin
	@echo "===============RUNNING BENCHMARK TESTS WITH NATS==============="
	@echo "--- starting testing servers"
	@echo "--- sleeping for 5 seconds"
	@make bench-nats-sv
	@sleep 5
	@go test ./benchmark/benchmark_test.go -bench=.
	@echo "--- killing testing servers"
	@kill `cat back.PID` && rm back.PID
	@kill `cat front.PID` && rm front.PID

benchmark-test-grpc: ensure-e2e-deps-grpc ensure-testing-bin
	@echo "===============RUNNING BENCHMARK TESTS WITH GRPC==============="
	@echo "--- starting testing servers"
	@echo "--- sleeping for 5 seconds"
	@make bench-grpc-sv
	@sleep 5
	@go test ./benchmark/benchmark_test.go -bench=.
	@echo "--- killing testing servers"
	@kill `cat back.PID` && rm back.PID
	@kill `cat front.PID` && rm front.PID

unit-test-coverage: kill-testing-deps
	@echo "===============RUNNING UNIT TESTS==============="
	@go test $(TESTABLE_PACKAGES) -coverprofile coverprofile.out

test: kill-testing-deps test-coverage
	@make rm-test-temp-files
	@make ensure-testing-deps
	@sleep 10
	@make e2e-test

test-coverage: unit-test-coverage
	@make rm-test-temp-files

test-coverage-html: test-coverage
	@go tool cover -html=coverprofile.out

merge-profiles:
	@rm -f coverage-all.out
	@gocovmerge *.out > coverage-all.out

test-coverage-func coverage-func: test-coverage merge-profiles
	@echo
	@echo "\033[1;34m=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\033[0m"
	@echo "\033[1;34mFunctions NOT COVERED by Tests\033[0m"
	@echo "\033[1;34m=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-\033[0m"
	@go tool cover -func=coverage-all.out | egrep -v "100.0[%]"

mocks: agent-mock session-mock networkentity-mock nano-mock serializer-mock metrics-mock acceptor-mock

agent-mock:
	@mockgen github.com/nut-game/nano/agent Agent,AgentFactory | sed 's/mock_agent/mocks/' > pkg/agent/mocks/agent.go

session-mock:
	@mockgen github.com/nut-game/nano/session Session,SessionPool | sed 's/mock_session/mocks/' > pkg/session/mocks/session.go

networkentity-mock:
	@mockgen github.com/nut-game/nano/networkentity NetworkEntity | sed 's/mock_networkentity/mocks/' > pkg/networkentity/mocks/networkentity.go

nano-mock:
	@mockgen github.com/nut-game/nano Nano | sed 's/mock_pkg/mocks/' > pkg/mocks/app.go

metrics-mock:
	@mockgen github.com/nut-game/nano/metrics Reporter | sed 's/mock_metrics/mocks/' > pkg/metrics/mocks/reporter.go
	@mockgen github.com/nut-game/nano/metrics Client | sed 's/mock_metrics/mocks/' > pkg/metrics/mocks/statsd_reporter.go

serializer-mock:
	@mockgen github.com/nut-game/nano/serialize Serializer | sed 's/mock_serialize/mocks/' > pkg/serialize/mocks/serializer.go

acceptor-mock:
	@mockgen github.com/nut-game/nano/acceptor PlayerConn,Acceptor | sed 's/mock_acceptor/mocks/' > pkg/mocks/acceptor.go

worker-mock:
	@mockgen github.com/nut-game/nano/worker RPCJob | sed 's/mock_worker/mocks/' > pkg/worker/mocks/rpc_job.go
