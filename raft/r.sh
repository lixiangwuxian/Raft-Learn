#!/bin/bash

# 启动程序并传入参数，然后放入后台运行
./test/1233/m 127.0.0.1:1233 ./test/1233/config.yaml &
./test/1234/m 127.0.0.1:1234 ./test/1234/config.yaml &
# ./test/1235/m 127.0.0.1:1235 ./test/1235/config.yaml &
# ./test/1235/m 127.0.0.1:1236 ./test/1236/config.yaml &
# ./test/1235/m 127.0.0.1:1237 ./test/1237/config.yaml &

echo "三个程序已在后台运行"