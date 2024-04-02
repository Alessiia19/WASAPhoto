<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
			isLogged: false,
			userID: localStorage.getItem('userID'),
		}
	},

	methods:{

		doLogout(loginStatus){
			console.log('Logging out:', loginStatus);
			this.isLogged = loginStatus;
			this.userID = null;
			this.$router.replace("/session");
		},

		changeLoginStatus(loginStatus){
			console.log('Changing login status:', loginStatus);
			this.isLogged = loginStatus;
		},


	},

	// Perform initial setup for the application
	created() {
		if (!localStorage.getItem('appInitialized')) {
			localStorage.clear();
			console.log("Contenuto attuale del localStorage:", localStorage);
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
					<Sidebar :currentRoute="$route.name" :userID="userID"/>
					
					<RouterView
					@changeLoginStatusmethod="changeLoginStatus"/>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
