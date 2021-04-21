fmt:
	go fmt ./...
run:fmt
	go run . dataset/iris.csv dataset/rb.csv
