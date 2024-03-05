-- 你的验证码在redis上的key
-- phone_code:login:152xxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码最多重复三次
-- 拼接字符串
local cntKey = key..":cnt"
-- 验证码 123456
local val = ARGV[1]
-- 过期时间
-- 当 key 不存在时，返回 -2。 当 key 存在但没有设置剩余生存时间时，返回 -1。 否则，以秒为单位，返回 key 的剩余生存时间
local ttl = tonumber(redis.call("ttl", key))

if ttl == -1 then
    return -2
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 0
else
    -- 发送太频繁
    return -1
end