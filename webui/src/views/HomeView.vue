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
					username: '',
					imageData: '',
					uploadDate: '',
					likesCount: 0,
					likes: [],
					commentsCount: 0,
					comments: [],
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

		async likePhoto(photo) {
			try {
				let response = await this.$axios.post('/users/' + this.userID + '/photos/' + photo.photoID + '/likes', {}, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				photo.isLiked = true;
				this.loadStreamData();
			} catch (error) {
				console.error('Error while attempting to like the photo:', error);
			}
		},

		async loadStreamData() {
			try {
				let response = await this.$axios.get('/users/' + this.userID + '/stream', {
					headers: {
						Authorization: "Bearer " + this.userID
					}
				});
				this.photos = response.data
				if (this.photos) {
					this.photos.map(photo => {
						if (photo.likes) {
							photo.isLiked = photo.likes.some(like => like.userID === parseInt(this.userID));
							return photo;
						}
					});
				}

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

		async unlikePhoto(photo) {
			// Trova il like specifico fatto dall'utente loggato
			const userLike = photo.likes.find(like => like.userID === parseInt(this.userID));
			try {
				let response = await this.$axios.delete('/users/' + this.userID + '/photos/' + photo.photoID + '/likes/' + userLike.likeID, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				photo.isLiked = false;
				this.loadStreamData();
			} catch (error) {
				console.error('Error while attempting to unlike the photo:', error);
			}
		}


	},

	mounted() {
		this.loadStreamData();
		document.addEventListener('click', (event) => {
			if (!this.$el.contains(event.target)) {
				this.clearSearch();
			}
		});
	},

	beforeUnmount() {
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
				<div class="stream-content">

					<!-- Photo card -->
					<div class="stream-photo-card" v-for="photo in this.photos" :key="photo.photoID">
						<div class="stream-photo-image-container">
							<img :src="'data:image/jpeg;base64,' + photo.imageData"
								alt="Photo by {{ photo.username }}" />
						</div>

						<!-- Photo infos -->
						<div class="stream-photo-info">
							<div class="photo-author">{{ photo.username }}</div>
							<div class="photo-engagement-stats">
								<div class="heart-icon">
									<svg v-if="photo.isLiked" @click="unlikePhoto(photo)"
										xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
										class="bi bi-heart-fill" viewBox="0 0 16 16">
										<path fill-rule="evenodd"
											d="M8 1.314C12.438-3.248 23.534 4.735 8 15-7.534 4.736 3.562-3.248 8 1.314" />
									</svg>
									<svg v-else @click="likePhoto(photo)" xmlns="http://www.w3.org/2000/svg" width="16"
										height="16" fill="currentColor" class="bi bi-heart" viewBox="0 0 16 16">
										<path
											d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143q.09.083.176.171a3 3 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15" />
									</svg>
								</div>
								<p>{{ photo.likesCount }} Likes</p>
								<div class="comment-icon">
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
										class="bi bi-chat" viewBox="0 0 16 16">
										<path
											d="M2.678 11.894a1 1 0 0 1 .287.801 11 11 0 0 1-.398 2c1.395-.323 2.247-.697 2.634-.893a1 1 0 0 1 .71-.074A8 8 0 0 0 8 14c3.996 0 7-2.807 7-6s-3.004-6-7-6-7 2.808-7 6c0 1.468.617 2.83 1.678 3.894m-.493 3.905a22 22 0 0 1-.713.129c-.2.032-.352-.176-.273-.362a10 10 0 0 0 .244-.637l.003-.01c.248-.72.45-1.548.524-2.319C.743 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.52.263-1.639.742-3.468 1.105" />
									</svg>
								</div>
								<p>{{ photo.commentsCount }} Comments</p>
							</div>
						</div>
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

.comment-icon {
	margin-bottom: 15px;
}

.header-home-content {
	display: flex;
	align-items: center;
	justify-content: space-between;
	margin-bottom: 20px;
}

.heart-icon {
	margin-bottom: 13px;
}

.heart-icon:hover {
	transform: scale(1.2);
	transition: 0.1s ease-in-out;
}

.photo-author {
	font-weight: bold;
	font-size: 18px;
	width: 100%;
}

.photo-engagement-stats {
	display: flex;
	align-items: center;
	justify-content: start;
	gap: 10px;
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

.stream-content {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: 100%;
}

.stream-photo-card {
	display: flex;
	flex-direction: row;
	align-items: center;
	background: #fff;
	border-radius: 30px;
	box-shadow: 0 2px 12px rgba(0, 0, 0, 0.25);
	padding: 10px;
	width: 1000px;
	margin-left: auto;
	margin-right: auto;
	margin-bottom: 30px;
	margin-top: 10px;

}

.stream-photo-image-container {
	flex: 1 1 auto;
	width: 380px;
	height: 380px;
	margin: 15px;
	border-radius: 15px;
	overflow: hidden;
}

.stream-photo-image-container img {
	width: 100%;
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