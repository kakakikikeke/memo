# memo

![screen_shot.png](https://raw.githubusercontent.com/kakakikikeke/memo/master/images/screen_shot.png)

This is a sample application for docker-compose.  
The "memo" is easy to save your memory wherever you want.  
This app is powerd by beego and redis.

## How to build

### localhost on standard golang
If you haven't installed redis-server yet, first you should install redis-server.  
The "memo" is connected to localhost:6379 redis-server.

Install libraries.

* go get -u github.com/go-redis/redis
* go get github.com/astaxie/beego
* go get github.com/astaxie/beego/logs

And run it.

* go run memo.go

You can access with your browser to "http://localhost:8080/".

### localhost on Docker (docker-compose)
If you haven't installed docker and docker-compose yet, first you should install them.

* docker-compose up -d

You can access with your browser to "http://localhost/".

## Heroku Container Registry

```
heroku container:login
heroku create -a memo-app-12345
(heroku git:remote -a memo-app-12345)
docker build -f Dockerfile -t registry.heroku.com/memo-app-12345/web .
docker push registry.heroku.com/memo-app-12345/web
heroku addons:create heroku-redis:hobby-dev
(heroku config:set REDIS_URL=redis://user:pass@ec2-00-000-000-000.compute-1.amazonaws.com:12345)
heroku container:release web
heroku open -a memo-app-12345
```

You can acces with your browser to "https://memo-app-12345.herokuapp.com/".