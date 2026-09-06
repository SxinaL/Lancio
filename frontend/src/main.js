import { createApp } from 'vue';
import App from './App.vue';
// import router from './router' 
// createApp(App).mount('#app');


// main.js
// import { createApp } from 'vue'
// import App from './App.vue'
import router from './router' // 确保导入了 router

const app = createApp(App)
app.use(router) // 这一步至关重要！必须执行
app.mount('#app')