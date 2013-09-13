<?php

// php -S 0.0.0.0:80 router.php
$matches = array();

if (preg_match('/\/quote\/([^\/]+)/', $_SERVER["REQUEST_URI"], $matches)) {
    $author = $matches[1];
    $redisUrl = parse_url(getenv('REDISTOGO_URL'));

    $redis = new Redis();
    $redis->pconnect($redisUrl['host'], $redisUrl['port']);
    $redis->auth($redisUrl['pass']);

    $redisKey = "quote:" . $author;
    $result = $redis->sRandMember($redisKey);

    if($result) {
        header("Content-type: application/json");
        header("Content-Length: " . strlen($result));
        echo $result;
    } else {
        return false;
    }
} else {
    return false;
}
?>
