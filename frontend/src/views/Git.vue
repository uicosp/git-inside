<template>
  <div class="git">
    <div class="git-repo">
      <input placeholder="请输入 git 仓库路径" v-model="repoPathNew" style="width: 400px"/>
      <button style="margin-left: 10px" v-on:click="gotoRepo">GO</button>
    </div>
    <div class="git-content" v-if="repoPath!==''">
      <div>
        <h2>./ 当前工作目录</h2>
        <file class="file" v-for="f in files" :key="f.Name" :name="relativeWdPath(f.Path)"
              head-style="background: rgb(158 158 158 / 20%)" body-style="background: rgb(158 158 158 / 3%)">
          {{ f.Content }}
        </file>
      </div>
      <div>
        <h2>.git/ 暂存区 & HEAD</h2>
        <file :name="index.Name" v-if="index"
              head-style="background: #8bc34a8a" body-style="background: #8bc34a4f">
          {{ index.Content }}
        </file>
        <file :name="HEAD.Name" id="HEAD" v-if="HEAD"
              head-style="background: rgb(255 235 59 / 60%)" body-style="background: rgb(255 235 59 / 30%)">
          {{ HEAD.Content }}
        </file>
      </div>
      <div>
        <h2>.git/refs/heads/ 分支</h2>
        <file v-for="head in heads" :key="head.Name" :name="head.Name"
              :id="relativePath(head.Path)"
              head-style="background: rgb(255 235 59 / 60%)" body-style="background: rgb(255 235 59 / 30%)">
          {{ head.Content }}
        </file>
      </div>
      <div>
        <h2>.git/objects/ 对象仓库</h2>
        <div>
          <h3>commit</h3>
          <file v-for="obj in commits" :key="obj.Hash" :name="obj.Hash" :id="obj.Hash"
                head-style="background: rgb(255 193 7 / 75%)" body-style="background: rgb(255 193 7 / 30%)">
          <div v-html="renderCommit(obj)"></div>
          </file>
        </div>
        <div title="tree">
          <h3>tree</h3>
          <file v-for="obj in trees" :key="obj.Hash" :name="obj.Hash" :id="obj.Hash"
                head-style="background: rgb(139 195 74 / 84%)" body-style="background: #8bc34a63">
              <div v-for="row, index in obj.Content.split('\n')" :key="index">
                {{ void (it = row.split(/\s+/)) }}
                {{ it[0] }} {{ it[1] }} <span :id="obj.Hash+':'+it[2]">{{ it[2] }}</span> {{ it[3] }}
              </div>
          </file>
        </div>
        <div title="blob">
          <h3>blob</h3>
          <file v-for="obj in blobs" :key="obj.Hash" :name="obj.Hash" :id="obj.Hash"
                head-style="background: rgb(33 150 243 / 50%)" body-style="background: rgb(33 150 243 / 19%)">
            {{ obj.Content }}
          </file>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {jsPlumb} from 'jsplumb'
import Wails from "@wailsapp/runtime";
import File from "@/components/File";

// jsplumb 概念介绍
// https://docs.jsplumbtoolkit.com/community/lib/basic-concepts

let plumbIns = jsPlumb.getInstance()
let defaultConfig = {
  // 对应上述基本概念
  anchor: ['Left', 'Right', 'Top', 'Bottom', [0.3, 0, 0, -1], [0.7, 0, 0, -1], [0.3, 1, 0, 1], [0.7, 1, 0, 1]],
  connector: ['Bezier'],
  endpoint: 'Blank',
  // 添加样式
  paintStyle: {stroke: '#909399', strokeWidth: 2}, // connector
  // endpointStyle: { fill: 'lightgray', outlineStroke: 'darkgray', outlineWidth: 2 } // endpoint
  // 添加 overlay，如箭头
  overlays: [['PlainArrow', {width: 8, length: 8, location: 1}]] // overlay
}
let commit2parentArrow = {
  ...defaultConfig,
  paintStyle: {stroke: 'orange', strokeWidth: 1}, // connector
  connector: ['Straight'],
}
let commit2treeArrow = {
  ...defaultConfig,
  paintStyle: {stroke: 'green', strokeWidth: 1}, // connector
}
let red = {
  ...defaultConfig,
  paintStyle: {stroke: 'red', strokeWidth: 1}, // connector
}
let tree2blobArrow = {
  ...defaultConfig,
  paintStyle: {stroke: 'red', strokeWidth: 1}, // connector
}

export default {
  name: "Git",
  components: {
    File
  },
  data: function () {
    return {
      repoPath: "",
      repoPathNew: "",
      files: [],
      index: null,
      HEAD: null,
      heads: [],
      objects: []
    }
  },
  computed: {
    commits: function () {
      return this.objects.filter(o => o.Type === 'commit')
    },
    trees: function () {
      return this.objects.filter(o => o.Type === 'tree')
    },
    blobs: function () {
      return this.objects.filter(o => o.Type === 'blob')
    }
  },
  methods: {
    renderCommit(obj) {
      const commitHash = obj.Hash
      const rows = obj.Content.split('\n')

      let str=''
      let ok = false
      for (let i = 0; i < rows.length; i++) {
        const row = rows[i]
        if(ok) {
          str += `${row}\n`
          continue
        }
        if (row.startsWith('tree')) {
          const treeHash = row.split(" ")[1]
          str += `tree <span id="${commitHash}:${treeHash}">${treeHash}</span>\n`
        } else if (row.startsWith('parent')) {
          const parentHash = row.split(" ")[1]
          str += `parent <span id="${commitHash}:${parentHash}">${parentHash}</span>\n`
        } else {
          ok = true
          str += `${row}\n`
        }
      }
      return str
    },
    relativeWdPath(path) {
      return path.replace(this.repoPath + "/", "")
    },
    relativePath(path) {
      return path.replace(this.repoPath + "/.git/", "")
    },
    async gotoRepo() {
      this.repoPath = this.repoPathNew
      await window.backend.Git.SetRepoPath(this.repoPath)
      this.fetchData()
    },
    fetchData() {
      const p1 = window.backend.Git.Files().then((files) => {
        if (files == null) {
          files = []
        }
        this.files = files
      })
      const p2 = window.backend.Git.Index().then((index) => {
        this.index = index
      })
      const p3 = window.backend.Git.HEAD().then((HEAD) => {
        this.HEAD = HEAD
      })
      const p4 = window.backend.Git.Heads().then((Heads) => {
        if (Heads == null) {
          Heads = []
        }
        this.heads = Heads
      })
      const p5 = window.backend.Git.Objects().then((objects) => {
        if (objects == null) {
          objects = []
        }
        this.objects = objects
      })
      Promise.all([p1, p2, p3, p4, p5]).then(() => {
        this.$nextTick(() => {
          // 连线
          let relations = []
          // HEAD
          if (this.heads.length > 0) {
            const ref = this.HEAD.Content.trim().replace("ref: ", "")
            relations.push(['HEAD', ref, red])
          }
          // heads
          for (let head of this.heads) {
            const commit = head.Content.trim()
            relations.push([this.relativePath(head.Path), commit, red])
          }
          // commits
          for (let commit of this.commits) {
            const trees = commit.Content.split('\n')
                .filter(line => line.startsWith("tree "))
                .map(line => line.substring(5))
            for (let tree of trees) {
              relations.push([commit.Hash+":"+tree, tree, commit2treeArrow])
            }

            const parents = commit.Content.split('\n')
                .filter(line => line.startsWith("parent "))
                .map(line => line.substring(7))
            for (let parent of parents) {
              relations.push([commit.Hash+":"+parent, parent, commit2parentArrow])
            }
          }
          // trees
          for (let tree of this.trees) {
            const blobs = tree.Content.split('\n')
              .filter(line => line.substr(7, 4) === 'blob')
              .map(line => line.substr(12, 40))
            for (let blob of blobs) {
              var rand = blob.substring(0,20).split("").map(it => it.charCodeAt()).reduce((a, b) => a+b, 0) / (20*128)
              var randomColor = Math.floor(rand * 16777215).toString(16);
              relations.push([tree.Hash + ":" + blob, blob, {
                ...tree2blobArrow,
                paintStyle: { stroke: '#' + randomColor, strokeWidth: 1 }, // connector
              }])
            }
            const subTrees = tree.Content.split('\n')
              .filter(line => line.substr(7, 4) === 'tree')
              .map(line => line.substr(12, 40))
            for (let subTree of subTrees) {
              relations.push([tree.Hash, subTree, commit2parentArrow])
            }
          }
          // 绘制连线
          plumbIns.deleteEveryConnection()
          plumbIns.deleteEveryEndpoint()
          plumbIns.ready(function () {
            for (let item of relations) {
              plumbIns.connect({
                source: item[0],
                target: item[1]
              }, item[2] ? item[2] : defaultConfig)
            }
          })
        })
      })
    }
  },
  async mounted() {
    Wails.Events.On("file_changed", () => {
      console.log("file_changed")
      this.fetchData()
    })
  },
  beforeRouteEnter(to, from, next) {
    const els = document.querySelectorAll("[class^='jtk']")
    els.forEach((el) => el.style.display = '')
    next()
  },
  beforeRouteLeave(to, from, next) {
    // 修复路由切换时 jsplumb 连线出现在其他页面的问题
    const els = document.querySelectorAll("[class^='jtk']")
    els.forEach((el) => el.style.display = 'none')
    next()
  }
}
</script>

<style scoped>
.git {
  text-align: left;
  margin: 10px;
}
</style>
