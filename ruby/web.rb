require 'sinatra'
require 'redis'
require 'uri'

redisUrl = URI(ENV['REDISTOGO_URL'])
redis = Redis.new(:driver => :hiredis,
                  :host => redisUrl.host,
                  :port => redisUrl.port,
                  :password => redisUrl.password)

get '/quote/:author' do |author|
    redisKey = 'quote:' + author
    result = redis.srandmember(redisKey)
end
