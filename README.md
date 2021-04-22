
## metabase-quick

Purpose: Visualize local csv files via Metabase, without user login and permission check, without database setup.
Use it and forget it...

Metabase is wonderful! It helps me a lot. :) 

But I have to do a lot of work for setting up Metabase: create account, create database, including import data, if I don't have one at hand. Especially when I just want to do some tiny work of data visualization, it's really a pain.

## Design

### Cross-platform

I will build it using Go so that it can easily run on Windows, MacOS, Linux.

### Built-in SQL DB engine

No need to setup a standalone database, e.g., MySQL.

### User Interface

Native query page

It's simple and useful enough, so that I only need mock 2 apis. See below.

### Metabase API Mock

1. /api/database
2. /api/dataset

## WIP

- [x] import local csv files into local in-memory sql engine tables, and serving, based on project https://github.com/dolthub/go-mysql-server
- [x] simplify metabase frontend, currently use native query as home page
- [x] go server: mock metabase key api logic
- [ ] more metabase api support, so that reuse more metabase frontend

### TODO

- [ ] infer more column type， now based on project https://github.com/go-gota/gota with limited column types supported
- [ ] db engine for handling bigger dataset


## Dev Setup

### Metabase dev setup

https://www.metabase.com/docs/latest/developers-guide.html

server:

```
lein run
```

visit:
metabase-server:3000

frontend:

```
yarn build-hot
```

visit:
metabase-server:3000 -> webpack-server:8080

### dev setup

```
func ReverseProxy() gin.HandlerFunc {
	target := "localhost:3000"

	return func(c *gin.Context) {
		director := func(req *http.Request) {
			req.URL.Scheme = "http"
			req.URL.Host = target
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// default route
router.NoRoute(ReverseProxy())
```

visit frontend:
go-server:8000 -> metabase-server:3000 

visit backend:
go-server:8000

## Build

in Metabase code repo 

```
yarn build
```

copy whole resources/frontend_client into pkg/metabase/frontend_client (metabase-quick code repo)


replace pkg/metabase/frontend_client/index.html with index.html (metabase backend generated)

Build metabase-quick (package frontend_client folder into go binary file)

```
go build
```

