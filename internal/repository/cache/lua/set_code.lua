-- code:biz:phone
local key = KEYS[1]
local cntKey = key..":cnt"
-- verification code we prepared
local val = ARGV[1]
-- code is valid for 10 min, 600s
local ttl = tonumber(redis.call("ttl", key))
-- -1 = key exists, but no expiration time
if ttl == -1 then
    return -2
-- -2 = key doesn't exist, ttl < 540 = sent a code more than 60s ago
elseif ttl == -2 or ttl < 540 then
    -- good to send verification code
    redis.call("set", key, val)
    -- 600 seconds
    redis.call("expire", key, 600)
    -- maximum verification time 3
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- send too frequent. a code is sent less then 1 min ago
    return -1
end


