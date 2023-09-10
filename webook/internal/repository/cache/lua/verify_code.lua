local key=KEYS[1]
-- 用户输入的code
local expectedCode=ARGV[1]

local cntKey=key..":cnt"
-- 转成一个数字
local cnt=tonumber(redis.call("get",cntKey))
local code = redis.call("get",key)
if cnt<=0 then
    -- 说明用户一直输错
    return -1
elseif expectedCode==code then
    -- 验证码正确
    redis.call("del",key)
    redis.call("del",cntKey)
    return 0
else
    -- 用户手一抖输错了
    -- 可验证次数-1
    redis.call("decr",cntKey)
    return -2
end