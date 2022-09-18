<template>
  <div class="about">
    <p>Git 是当下每一个程序员的必备技能，网络上关于 git 的使用教程有很多，只要仔细阅读就能掌握常用命令以应付工作。</p>
    <p>
      除了会用，git 内部是怎么运行的？想必大家也有过这样的好奇，但网络上关于内部原理的介绍大多让读者看了之后似懂非懂，
      没有一篇能完全透彻地将其讲明白。
    </p>
    <p>直到我看到了这篇文章
      <a href="https://www.jiqizhixin.com/articles/2019-12-20">这才是真正的Git——Git内部原理揭秘！</a>。
      这次，我想我是真的懂了。
    </p>
    <p>但如何证明作者说的是对的呢？他文章中的动图是手动画的，并不一定能真实反应 git 操作过程中文件的实际变化情况。</p>
    <p> 💡 这时我想到了一个 idea：为什么不通过读取目录树并展示出来的方法来验证呢！ </p>
    <p>说干就干</p>
    <p>凭借自己对 web 前端技术的熟悉，呈现方式就选用 html+css+js，前端框架选择 vue，简单够用。</p>
    <p>现在还剩下一个问题，出于安全原因浏览器中运行的 js 是不能直接访问本地文件的，像上传文件必须由用户事件来触发。</p>
    <p>此时，我想到了之前用过的一个 golang 桌面应用框架 <a href="https://wails.app/about/">wails</a>。</p>
    <p>这是一个神奇的框架，官网的介绍是 "Build desktop apps using Go & Web Technologies"。前端使用 webview，通过
    html/js/css 编写界面；后端使用 go 语言编写业务逻辑。</p>
    <p>万事俱备，只差时间</p>
    <p>一天过去了...</p>
    <p><router-link to="/git">完工验货</router-link></p>

  </div>
</template>

<style scoped>
.about {
  text-align: left;
}
</style>