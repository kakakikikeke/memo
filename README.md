# memo

This is a sample application for docker-compose.  
The "memo" is easy to save your memory wherever you want.  
This app is powerd by beego and redis.

## How to build

### localhost on standard golang
If you haven't installed redis-server yet, first you should install redis-server.  
The "memo" is connected to localhost:6379 redis-server.

* go run memo.go

You can access with your browser to "http://localhost:8080/".

### localhost on Docker (docker-compose)
If you haven't installed docker and docker-compose yet, first you should install them.

* docker-compose up -d

You can access with your browser to "http://localhost/".

### Heroku Container Registry
* git clone https://github.com/kakakikikeke/memo.git
* cd memo
* git checkout -b ver_golang origin/ver_golang
* heroku create -a memo-app-12345
  * or `heroku git:remote --app memo-app-12345`
* heroku addons:create heroku-redis:hobby-dev
* heroku config:set REDIS_URL=redis://user:pass@ec2-00-000-000-000.compute-1.amazonaws.com:12345
* heroku container:push web
* heroku open
