fmt:
	go fmt ./...
run:fmt
	PORT=8000 go run . \
		dataset/sample-dataset/orders.csv \
		dataset/sample-dataset/products.csv \
		dataset/sample-dataset/reviews.csv \
		dataset/sample-dataset/people.csv
# golang project hot reload, using gin
# go get github.com/codegangsta/gin
run-hot:fmt
	gin -p 8000 run main.go \
		dataset/sample-dataset/orders.csv \
		dataset/sample-dataset/products.csv \
		dataset/sample-dataset/reviews.csv \
		dataset/sample-dataset/people.csv
# use default;
# show tables;
# join test
# select * from movies left join tags on movies.movieId = tags.movieId limit 1;
sqlClient:
	mysql -h 127.0.0.1 -u root
