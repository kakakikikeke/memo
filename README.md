# memo

![screen_shot.png](https://raw.githubusercontent.com/kakakikikeke/memo/master/images/screen_shot.png)

This is a sample application for docker-compose.  
The "memo" is easy to save your memory wherever you want.  
This app is powerd by beego and redis.

## localhost on standard golang
If you haven't installed redis-server yet, first you should install redis-server.  
The "memo" is connected to localhost:6379 redis-server.

Install libraries.

* go mod tidy

And build and run it.

* go fmt ./... && go build && ./memo

You can access with your browser to "http://localhost:8080/".

Update libraries.

* go get -u && go mod tidy

## localhost on Docker (docker-compose)
If you haven't installed docker and docker-compose yet, first you should install them.

* docker compose up -d

You can access with your browser to "http://localhost/".

## Test
* go test ./...
