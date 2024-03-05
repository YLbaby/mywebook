local key = KEYS[1]

local expectedCode = ARGV[1]
local code = redis.call("get", key)
local cntKey = key..":cnt"

local cnt = tonumber(redis.call("get", cntKey))
if cnt <= 0 then
    -- 说明有人一直输错
    return -1
elseif expectedCode == code then
    -- 输入对了
    -- 用完，不能再用了
    redis.call("set", cntKey, -1)
    return 0
else
    -- 用户输入错误
    redis.call("decr", cntKey)
    return -2
end