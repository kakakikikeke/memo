# memo

This is a sample application for docker-compose.  
The "memo" is easy to save your memory wherever you want.  
This app is powerd by sinatra and redis.

## How to build

### localhost on standard ruby
If you haven't installed redis-server yet, first you should install redis-server.  
The "memo" is connected to localhost:6379 redis-server.

* bundle install && bundle exec rackup config.ru -o 0.0.0.0

You can access with your browser to "http://localhost:9292/".

### localhost on Docker (docker-compose)
If you haven't installed docker and docker-compose yet, first you should install them.

* docker-compose up -d

You can access with your browser to "http://localhost/".

### Heroku Container Registry
