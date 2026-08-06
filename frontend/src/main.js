import { createApp } from 'vue'
import { createPinia } from 'pinia'
import {
	Button,
	Card,
	Form,
	Input,
	Alert,
	Menu,
	Layout,
	Table,
	Modal,
	Select,
	Space,
	Tag,
	Pagination,
	Checkbox,
	Breadcrumb,
	Tooltip,
	Radio,
	Switch,
	Drawer,
	List
} from '@arco-design/web-vue'
import '@arco-design/web-vue/es/menu/style/css'
import App from './App.vue'
import router from './router'
import './assets/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 按需注册组件
app.use(Button)
app.use(Card)
app.use(Form)
app.use(Input)
app.use(Alert)
app.use(Menu)
app.use(Layout)
app.use(Table)
app.use(Modal)
app.use(Select)
app.use(Space)
app.use(Tag)
app.use(Pagination)
app.use(Checkbox)
app.use(Breadcrumb)
app.use(Tooltip)
app.use(Radio)
app.use(Switch)
app.use(Drawer)
app.use(List)

app.mount('#app')