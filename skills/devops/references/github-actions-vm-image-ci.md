# GitHub Actions + virt-customize 排障模板（QCOW2/云镜像）

适用场景：`libguestfs/virt-customize` 在 CI 中构建镜像失败。

## 快速诊断顺序

1. 先确认失败是否在 guest 内执行阶段  
- 关键字：`virt-customize: error`、`Running: apt-get/dnf`、`host cpu ... guest arch ... not compatible`。

2. 架构问题优先级最高  
- 现象：`host cpu (x86_64) and guest arch (aarch64) are not compatible`。  
- 修复：按 matrix 架构选择 runner。  
  `amd64 -> ubuntu-22.04`，`arm64 -> ubuntu-22.04-arm`。

3. runner 版本稳定性  
- 对 `libguestfs` 任务固定 runner，不用 `ubuntu-latest`。  
- 推荐：`ubuntu-22.04`（避免镜像更新引入不确定性）。

4. libguestfs 网络后端问题  
- 现象：`passt exited with status 1`。  
- 修复：CI 中移除 `passt`，回退 `slirp`（仅 CI）。

5. 包管理失败要区分“网络问题”与“包不存在”  
- DNS 问题：`Temporary failure resolving ...`。  
- 包不存在：`No match for argument: xxx`。  
- 修复策略：  
  - 网络问题：先用官方源安装，后切国内源；切源后不做 `apt update/dnf clean`。  
  - 包不存在：按发行版维护独立包清单（如 Rocky 10 移除 `htop`）。

6. matrix 设计  
- 默认用双维：`image x arch`。  
- artifact 名称必须带架构，避免覆盖：`${image}-${arch}-qcow2`。

## Release 日志建议

- 使用 `git-cliff` 生成 `RELEASE_NOTES.md`（替代手写 `git log`）。
- `release` job 需要 `fetch-depth: 0` 与 `fetch-tags: true`。
- 将 `body_path: RELEASE_NOTES.md` 传给 release action。

## 参考 workflow 片段

```yaml
jobs:
  build:
    runs-on: ${{ matrix.arch == 'arm64' && 'ubuntu-22.04-arm' || 'ubuntu-22.04' }}
    strategy:
      matrix:
        image: [debian12, ubuntu2404]
        arch: [amd64, arm64]
```
