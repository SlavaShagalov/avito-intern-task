#!/bin/bash

# Получение всех баннеров (3 штуки).

ab -n 1000 -c 1000 \
-H "token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMxNjU5ODQsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoxfQ.DxZ9qQW2ydo6gsH22EqEkkKavVntM2XJpsay9Wa1y5M" \
"http://localhost:8000/api/v1/banner"
