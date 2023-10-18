# Makefile

# 编译目标
TARGET = raft.exe

# Go源文件列表
SOURCES = candidate.go raft.go conf.go follower.go leader.go log.go networkAdapter.go store.go

# 编译命令
build:
	cd raft
	go build -o $(TARGET) $(SOURCES)

# 清理生成的可执行文件
clean:
	cd raft
	rm -f $(TARGET)
