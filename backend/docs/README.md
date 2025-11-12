# 📚 项目文档目录

欢迎查阅 GoHotel 酒店管理系统的项目文档！

---

## 📖 文档列表

### 1. [Git 提交规范](./COMMIT_CONVENTION.md)
- Conventional Commits 规范
- 提交类型说明（feat, fix, docs 等）
- 提交信息示例
- 配置 commitlint

### 2. [分支管理策略](./BRANCH_STRATEGY.md)
- Git Flow 工作流
- 分支类型详解（main, develop, feature, release, hotfix）
- 版本标签规范（Semantic Versioning）
- 实战场景示例
- 最佳实践

---

## 🚀 快速开始

### 初次开发

```bash
# 1. 克隆项目
git clone git@github.com:yichen11818/gohotel.git
cd gohotel/backend

# 2. 切换到 develop 分支
git checkout develop

# 3. 创建功能分支
git checkout -b feature/your-feature-name

# 4. 开发并提交
git add .
git commit -m "feat(module): 你的功能描述"

# 5. 推送并创建 Pull Request
git push origin feature/your-feature-name
```

---

## 📝 常用命令速查

### 分支操作
```bash
# 查看所有分支
git branch -a

# 创建功能分支
git checkout -b feature/new-feature

# 切换分支
git checkout develop

# 删除分支
git branch -d feature/old-feature
```

### 提交操作
```bash
# 查看状态
git status

# 提交代码
git add .
git commit -m "feat(booking): 添加预订功能"

# 推送到远程
git push origin feature/booking
```

### 同步更新
```bash
# 更新本地分支
git checkout develop
git pull origin develop

# 将 develop 合并到功能分支
git checkout feature/booking
git merge develop
```

---

## 🎯 工作流程图

```
开发新功能:
  develop → feature/xxx → 开发 → Pull Request → develop

准备发布:
  develop → release/vX.Y.Z → 测试修复 → main + tag → develop

紧急修复:
  main → hotfix/vX.Y.Z → 修复 → main + tag → develop
```

---

## 🔗 相关链接

- **项目仓库**: https://github.com/yichen11818/gohotel
- **Issue 追踪**: https://github.com/yichen11818/gohotel/issues
- **Pull Requests**: https://github.com/yichen11818/gohotel/pulls

---

## ❓ 常见问题

### Q: 我应该从哪个分支创建功能分支？
A: 始终从 `develop` 分支创建。

### Q: 什么时候使用 hotfix 分支？
A: 只有在生产环境（main 分支）出现紧急 bug 时使用。

### Q: 如何命名我的功能分支？
A: 使用 `feature/模块-功能描述`，例如：`feature/user-login`

### Q: 提交信息写什么？
A: 遵循 Conventional Commits 规范，例如：`feat(user): 添加用户注册功能`

---

## 📞 需要帮助？

如有疑问，请：
1. 查阅相关文档
2. 在 Issues 中提问
3. 联系项目维护者

---

**Happy Coding! 🎉**


















