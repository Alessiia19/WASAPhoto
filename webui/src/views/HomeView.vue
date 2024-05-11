<script>
import { RouterLink } from 'vue-router'

export default {
	data() {
		return {
			errormsg: null,
			loading: false,
			userID: localStorage.getItem('userID'),
			username: localStorage.getItem('username'),
			searchQuery: "",
			users: [],
		}
	},

	methods: {
		async refresh() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},

		async searchUsers() {
			if (!this.searchQuery) {
				this.users = [];
				return;
			}
			try {
				const response = await this.$axios.get(`/users?username=${this.searchQuery}`);
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.toString();
				this.users = [];
			}
		},
	},

	watch: {
		searchQuery() {
			this.searchUsers(); // Chiama searchUsers ogni volta che searchQuery cambia
		}
	}
};
</script>

<template>
	<div>
		<!-- Header -->
		<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
			<h1 class="app-title unselectable">WASAPhoto</h1>
		</header>

		<main class="col-md-9 ms-sm-auto mt-3 col-lg-10 px-md-4 main-content">
			<div class="home-content">
				<div class="header-content">
					<h1 class="h2 unselectable" style="font-weight: bold; white-space: nowrap;">Home page</h1>
					<div class="search-container">
						<input v-model="searchQuery" type="text" class="form-control" placeholder="Search"
							aria-label="Search" name="searchbar">
						<ul v-if="users.length" class="search-dropdown">
							<li v-for="user in users" :key="user.userID" @click="goToUserProfile(user.userID)">
								{{ user.username }}
							</li>
						</ul>
					</div>
				</div>
			</div>
		</main>

	</div>
</template>


<style>
.header {
	background-image: linear-gradient(to bottom right, #f5dd90, #b97b90, #446ca0);
	height: 70px;
}


.app-title {
	color: #00264d;
	/* Colore del titolo dell'applicazione */
	font-weight: bold;
	/* Rende il titolo in grassetto */
	font-size: 26px;
	/* Dimensione del titolo dell'applicazione */
	margin-left: 17px;
	margin-top: 7px;
}

.header-content {
	display: flex;
	align-items: center;
	/* Allinea verticalmente al centro */
	justify-content: space-between;
	/* Spazia gli elementi uniformemente */
	margin-bottom: 20px;
	/* Distanza dalla parte inferiore */
}

.form-control {
	border-radius: 20px;
	/* Bordo arrotondato per la barra di ricerca */
	flex-grow: 1;
	/* Permette alla barra di ricerca di espandersi */
	margin-left: 20px;
	/* Distanza dalla scritta Home Page */
}

.search-container {
	position: relative;
	width: 100%;
	/* Larghezza massima possibile all'interno del suo contenitore */
}

.search-dropdown {
	list-style-type: none;
	margin: 0;
	padding: 0;
	background: white;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	position: absolute;
	width: 100%;
	top: 100%;
	left: 0;
	z-index: 1000;
	max-height: 300px;
	overflow-y: auto;
}

.search-dropdown li {
	padding: 8px 12px;
	cursor: pointer;
	border-bottom: 1px solid #ccc;
}

.search-dropdown li:hover {
	background-color: #f5f5f5;
}

.unselectable {
	user-select: none;
}
</style>