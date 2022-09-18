# 1 `.git` 目录结构

```
.
├── COMMIT_EDITMSG
├── FETCH_HEAD
├── HEAD // 当前指针
├── config
├── description
├── hooks
├── index // 暂存区
├── info
│   └── exclude
├── logs
│   ├── HEAD
│   └── refs
│       └── heads
│           └── master
├── objects // 对象存储仓库
│   ├── 4f
│   │   └── 2ee5ef94b02b742696aca3b202b30e14d53c28
│   ├── aa
│   │   └── a96ced2d9a1c8e72c56b253a0e2fe78393feb7
│   ├── ce
│   │   └── 013625030ba8dba906f756967f9e9ca394464a
│   ├── info
│   └── pack
└── refs
    ├── heads
    │   └── master // 主分支
    └── tags

14 directories, 26 files

```

`.git` 的目录结构如上图所示。

# 2 文件存储格式

`.git` 目录中除了 `objects` 和 `index` 是按二进制格式存储外，其他所有文件都是以纯文本格式存储的。

之所以 `objects` 和 `index` 没有直接用文本方式存储，是出于性能的考虑。

git 提供了对应的命令来读取其中的内容：

```
// 读取 obejcts 内容
git cat-file

// 读取 index 文件
git ls-file
```

# 3 objects 对象存储
## 3.1 对象的分类

所有对象存储在 objects 目录下，共分为三种类型：
- commit
- tree
- blob

可以用 `git cat-file -t` 看一下文件类型：

```bash
git cat-file -t 4f2e
// commit

git cat-file -t aaa9
// tree

git cat-file -t ce01
// blob
```

## 3.2 对象中存储的内容

在用 `git cat-file -p` 命令看下具体内容：

### 3.2.1 commit

```bash
git cat-file -p 4f2e
```

```
tree aaa96ced2d9a1c8e72c56b253a0e2fe78393feb7
author xiayuxiao <xiayuxiao@weidian.com> 1657212968 +0800
committer xiayuxiao <xiayuxiao@weidian.com> 1657212968 +0800

first commit
```

### 3.2.2 tree

```bash
git cat-file -p aaa9
```

```
100644 blob ce013625030ba8dba906f756967f9e9ca394464a hello.txt
```

### 3.2.3 blob

```bash
git cat-file -p ce01
```

```
hello
```


## 3.3 git cat-file 命令原理

```bash
hexdump -v -e '"\\\x" 1/1 "%02x"' objects/ce/013625030ba8dba906f756967f9e9ca394464a
```

```bash
\x78\x01\x4b\xca\xc9\x4f\x52\x30\x63\xc8\x48\xcd\xc9\xc9\xe7\x02\x00\x1d\xc5\x04\x14
```

```ruby
❯ irb          

irb(main):001:0> require 'zlib'
=> true
irb(main):002:0> Zlib::Inflate.inflate("\x78\x01\x4b\xca\xc9\x4f\x52\x30\x63\xc8\x48\xcd\xc9\xc9\xe7\x02\x00\x1d\xc5\x04\x14")
=> "blob 6\x00hello\n"
irb(main):003:0> 

```

```
type = blob|tree|commit 三选一
header = "#{type} #{content.length}\0"
store = header + content

// hash长度为40位，前2位目录名，后38位文件名
sha1 = Digest::SHA1.hexdigest(store)

// 为节省存储空间，在存入磁盘前进行压缩
zlib_content = Zlib::Deflate.deflate(store)
```

# 4 index 文件
[Git: Understanding the Index File](https://mincong.io/2018/04/28/git-index/)

# 5 git 命令引发的内部变化

这一节我们来看看 git 命令执行后会引发 `.git` 目录下发生什么变化。

## 5.1 git add 
1. 将对象存储到 ./git/objects/ 下，类型为 blob
2. 同时更新 index 文件（也叫暂存区），不存在的话会创建

## 5.2 git status
1. 对比 index 和 tree（Changes to be committed）
2. 对比 index 和 当前工作目录（Untracked files）

## 5.3 git commit
1. 根据 index 文件生成 tree 对象并存储到 .git/objects/ 下
2. 生成 commit 对象并存储到  .git/objects/ 下
3. 将当前分支指向生成的 commit

## 5.4 git log 是如何获取历史记录的？
commit 通过 parent 指针组成了一条单链表，`git log` 就是通过遍历链表来获取历史记录的。

## 5.5 git branch 时发生了什么？
1. 在 .git/refs/heads/ 下创建新的分支文件
2. 将新创建的分支指向 HEAD 所指向的 commit

## 5.6 git checkout 时发生了什么？
1. 将 HEAD 指向签出的分支文件
2. 更新 working directory

### 5.6.1 为什么有时候 checkout 会卡顿？
因为当签出分支与当前分支内容差异较大时，涉及大量的 IO 操作，可能会有卡顿感。

### 5.6.2 与 git checkout commit 的区别？
签出某个具体 commit 时，HEAD 跳过分支直接指向了 commit，这时会提示

>You are in 'detached HEAD' state

意思就是 HEAD 指针游离于分支之外了。

## 5.7 git merge
1. 合并产生的 commit 文件会有两个 parent 指针
2. 产生冲突的文件会残留下来，git gc 会定期进行垃圾回收，我们可以通过 git prune 命令手动清除 https://git-scm.com/docs/git-gc

## 5.8 git reset
1. 改变 HEAD 所指向的 commit (--soft)
2. 将 index 区域更新为 HEAD 所指向的 commit 对应的 tree 里的内容 (--mixed)
3. 将 working directory 更新为 tree 包含的文件 (--hard)
https://git-scm.com/docs/git-reset

## 5.9 git revert
撤销某一次 commit 的操作，即对某次 commit 进行取反，比如：
- 新增文件 -> 删除文件
- 删除文件 -> 新增文件
- 文件中增加一行 -> 文件中删除一行
- 文件中删除一行 -> 文件中增加一行
- ...
然后生成一个新的 commit

## 5.10 git rebase
https://www.atlassian.com/git/tutorials/merging-vs-rebasing

# 6 数据一致性&安全

## 6.1 默克尔树 - Merkel Tree
git 使用默克尔树来保证分布式数据的一致性
https://blog.csdn.net/itworld123/article/details/115042033

## 6.2 如何在 git 中删除敏感数据
错误做法：`git rm`
正确做法：
https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository

# 7 git 底层命令
https://git-scm.com/book/zh/v2/Git-内部原理-底层命令与上层命令
