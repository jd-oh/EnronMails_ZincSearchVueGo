import './assets/main.css'
import './index.css'
import 'simple-datatables/dist/style.css';

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia';

import { DataTable } from 'simple-datatables';


const app = createApp(App)
const pinia = createPinia();

app.use(pinia);

app.use(router)

app.mount('#app')
