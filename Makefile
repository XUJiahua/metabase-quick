fmt:
	go fmt ./...
run:fmt
	go run . dataset/ml-latest-small/movies.csv dataset/ml-latest-small/tags.csv
# use default;
# show tables;
# join test
# select * from movies left join tags on movies.movieId = tags.movieId limit 1;
sqlClient:
	mysql -h 127.0.0.1 -u root
