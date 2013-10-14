require 'sinatra'
require 'redis'
require 'uri'

set :logging, false
set :port, 80
set :bind, '0.0.0.0'
set :environment, :production

redisUrl = URI(ENV['REDISTOGO_URL'])
redis = Redis.new(:driver => :hiredis,
                  :host => redisUrl.host,
                  :port => redisUrl.port,
                  :password => redisUrl.password)

get '/quote/:author' do |author|
    redisKey = 'quote:' + author
    content_type 'application/json'
    result = redis.srandmember(redisKey)
    halt 404 if result == nil
    result
end
