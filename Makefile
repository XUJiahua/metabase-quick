fmt:
	go fmt ./...
run:fmt
	go run . \
		dataset/sample-dataset/products.csv \
		dataset/sample-dataset/reviews.csv \
		dataset/sample-dataset/people.csv
dev:fmt
	go run . -d -v \
		dataset/sample-dataset/orders.csv \
		dataset/sample-dataset/products.csv \
		dataset/sample-dataset/reviews.csv \
		dataset/sample-dataset/people.csv
git_submodule:
	git submodule update --init

clean:
	git clean -f -x -d
build_frontend:
	cd metabase && git clean -f -x -d && yarn build
cp_frontend:build_frontend
	cp -r metabase/resources/frontend_client/* pkg/metabase/frontend_client
build_dir:
	mkdir -p build
build_mac:
	env GOOS=darwin GOARCH=amd64 go build -o build/metabase-quick-darwin
build_win:
	env GOOS=windows GOARCH=amd64 go build -o build/metabase-quick-win.exe
build_linux:
	env GOOS=linux GOARCH=amd64 go build -o build/metabase-quick-linux
build:clean cp_frontend build_dir build_mac build_win build_linux
	echo 'done'