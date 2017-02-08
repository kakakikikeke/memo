require 'sinatra'
require 'redis'

rs = ENV['REDIS_SERVER'] || 'localhost'
rsp = ENV['REDIS_SERVER_PORT'] || '6379'
redis = Redis.new host: rs, port: rsp

get '/' do
  @memos = redis.lrange :memos, 0, -1
  erb :memo_list
end

post '/save' do
  msg = params[:msg]
  redis.lpush :memos, msg
  redirect '/'
end

post '/clear' do
  redis.del :memos
  redirect '/'
end
