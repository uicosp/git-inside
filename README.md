这是一个用来揭示 git 内部运作原理的演示工具

1. 先执行 `wails build` 构建程序
2. 运行 `wails serve` 启动后端程序
3. 运行 `cd frontend && npm run serve` 启动前端程序
4. 在浏览器打开 `http://localhost:8080/`
5. 点击 `演示` tab
6. 输入本地 git 仓库路径，点击 `go` 便会展示 git 仓库内部数据组织形式
7. 开启程序后会实时监听目录文件变化，并渲染最新的数据到页面上

⚠️注意：由于本程序只是一个演示工具，如果打开真实 git 项目会遇到渲染性能问题，
建议新建空的 git 仓库，通过一步步 `git add` \ `git commit` 来观察 git 仓库内部的变化。

👉 推荐跟着 [GIT 原理讲解](/git_inside.md) 一步步执行 git 命令以窥探其内部的数据变化