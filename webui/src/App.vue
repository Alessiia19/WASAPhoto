<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
			isLogged: false,
		}
	},

	methods:{

		doLogout(loginStatus){
			this.isLogged = loginStatus
			this.$router.replace("/session")
		},

		changeLoginStatus(loginStatus){
			this.isLogged = loginStatus
		},


	},

	// Perform initial setup for the application
	created() {
		if (!localStorage.getItem('appInitialized')) {
			localStorage.clear();
			localStorage.setItem('appInitialized', true);
		}
	},

	// Perform actions after the component has been mounted
	mounted() {
		if (!localStorage.getItem('userID')) {
			this.$router.replace("/session");
		} else {
			this.isLogged = true;
		}
	},
}
</script>

<template>

	<div class="container-fluid">
		<div class="row">
			<div class="col p-0">
				<main>
					<Sidebar v-if="isLogged"
					@doLogoutMethod="doLogout"/>
					
					<RouterView
					@changeLoginStatusmethod="changeLoginStatus"/>
				</main>
			</div>
		</div>
	</div>


	<!--
	<div id="app">
		 La rotta principale verrà gestita dal router 
		<router-view />
	</div>-->
</template>

<style>
</style>
