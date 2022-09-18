<template>
  <div class="file">
    <div class="file-head" :style="headStyle">
      <div class="file-name">
        {{ name }}
      </div>
    </div>
    <pre class="file-body" :style="bodyStyle"><slot></slot></pre>
  </div>
</template>

<script>
export default {
  name: "File",
  props: {
    name: String,
    headStyle: String,
    bodyStyle: String
  }
}
</script>

<style scoped>
.file {
  width: auto;
  border: 1px solid #f6f6f6;
  margin: 10px;
  display: inline-block;
  vertical-align: top;
}

.file-head {
  background: #f6f6f6;
}

.file-name {
  padding: 2px 10px;
  font-weight: 500;
}

.file-body {
  margin: 0;
  padding: 4px 10px;
}
</style>