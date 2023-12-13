createRepository:
	go run . createRepository --input_file="example/domain/user.go" --output_file="example/infrastructure/user_repository.go"
	goimports -w example/infrastructure/user_repository.go