-- 1, 2, 3, 4, 5, 6, 7 elements
-- ZREMRANGEBYSCORE key1 0 6
-- 7 after execution

-- rate limit object
local key = KEYS[1]
-- window size
local window = tonumber(ARGV[1])
-- threshold
local threshold = tonumber( ARGV[2])
local now = tonumber(ARGV[3])
-- window starting time
local min = now - window

redis.call('ZREMRANGEBYSCORE', key, '-inf', min)
local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
-- local cnt = redis.call('ZCOUNT', key, min, '+inf')
if cnt >= threshold then
    -- implement rate limit
    return "true"
else
    -- set score and member to now
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end