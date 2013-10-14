import os
import redis
from urlparse import urlparse
from bottle import route, run, response, abort

redisUrl = urlparse(os.environ['REDISTOGO_URL'])
pool = redis.ConnectionPool(host=redisUrl.hostname, port=redisUrl.port, db=0, password=redisUrl.password)

@route('/quote/<author>')
def get_quote(author):
    redisClient = redis.Redis(connection_pool=pool)
    redisKey = 'quote:' + author
    result = redisClient.srandmember(redisKey)
    response.content_type = 'application/json'
    if result:
        return result
    else:
        abort(404)

run(host='localhost', port=8080, debug=False)
