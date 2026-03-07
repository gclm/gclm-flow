# GitHub Actions VM Image CI Notes

用于 `libguestfs` / `virt-customize` / QCOW2 相关 GitHub Actions 镜像构建排障。

## 何时查看

- `virt-customize` 在 CI 中失败
- guest 架构与 runner 架构可能不一致
- VM 镜像构建偶发失败、网络后端异常或 artifact 被覆盖

## 诊断顺序

1. 先确认失败是否发生在 guest 内执行阶段
   - 关注：`virt-customize: error`、`Running: apt-get/dnf`
2. 架构问题优先级最高
   - 关注：`host cpu ... and guest arch ... are not compatible`
   - 处理：按 matrix 架构选择 runner，例如 `amd64 -> ubuntu-22.04`、`arm64 -> ubuntu-22.04-arm`
3. 固定 runner 版本，避免 `ubuntu-latest` 带来不确定性
4. 如果看到 `passt exited with status 1`，在 CI 中优先回退 `slirp`
5. 区分“网络问题”和“包不存在”
   - DNS 问题：`Temporary failure resolving ...`
   - 包不存在：`No match for argument: xxx`
6. matrix 设计要带架构维度，artifact 名称必须带架构，避免覆盖

## 发布补充

- 生成 `RELEASE_NOTES.md` 时，优先稳定的自动生成方案
- release job 需要完整 tags 和提交历史
