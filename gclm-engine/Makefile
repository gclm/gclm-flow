.PHONY: all build test clean install uninstall start stop restart status logs run help

# 变量定义
BINARY_NAME=gclm-engine
BUILD_DIR=build
WORKFLOWS_DIR=workflows
HOME_DIR=$(shell echo ~)
INSTALL_DIR=$(HOME_DIR)/.gclm-flow

# Go 相关
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFLAGS=-v

# 默认目标
all: build

## build: 编译 gclm-engine 二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: 运行所有测试
test:
	@echo "Running tests..."
	$(GOTEST) ./... -v -count=1

## clean: 清理构建文件
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## install: 安装 gclm-engine 到 ~/.gclm-flow (配置和工作流首次运行时自动导出)
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo ""
	@echo "Binary installed to: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "Quick start:"
	@echo "  make start           # Start daemon in background"
	@echo "  make stop            # Stop daemon"
	@echo "  make status          # Check status"
	@echo "  make logs            # View logs"
	@echo ""
	@echo "Or use directly:"
	@echo "  $(BINARY_NAME) serve         # Start in foreground"
	@echo "  $(BINARY_NAME) serve -d      # Start in background (daemon)"
	@echo "  $(BINARY_NAME) serve --stop  # Stop background server"
	@echo ""
	@echo "Add to PATH (optional, for convenience):"
	@echo "  export PATH=\"$(INSTALL_DIR):\$$PATH\""
	@echo ""
	@echo "First run will auto-initialize config and workflows."

## uninstall: 卸载 gclm-engine
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstall complete"
	@echo "Note: Config, workflows, and database at $(INSTALL_DIR) were preserved"
	@echo "      To remove everything: rm -rf $(INSTALL_DIR)"

## reset: 重置配置（删除配置文件和工作流，保留数据库）
reset:
	@echo "Resetting configuration..."
	@rm -f $(INSTALL_DIR)/gclm_engine_config.yaml
	@rm -rf $(INSTALL_DIR)/workflows
	@echo "Configuration reset."
	@echo "Next run will re-initialize from embedded defaults."

## purge: 完全清理（包括数据库）
purge:
	@echo "Purging all gclm-flow data..."
	@rm -rf $(INSTALL_DIR)
	@echo "Purge complete."

## run: 快速运行（前台模式）
run: build
	@echo "Running $(BINARY_NAME) in foreground..."
	@$(BUILD_DIR)/$(BINARY_NAME) serve

## start: 启动服务（后台守护进程模式）
start:
	@echo "Starting $(BINARY_NAME) in daemon mode..."
	@if [ -f "$(INSTALL_DIR)/$(BINARY_NAME)" ]; then \
		$(INSTALL_DIR)/$(BINARY_NAME) serve -d; \
	else \
		echo "Error: $(BINARY_NAME) not installed. Run 'make install' first."; \
		exit 1; \
	fi

## stop: 停止后台服务
stop:
	@echo "Stopping $(BINARY_NAME)..."
	@if [ -f "$(INSTALL_DIR)/$(BINARY_NAME)" ]; then \
		$(INSTALL_DIR)/$(BINARY_NAME) serve --stop; \
	else \
		echo "Error: $(BINARY_NAME) not installed."; \
		exit 1; \
	fi

## restart: 重启服务
restart: install stop start

## status: 查看服务状态
status:
	@if [ -f "$(INSTALL_DIR)/gclm-engine.pid" ]; then \
		pid=$$(cat $(INSTALL_DIR)/gclm-engine.pid); \
		if ps -p $$pid > /dev/null 2>&1; then \
			echo "$(BINARY_NAME) is running (PID: $$pid)"; \
			echo "URL: http://localhost:9988"; \
		else \
			echo "$(BINARY_NAME) is not running (stale PID file)"; \
			rm -f $(INSTALL_DIR)/gclm-engine.pid; \
		fi \
	else \
		echo "$(BINARY_NAME) is not running"; \
	fi

## logs: 查看服务日志
logs:
	@if [ -f "$(INSTALL_DIR)/gclm-engine.log" ]; then \
		tail -f $(INSTALL_DIR)/gclm-engine.log; \
	else \
		echo "No log file found at $(INSTALL_DIR)/gclm-engine.log"; \
	fi

## package: 打包为发布文件（用于 GitHub Releases）
package: build
	@echo "Creating release packages..."
	@mkdir -p dist
	@cp $(BUILD_DIR)/$(BINARY_NAME) dist/$(BINARY_NAME)-$(GOOS)-$(GOARCH)
	@cd $(WORKFLOWS_DIR) && tar -czf ../../dist/workflows.tar.gz *.yaml
	@echo "Package complete: dist/"

## help: 显示帮助信息
help:
	@echo "gclm-engine Makefile"
	@echo ""
	@echo "使用方法: make [target]"
	@echo ""
	@echo "可用目标:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
