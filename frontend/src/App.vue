<template>
  <RouterView />
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GetWindowMode } from '../wailsjs/go/handler/WindowHandler'

const router = useRouter()

onMounted(async () => {
  try {
    const mode = await GetWindowMode()
    if (mode === 'float') {
      router.push('/float')
    } else {
      router.push('/main/bank-manager')
    }
  } catch (err) {
    // 非 Wails 环境（如纯浏览器调试）默认走主窗口
    router.push('/main/bank-manager')
  }
})
</script>

<style>
/* ===== 全局样式 ===== */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}
/* 去掉所有 a 标签的下划线和默认颜色 */
a {
  text-decoration: none;
  color: inherit;
}

a:visited {
  color: inherit;
}
html,
body {
    width: 100%;
    height: 100%;
    overflow: hidden;
    background: transparent;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
        "PingFang SC", "Microsoft YaHei", sans-serif;
    user-select: none;
    -webkit-user-select: none;
}

#app {
    width: 100%;
    height: 100%;
}
</style>
