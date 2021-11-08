build:
	go build -o ./plugin/greeter ./plugin/greeter/greeter_implementation.go
	go build -o ./plugin/shard ./plugin/shard/shard.go
	go build -o basic .