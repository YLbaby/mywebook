-- 你的验证码在redis上的key
-- phone_code:login:152xxxxxxx
local key = KEYS[1]
-- 验证次数，一个验证码最多重复三次

local cntKey = key..":cnt"
-- 验证码 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
