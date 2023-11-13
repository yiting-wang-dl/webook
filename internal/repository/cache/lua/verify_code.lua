local key = KEYS[1]
local cntKey = key..":cnt"
-- user input verification code
local expectedCode =ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

-- not strictly distinguish between code not sent yet or varified too many times
-- exhault the verification time
if cnt == nil or cnt <= 0 then
    return -1
end

-- code is verified
if code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    -- code != expectedCode
    return -2
end