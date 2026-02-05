.PHONY: all build test clean install uninstall help

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

## install: 安装 gclm-engine 到 ~/.gclm-flow
## Note: 配置文件和工作流会在首次运行时自动从二进制中导出
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"
	@echo "  Binary: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo ""
	@echo "Add to PATH:"
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

## run: 快速运行
run: dev
	@echo "Running $(BINARY_NAME)..."
	@$(INSTALL_DIR)/$(BINARY_NAME)

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
