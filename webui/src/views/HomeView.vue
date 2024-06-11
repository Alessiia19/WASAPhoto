<script>
import { RouterLink } from 'vue-router'
import { LottiePlayer } from '@lottiefiles/vue-lottie-player'
import animationData from '@/assets/home_animation.json'
import defaultProfilePic from '@/assets/user_icon.svg'

export default {
	components: {
		LottiePlayer
	},

	data() {
		return {
			activeCommentMenu: null,
			animationData,
			defaultProfilePic,
			errormsg: null,
			showBackToTop: false,
			loading: false,
			userID: localStorage.getItem('userID'),
			username: localStorage.getItem('username'),
			searchQuery: "",
			users: [],
			photos: [
				{
					photoID: 0,
					userID: 0,
					username: '',
					imageData: '',
					uploadDate: '',
					likesCount: 0,
					likes: [],
					newComment: '',
					commentsCount: 0,
					comments: [
						{
							commentID: 0,
							authorID: 0,
							authorUsername: '',
							photoID: 0,
							commentText: '',
							uploadDate: '',
							isMyComment: false,
						}
					],
					isLiked: false,
				}
			],
		}
	},

	methods: {

		backToTop() {
			window.scrollTo({
				top: 0,
				behavior: 'smooth'
			});
		},

		async clearSearch() {
			this.searchQuery = '';
			this.users = [];
		},

		async deleteComment(photoID, commentID) {
			try {
				await this.$axios.delete('/users/' + this.userID + '/photos/' + photoID + '/comments/' + commentID, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.loadStreamData();
			} catch (error) {
				console.error('Error deleting comment:', error);
			}
		},

		formatDate(value) {
			if (value) {
				return new Date(value).toLocaleDateString('en-US', {
					year: 'numeric', 
					month: 'long', 
					day: 'numeric', 
					hour: '2-digit', 
					minute: '2-digit', 
					hour12: true 
				});
			}
			return '';
		},

		async goToUserProfile(userID, username) {
			console.log(userID, username)
			if (userID) {
				localStorage.setItem("userToSearchID", userID)
				this.$router.push({ path: `/users/${username}` });
			} else {
				console.error('Errore: username not provided.');
			}
		},

		handleOutsideClick(event) {
			// Check if the click is outside the search bar and clear the search if true.
			const searchBar = this.$refs.searchBar;
			if (searchBar && !searchBar.contains(event.target)) {
				this.clearSearch();
			}
		},

		handleScroll() {
			// Show the "Back to Top" button when the user scrolls down 200 pixels.
			this.showBackToTop = window.scrollY > 200;
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
					this.photos.forEach(photo => {
						if (photo.comments) {
							// Check and update each comment to see if it belongs to the logged-in user.
							photo.comments.forEach(comment => {
								comment.isMyComment = comment.authorID === parseInt(this.userID);
							});
						}
						if (photo.likes) {
							// Check and update each photo to see if it is liked by the logged-in user
							photo.isLiked = photo.likes.some(like => like.userID === parseInt(this.userID));
						}
					});
				}
			} catch (error) {
				console.error('Error while retrieving user stream: ', error);
				this.photos = []
			}
		},

		async postComment(photo) {
			if (!photo.newComment.trim()) {
				return;
			}
			try {
				// Inviare il nuovo commento al server
				let response = await this.$axios.post('/users/' + this.userID + '/photos/' + photo.photoID + '/comments', { commentText: photo.newComment }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				photo.newComment = '';
				this.loadStreamData();
			} catch (error) {
				console.error('Error posting comment:', error);
			}
		},

		async searchUsers() {
			if (!this.searchQuery.trim()) {
				this.users = [];
				return;
			}
			try {
				const response = await this.$axios.get(`/users?username=${this.searchQuery}`, {
					headers: {
						Authorization: "Bearer " + this.userID
					}
				});
				this.users = response.data || [];
			} catch (e) {
				this.errormsg = e.toString();
				this.users = [];
			}
		},

		toggleCommentMenu(commentID) {
			this.activeCommentMenu = this.activeCommentMenu === commentID ? null : commentID;
		},

		async unlikePhoto(photo) {
			// Find the like by the logged-in user on the specified photo.
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
		},

	},

	mounted() {
		this.loadStreamData();
		document.addEventListener('click', this.handleOutsideClick);
		window.addEventListener('scroll', this.handleScroll);
	},

	beforeUnmount() {
		document.removeEventListener('click', this.handleOutsideClick);
		window.removeEventListener('scroll', this.handleScroll);
	},

	watch: {
		searchQuery() {
			this.searchUsers(); // Call searchUsers whenever the searchQuery changes.
		}
	}
};
</script>

<template>
	<div>
		<!-- Header -->
		<header class="navbar sticky-top flex-md-nowrap p-0 header">
			<h1 class="app-title unselectable">WASAPhoto</h1>
		</header>

		<!-- Main content area -->
		<main class="col-md-9 ms-sm-auto mt-3 col-lg-10 px-md-4 main-content">
			<div class="home-content">
				<div class="header-home-content border-bottom pb-3">
					<h1 class="h2 unselectable" style="font-weight: bold; white-space: nowrap; margin-right: 15px;">Home
					</h1>

					<!-- Search bar -->
					<div class="search-container" @click.stop ref="searchBar">
						<input v-model="searchQuery" type="text" class="form-control custom"
							:class="{ 'top-rounded': users.length, 'all-rounded': !users.length }" placeholder="Search"
							aria-label="Search" name="searchbar" autocomplete="off">
						<svg v-if="users.length" xmlns="http://www.w3.org/2000/svg" width="16" height="16"
							fill="currentColor" class="bi bi-x-circle-fill clear-search-icon" viewBox="0 0 16 16"
							@click.stop="clearSearch">
							<path
								d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0M5.354 4.646a.5.5 0 1 0-.708.708L7.293 8l-2.647 2.646a.5.5 0 0 0 .708.708L8 8.707l2.646 2.647a.5.5 0 0 0 .708-.708L8.707 8l2.647-2.646a.5.5 0 0 0-.708-.708L8 7.293z" />
						</svg>
						<ul v-if="users.length" class="search-dropdown pl-3">
							<li v-for="user in users" :key="user.userID" @click="goToUserProfile(user.userID, user.username)">
								{{ user.username }}
							</li>
						</ul>
					</div>
				</div>

				<!-- Stream content -->
				<!-- No content message -->
				<div v-if="!photos" class="no-content-message">
					<lottie-player :src="animationData" background="transparent" speed="0.5"
						style="width: 300px; height: 300px;" loop autoplay></lottie-player>
					<p class="no-content-text">There are no photos to display. Start following other users to see their
						photos here!</p>
				</div>

				<!-- Stream photos -->
				<div class="stream-content">

					<!-- Photo card -->
					<div class="stream-photo-card" v-for="photo in this.photos" :key="photo.photoID">
						<div class="stream-photo-image-container">
							<img :src="'data:image/jpeg;base64,' + photo.imageData">
						</div>

						<!-- Photo infos -->
						<div class="stream-photo-info">
							<!-- Photo author username  -->
							<div class="photo-author">
								<img class="author-image" :src="defaultProfilePic">
								<p class="photo-author-username" @click="goToUserProfile(photo.userID, photo.username)">
									{{ photo.username}}
								</p>
							</div>

							<!-- Comment section -->
							<div class="photo-comments">
								<div class="comments-list">
									<div v-for="comment in photo.comments" :key="comment.commentID"
										class="comment-container">

										<div style="display: flex; flex-direction: row; justify-content: space-between; width: 100%">

											<!-- Comment text -->
											<p class="comment-text">
												<strong>{{ comment.authorUsername }}:</strong> {{ comment.commentText }}
											</p>
											<!-- Delete comment section -->
											<div class="comment-menu-icon" v-if="comment.isMyComment"
												@click="toggleCommentMenu(comment.commentID)">
												<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
													fill="currentColor" class="bi bi-three-dots" viewBox="0 0 16 16">
													<path
														d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3" />
												</svg>
												<ul v-if="activeCommentMenu === comment.commentID"
													class="dropdown-menu">
													<li @click="deleteComment(photo.photoID, comment.commentID)">delete
													</li>
												</ul>
											</div>
										</div>


										<!-- Comment Upload Date -->
										<div class="photo-upload-date">{{ formatDate(comment.uploadDate) }}</div>
									</div>
								</div>

								<!-- Add comment section -->
								<div class="comment-input-container">
									<textarea v-model="photo.newComment" placeholder="Add a comment..."
										class="comment-input"></textarea>
									<div @click="postComment(photo)" class="send-icon">
										<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16"
											fill="currentColor" class="bi bi-send-fill" viewBox="0 0 16 16">
											<path
												d="M15.964.686a.5.5 0 0 0-.65-.65L.767 5.855H.766l-.452.18a.5.5 0 0 0-.082.887l.41.26.001.002 4.995 3.178 3.178 4.995.002.002.26.41a.5.5 0 0 0 .886-.083zm-1.833 1.89L6.637 10.07l-.215-.338a.5.5 0 0 0-.154-.154l-.338-.215 7.494-7.494 1.178-.471z" />
										</svg>
									</div>
								</div>
							</div>

							<!-- Engagement stats -->
							<div class="photo-engagement-stats">

								<!-- Like section -->
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

								<!-- Likes count -->
								<p>{{ photo.likesCount }} Likes</p>
								<div class="comment-icon">
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
										class="bi bi-chat" viewBox="0 0 16 16">
										<path
											d="M2.678 11.894a1 1 0 0 1 .287.801 11 11 0 0 1-.398 2c1.395-.323 2.247-.697 2.634-.893a1 1 0 0 1 .71-.074A8 8 0 0 0 8 14c3.996 0 7-2.807 7-6s-3.004-6-7-6-7 2.808-7 6c0 1.468.617 2.83 1.678 3.894m-.493 3.905a22 22 0 0 1-.713.129c-.2.032-.352-.176-.273-.362a10 10 0 0 0 .244-.637l.003-.01c.248-.72.45-1.548.524-2.319C.743 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.52.263-1.639.742-3.468 1.105" />
									</svg>
								</div>
								<!-- Comments count -->
								<p>{{ photo.commentsCount }} Comments</p>
							</div>

							<!-- Photo Upload Date -->
							<div class="photo-upload-date">{{ formatDate(photo.uploadDate) }}</div>
						</div>
					</div>
				</div>

			</div>
		</main>

		<!-- Back To Top Floating Button -->
		<button v-show="showBackToTop" @click="backToTop" class="back-to-top">
			<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-arrow-up"
				viewBox="0 0 16 16">
				<path fill-rule="evenodd"
					d="M8 15a.5.5 0 0 0 .5-.5V2.707l3.146 3.147a.5.5 0 0 0 .708-.708l-4-4a.5.5 0 0 0-.708 0l-4 4a.5.5 0 1 0 .708.708L7.5 2.707V14.5a.5.5 0 0 0 .5.5" />
			</svg>
		</button>
	</div>
</template>


<style>
.header {
	background-image: none;
	background-color: #fff;
	height: 70px;
	box-shadow: 0 3px 5px rgba(0, 0, 0, 0.2);
}

.all-rounded {
	border-radius: 20px;
}

.app-title {
	color: #333333;
	font-weight: bold;
	font-size: 26px;
	margin-left: 17px;
	margin-top: 7px;
}

.author-image {
	width: 30px;
	height: 30px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 10px;
}

.back-to-top {
	position: fixed;
	bottom: 30px;
	right: 30px;
	width: 50px;
	height: 50px;
	border-radius: 50%;
	background-color: #446ca0;
	color: white;
	border: none;
	box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.back-to-top:hover {
	background-color: #365880;
}

.clear-search-icon {
	position: absolute;
	right: 13px;
	top: 50%;
	transform: translateY(-50%);
	cursor: pointer;
	color: #ccc;
}

.clear-search-icon:hover {
	color: #000;
}

.comment-container {
	display: flex;
	flex-direction: column;
	margin-bottom: 20px;
	padding-right: 15px;
	position: relative;
}

.comment-container:hover .comment-menu-icon {
	visibility: visible;
}

.comment-icon {
	margin-bottom: 15px;
}

.comment-input {
	border-radius: 10px;
	padding: 10px;
	padding-right: 50px;
	width: 100%;
	border: 1px solid #ccc;
	resize: none;
}

.comment-input-container {
	position: relative;
	overflow-x: auto;
}

.comments-list {
	height: 140px;
	margin-bottom: 15px;
	margin-right: 2px;
	overflow-y: auto;
	word-wrap: break-word;
}

.comment-menu-icon {
	flex-shrink: 0;
	cursor: pointer;
	align-self: flex-start;
	margin-top: 2px;
	height: 16px;
	visibility: hidden;
}

.comment-text {
	flex: 1;
	word-wrap: break-word;
	max-width: 493px;
	margin-bottom: 0;
}

.dropdown-menu {
	list-style: none;
	position: absolute;
	background-color: #fff;
	border: 1px solid #ccc;
	right: 0;
	z-index: 100;
	padding: 5px 10px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	display: block;
}

.dropdown-menu li {
	cursor: pointer;
	color: red;
}

.dropdown-menu:hover {
	background-color: #f0f0f0;
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

.no-content-message {
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	height: 80vh;
	text-align: center;
}

.no-content-text {
	font-family: 'Fredoka', sans-serif;
	font-weight: 400;
	font-size: 24px;
	color: #444;
	margin-top: 20px;
}

.photo-author {
	font-weight: bold;
	font-size: 18px;
	width: 100%;
	display: flex;
	flex-direction: row;
	margin-bottom: 10px;
}

.photo-author-username {
	cursor: pointer;
}

.photo-comments {
	border-left: 2px solid #e0e0e0;
	border-right: 2px solid #e0e0e0;
	display: flex;
	flex-direction: column;
	margin-bottom: 10px;
}

.photo-engagement-stats {
	display: flex;
	align-items: center;
	justify-content: start;
	gap: 10px;
}

.photo-upload-date {
	font-size: 14px;
	color: #6c757d;
	font-style: italic;
	margin-top: 5px;
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

.send-icon {
	margin-right: 10px;
	border-radius: 5px;
	position: absolute;
	top: 15px;
	right: 15px;
	cursor: pointer;
	transition: transform 0.2s ease-in-out;
}

.send-icon:hover {
	transform: scale(1.2);
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
	width: 570px;
}

.stream-photo-info,
.photo-comments {
	flex: 1;
	padding: 10px;
}

.top-rounded {
	border-radius: 20px 20px 0px 0px;
}

.trash-icon:hover {
	color: red
}

.unselectable {
	user-select: none;
}
</style>