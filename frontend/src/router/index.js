import { createRouter, createWebHashHistory } from 'vue-router'
const router = createRouter({
    // 使用 hash 模式：Wails 桌面环境下刷新/返回时始终请求根路径，
    // 避免 history 模式深层路径刷新导致 window.go 运行时注入丢失
    history: createWebHashHistory(import.meta.env.BASE_URL),
    routes: [
        // 根路径 —— 由 App.vue 根据 GetWindowMode() 动态跳转
        {
            path: '/',
            name: 'Home',
        },
        // ===== 主窗口路由 =====
        {
            path: '/main',
            name: 'MainWin',
            component: () => import('@/views/MainWin.vue'),
            children: [
                {
                    path: 'bank-manager',
                    name: 'BankManager',
                    component: () => import('@/views/BankManager.vue'),
                },
                {
                    path: 'voc-manager/:id',
                    name: 'VocManager',
                    component: () => import('@/views/VocManager.vue'),
                    props: true,
                },
                {
                    path: 'float-manager',
                    name: 'FloatWinManager',
                    component: () => import('@/views/FloatWinManager.vue'),
                },
                
            ],
        },
        // ===== 浮窗路由 =====
        {
            path: '/float',
            name: 'VocFloatWin',
            component: () => import('@/views/VocFloatWin.vue'),
        },
    ],
})

export default router
