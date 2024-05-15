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
			if (!this.searchQuery.trim()) {
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

		goToUserProfile(userToSearch) {
			if (userToSearch) {
				localStorage.setItem("userToSearchID", userToSearch.userID)
				this.$router.push({ path: `/users/${userToSearch.username}` });
			} else {
				console.error('Errore: Nome utente non fornito.');
			}
		},

		clearSearch() {
			this.searchQuery = '';
			this.users = [];
		},
	},

	mounted() {
		document.addEventListener('click', (event) => {
			if (!this.$el.contains(event.target)) {
				this.clearSearch();
			}
		});
	},

	beforeDestroy() {
		document.removeEventListener('click', this.clearSearch);
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
				<div class="header-home-content border-bottom pb-3">
					<h1 class="h2 unselectable" style="font-weight: bold; white-space: nowrap; margin-right: 15px;">Home
					</h1>
					<div class="search-container" @click.stop>
						<input v-model="searchQuery" type="text" class="form-control custom"
							:class="{ 'top-rounded': users.length, 'all-rounded': !users.length }" placeholder="Search"
							aria-label="Search" name="searchbar" autocomplete="off" @click.stop>
						<ul v-if="users.length" class="search-dropdown pl-3">
							<li v-for="user in users" :key="user.userID" @click="goToUserProfile(user)">
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

.all-rounded {
	border-radius: 20px;
}

.app-title {
	color: #00264d;
	font-weight: bold;
	font-size: 26px;
	margin-left: 17px;
	margin-top: 7px;
}

.header-home-content {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 20px;
}

.search-container {
	position: relative;
	width: 100%;
	transition: all 0.3s;
}

.search-container .custom {
	border: 1px solid #ccc;
	width: 100%;
	box-sizing: border-box;
	transition: border-color 0.2s, box-shadow 0.2s;
}

.search-container:focus-within .custom,
.search-container.show-results .custom {
	box-shadow: 0 2px 6px rgba(0, 0, 0, 0.25);
}

.search-dropdown {
	list-style-type: none;
	margin: 0;
	padding: 0;
	background: white;
	position: absolute;
	width: 100%;
	top: 100%;
	left: 0;
	z-index: 1000;
	max-height: 300px;
	overflow-y: auto;
	border: 1px solid #ccc;
	border-top: none;
	border-bottom: none;
	border-radius: 0 0 20px 20px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.25);
}

.search-dropdown li {
	padding: 8px 12px;
	cursor: pointer;
	border-bottom: 1px solid #ccc;
}

.search-dropdown li:hover {
	background-color: #f5f5f5;
}

.top-rounded {
	border-radius: 20px 20px 0px 0px;
}

.unselectable {
	user-select: none;
}
</style>