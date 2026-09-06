<template>
    <div class="main-window">
        <!-- ===== 顶部标题栏 ===== -->
        <ToolBar />

        <!-- ===== 主体区域 ===== -->
        <div class="main-body">
            <!-- 左侧工具栏 -->
            <SiderBar />

            <!-- 右侧内容区 -->
            <main class="content-area">
                 <RouterView />
            </main>
        </div>
    </div>
</template>

<script setup>
import ToolBar from '@/components/mainwin/ToolBar.vue'
import SiderBar from '@/components/mainwin/SiderBar.vue'
</script>

<style scoped>
/* ===== 主窗口容器 ===== */
.main-window {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #1a1a2e;
    color: rgba(255, 255, 255, 0.85);
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
        "PingFang SC", "Microsoft YaHei", sans-serif;
    overflow: hidden;

}

/* ===== 主体 ===== */
.main-body {
    height: calc(100% - 40px);
    width: 100%;
    display: flex;
    /* overflow: hidden; */
    /* border: 1px solid rgb(255, 0, 0); */
}

/* ===== 右侧内容区 ===== */
.content-area {
    width: calc(100% - 180px);
    height: 100%;
    display: flex;

    flex-direction: row;
    /* align-items: stretch; */
    justify-content: flex-start;
    background: rgba(0, 0, 0, 0.15);
    overflow: auto;
}
</style>