import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import UploadPhotoView from '../views/UploadPhotoView.vue'
import ProfileView from '../views/ProfileView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/',
			redirect: '/session'	// Default redirect to /session
		},
		{ 
			path: '/session',
			component: LoginView 	// Login screen
		}, 
    	{ 
			path: '/home',
			component: HomeView 	// Home screen
		}, 
		{	path: '/link1', 
			component: UploadPhotoView	// Screen where you can upload a new photo
		},
		{
			path: '/some/:id/link',
			component: ProfileView	// User profile screen
		},
	]
})

export default router	
