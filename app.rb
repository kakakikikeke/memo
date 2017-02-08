require 'sinatra'
require 'redis'

url = ENV['REDIS_URL'] || 'redis://localhost:6379'
redis = Redis.new :url => url

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
