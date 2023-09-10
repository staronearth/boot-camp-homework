--验证码在redis上的key
local key=KEYS[1]
-- 使用次数，也就是验证次数
local cntKey=key..":cnt"
-- 你的验证码
local val = ARGV[1]
-- 验证码的有效时间是十分钟，600秒
local ttl=tonumber(redis.call("ttl",key))

-- -1 是key存在，但是没有过期时间
if ttl==-1 then
    --key存在没有过期时间
    return -2
    -- -2代表key不存在，ttl<540是发了一个验证码。已经超过一分钟
elseif ttl==-2 or ttl<540 then
    redis.call("set",key,val)
    redis.call("expire",key,600)
    redis.call("set",cntKey,3)
    redis.call("expire",cntKey,600)
    return 0
else
    -- 发送太频繁
    return -1
end