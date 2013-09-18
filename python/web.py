# gunicorn -w 4 -b 0.0.0.0:80 uni:app
import os
import re
import redis
from urlparse import urlparse
from pprint import pprint

redisUrl = urlparse(os.environ['REDISTOGO_URL'])
pool = redis.ConnectionPool(host=redisUrl.hostname, port=redisUrl.port, db=0, password=redisUrl.password)

def app(environ, start_response):
        expression = re.compile('/quote/(\w+)')
        match = expression.findall(environ['PATH_INFO'])
        if match:
                author = match[0]
                redisClient = redis.Redis(connection_pool=pool)
                redisKey = 'quote:' + author
                result = redisClient.srandmember(redisKey)
                if result:
                        response_headers = [
                                ('Content-type','application/json'),
                                ('Content-Length', str(len(result)))
                        ]
                        start_response('200 OK', response_headers)
                        return result
                else:
                        start_response('404 Not Found', [])
                        return ''
        else:
                start_response('404 Not Found', [])
                return ''
