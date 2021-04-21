
## metabase-quick

Purpose: Visualize local csv files via Metabase, without user login and permission check, without database setup.
Use it and forget it...

Metabase is wonderful! It helps me a lot. :) 

But I have to do a lot of work for setting up Metabase: create account, create database, including import data, if I don't have one at hand. Especially when I just want to do some tiny work of data visualization, it's really a pain.

## Design

### Cross-platform

I will build it using Go so that it can easily run on Windows, MacOS, Linux.

## WIP

- [x] import local csv files into local in-memory sql engine, and serving
- [ ] simplify metabase frontend, do a lot
- [ ] go server: mock metabase key api logic, the hardest... 

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
