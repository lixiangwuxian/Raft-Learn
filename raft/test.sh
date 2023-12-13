#!/bin/bash

# 编译 Go 程序
go build -o m

# 检查编译是否成功
if [ ! -f "./m" ]; then
    echo "编译失败，未找到可执行文件"
    exit 1
fi

# 创建目标文件夹（如果它们不存在）并复制文件
for dir in test/1233 test/1234 test/1235 test/1236 test/1237; do
    mkdir -p "$dir"
    cp ./m "$dir/"
done

echo "可执行文件已成功复制到目标文件夹"
