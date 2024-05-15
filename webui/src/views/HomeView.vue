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
			photos: [
				{
					photoID: 0,
					authorID: 0,
					authorUsername: '',
					imageData: '',
					uploadDate: '',
					likesCount: 0,
					commentsCount: 0,
					isLiked: false,
				}
			],
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

		async clearSearch() {
			this.searchQuery = '';
			this.users = [];
		},

		async goToUserProfile(userToSearch) {
			if (userToSearch) {
				localStorage.setItem("userToSearchID", userToSearch.userID)
				this.$router.push({ path: `/users/${userToSearch.username}` });
			} else {
				console.error('Errore: Nome utente non fornito.');
			}
		},

		async loadStreamData() {
			try {
				let response = await this.$axios.get('/users/' + this.userID + '/stream', {
					headers: {
						Authorization: "Bearer " + this.userID
					}
				});
				this.photos = response.data;
				console.log(this.photos);

			} catch (error) {
				console.error('Error while retrieving user stream: ', error);
			}
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

	},

	mounted() {
		this.loadStreamData();
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

					<!-- Search bar -->
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

				<!-- Stream photos -->
				<div class="stream-photo-card" v-for="photo in this.photos" :key="photo.photoID">
					<div class="stream-photo-image-container">
						<img :src="'data:image/jpeg;base64,' + photo.imageData" alt="Photo by {{ photo.authorUsername }}" />
					</div>
					<div class="stream-photo-info">
						<h3>{{ photo.authorUsername }}</h3>
						<p>Likes: {{ photo.likesCount }}</p>
						<p>Comments: {{ photo.commentsCount }}</p>
						<svg v-if="photo.isLiked" @click="toggleLike(photo)" xmlns="http://www.w3.org/2000/svg"
							width="16" height="16" fill="currentColor" class="bi bi-heart-fill" viewBox="0 0 16 16">
							<path fill-rule="evenodd"
								d="M8 1.314C12.438-3.248 23.534 4.735 8 15-7.534 4.736 3.562-3.248 8 1.314" />
						</svg>
						<svg v-else @click="toggleLike(photo)" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
							fill="currentColor" class="bi bi-heart" viewBox="0 0 16 16">
							<path
								d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143q.09.083.176.171a3 3 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15" />
						</svg>
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

.stream-photo-card {
	display: flex;
    align-items: center;
    background: #fff;
    border-radius: 15px; 
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.25); 
    margin-bottom: 20px; 
    padding: 10px;
    transition: transform 0.2s ease-in-out;
}

.stream-photo-image-container {
	flex: 1 1 auto;
    width: 100%;
    height: 390px;
    margin: 10px;
    border-radius: 15px;
    overflow: hidden;
}

.stream-photo-image-container img {
	width: 380px;
	height: 100%;
	object-fit: cover;
}

.stream-photo-info {
	flex: 2 1 50%;
	display: flex;
	flex-direction: column;
}

.stream-photo-info>* {
	margin-bottom: 10px;
}


.top-rounded {
	border-radius: 20px 20px 0px 0px;
}

.unselectable {
	user-select: none;
}
</style>