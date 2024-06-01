import {createRouter, createWebHashHistory} from 'vue-router'
import LoginView from '../views/LoginView.vue'
import HomeView from '../views/HomeView.vue'
import UploadPhotoView from '../views/UploadPhotoView.vue'
import ProfileView from '../views/ProfileView.vue'
import NotFoundView from '../views/NotFoundView.vue'


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
			redirect: '/users/:username/stream'
		},
    	{ 
			path: '/users/:username/stream',
			name: 'home',
			component: HomeView 	// Home screen
		}, 
		{	path: '/users/:username/photos',
			name: 'uploadPhoto',
			component: UploadPhotoView	// Screen where you can upload a new photo
		},
		{
			path: '/users/:username',
			name: 'profile',
			component: ProfileView,	// User profile screen
			props: true
		},
		{
			path: '/not-found',
			name: 'notFoundPage',
			component: NotFoundView, // Not Found Page
		},
		{
			path: '/:catchAll(.*)',
			name: 'notFoundPage',
			component: NotFoundView, // Not Found Page
		},
	]
})

export default router	
