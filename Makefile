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
build_frontend:
	cd metabase && yarn build
cp:
	cp -r metabase/resources/frontend_client/* pkg/metabase/frontend_client
build:fmt
	go build
build-win:
	env GOOS=windows GOARCH=amd64 go build
