wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"
-- 这个要改为你的注册的数据
wrk.body='{"email":"123456@qq.com", "password": "123456@qq"}'