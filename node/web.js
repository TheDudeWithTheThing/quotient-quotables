// node web.js
var express = require('express'),
    app = express(),
    server = require('http').createServer(app),
    redisUrl = require('url').parse(process.env.REDISTOGO_URL),
    redis = require('redis').createClient(redisUrl.port, redisUrl.hostname);

redis.auth(redisUrl.auth.split(':')[1]);
app.use(express.errorHandler());

app.get('/quote/:author', function(req, res) {
  redis.srandmember('quote:' + req.params.author, function(err, result) {
    if (!result || err) {
      res.send(404);
    }
    else {
      res.type('json');
      res.send(result);
    }
  });
});

var port = 80;
server.listen(port, function() {
  console.log('Listening on ' + port);
});
