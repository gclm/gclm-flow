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

## install: 安装 gclm-engine 和 workflows 到 ~/.gclm-flow
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installing workflows to $(INSTALL_DIR)/workflows..."
	@mkdir -p $(INSTALL_DIR)/workflows
	@cp $(WORKFLOWS_DIR)/*.yaml $(INSTALL_DIR)/workflows/
	@echo "Installation complete!"
	@echo "  Binary: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "  Workflows: $(INSTALL_DIR)/workflows/"
	@echo ""
	@echo "Add to PATH:"
	@echo "  export PATH=\"$(INSTALL_DIR):\$$PATH\""

## uninstall: 卸载 gclm-engine
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@rm -rf $(INSTALL_DIR)/workflows
	@echo "Uninstall complete"
	@echo "Note: Database at $(INSTALL_DIR)/gclm-engine.db was preserved"

## dev: 开发模式 - 快速构建到 ~/.gclm-flow
dev:
	@echo "Building for development..."
	@mkdir -p $(INSTALL_DIR)
	$(GOBUILD) -o $(INSTALL_DIR)/$(BINARY_NAME) .
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@mkdir -p $(INSTALL_DIR)/workflows
	@cp $(WORKFLOWS_DIR)/*.yaml $(INSTALL_DIR)/workflows/
	@echo "Dev build complete!"
	@echo "  Binary: $(INSTALL_DIR)/$(BINARY_NAME)"
	@echo "  Workflows: $(INSTALL_DIR)/workflows/"
	@echo ""
	@echo "Run: $(INSTALL_DIR)/$(BINARY_NAME) workflow list"

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
