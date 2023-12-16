createRepository:
	go run . createRepository --input_file="examples/domain/user.go" --output_file="examples/infrastructure/user_repository.go"
	goimports -w examples/infrastructure/user_repository.go