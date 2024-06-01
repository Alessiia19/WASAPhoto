<script>
import { RouterLink } from 'vue-router'
import defaultProfilePic from '@/assets/user_icon.svg'

export default {
	data: function () {
		return {
			activeCommentMenu: null,
			defaultProfilePic,
			showDeleteModal: false,
			isEditingUsername: false,
			usernameWasModified: false,
			isMyProfile: false,
			isFollowed: false,
			isBanned: false,
			errormsg: null,
			isPhotoPopupOpen: false,
			newComment: '',
			selectedPhoto: {
				photoID: 0,
				userID: 0,
				imageData: '',
				uploadDate: '',
				likesCount: 0,
				likes: [],
				commentsCount: 0,
				comments: [
					{
						commentID: 0,
						authorID: 0,
						authorUsername: '',
						photoID: 0,
						commentText: '',
						uploadDate: '',
					}
				],
				isLiked: false,
			},
			loading: false,
			userID: localStorage.getItem('userID'),
			username: localStorage.getItem('username'),
			userToSearchID: localStorage.getItem('userToSearchID'),
			userProfile: {
				username: '',
				followers: [],
				following: [],
				followersCount: 0,
				followingCount: 0,
				uploadedPhotos: [
					{
						photoID: 0,
						userID: 0,
						imageData: '',
						uploadDate: '',
						likesCount: 0,
						likes: [],
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
				uploadedPhotosCount: 0,
			}

		}
	},

	async mounted() {
		this.checkIfOwnProfile();
		await this.loadProfileData();

	},

	watch: {
		'$route.params.username': 'handleRouteChange'
	},

	methods: {

		async banUser() {
			try {
				let response = await this.$axios.post('/users/' + this.userID + '/banned_users', { userid: parseInt(this.userToSearchID) }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isBanned = true;
				this.isFollowed = false;
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to ban:', error);
			}
		},

		async checkIfOwnProfile() {
			const routeUsername = this.$route.params.username;
			this.isMyProfile = (routeUsername === this.username);

		},

		closeDeleteModal() {
			this.showDeleteModal = false;
		},

		closePhotoPopup() {
			this.isPhotoPopupOpen = false;
			this.selectedPhoto = null;
		},

		confirmDelete() {
			if (this.selectedPhoto) {
				this.deletePhoto(this.selectedPhoto.photoID);
				this.closeDeleteModal();
			}
		},

		async deleteComment(photoID, commentID) {
			try {
				await this.$axios.delete('/users/' + this.userID + '/photos/' + photoID + '/comments/' + commentID, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.selectedPhoto.commentsCount -= 1;
				this.loadProfileData();
			} catch (error) {
				console.error('Error deleting comment:', error);
			}
		},

		async deletePhoto(photoID) {
			try {
				await this.$axios.delete('/users/' + this.userID + '/photos/' + photoID, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.closePhotoPopup();
				this.loadProfileData();
			} catch (error) {
				console.error('Error deleting photo:', error);
			}
		},

		enableEditing() {
			this.isEditingUsername = true;
		},

		async followUser() {
			try {
				let response = await this.$axios.post('/users/' + this.userID + '/following', { userid: parseInt(this.userToSearchID) }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isFollowed = true;
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to follow:', error);
			}
		},

		formatDate(value) {
			if (value) {
				return new Date(value).toLocaleDateString('en-US', {
					weekday: 'long', // "Monday"
					year: 'numeric', // "2021"
					month: 'long', // "July"
					day: 'numeric', // "19"
					hour: '2-digit', // "02"
					minute: '2-digit', // "00"
					hour12: true // use AM/PM
				});
			}
			return '';
		},

		async handleRouteChange() {
			if (!this.usernameWasModified) {
				this.checkIfOwnProfile();
				this.loadProfileData();
			}
			else {
				this.usernameWasModified = false;
			}
		},

		handleUsernameInput() {
			// Chiamata alla funzione di validazione durante ogni inserimento
			this.validateUsername();
		},

		async likePhoto(photo) {
			try {
				let response = await this.$axios.post('/users/' + this.userID + '/photos/' + photo.photoID + '/likes', {}, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				photo.isLiked = true;
				photo.likesCount += 1;
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to like the photo:', error);
			}
		},

		async loadProfileData() {
			if (this.isMyProfile) {
				try {
					let response = await this.$axios.get('/users/' + this.userID, {
						headers: {
							Authorization: "Bearer " + this.userID
						}
					});
					this.userProfile = response.data;
					this.username = response.data.username;
					localStorage.setItem("username", this.username)
					this.$router.push({ path: '/users/' + this.username })

				} catch (error) {
					console.error('Errore nel recupero dei dati del profilo:', error);
				}
			}

			else if (!this.isMyProfile) {
				try {
					// Verifica se l'utente attuale è stato bannato dall'utente di cui si sta cercando il profilo
					let responseBan = await this.$axios.get('/users/' + this.userToSearchID + '/banned_users/' + this.userID, {
						headers: {
							Authorization: "Bearer " + this.userToSearchID
						}
					});
					if (responseBan.data) {
						this.$router.replace({ path: '/not-found' });
						return;
					}
				} catch (error) {
					console.error('Errore durante il controllo del ban status:', error);
				}

				try {
					let response = await this.$axios.get('/users/' + this.userToSearchID, {
						headers: {
							Authorization: "Bearer " + this.userID
						}
					});
					this.userProfile = response.data;
					this.username = response.data.username;
					if (this.userProfile.followers) {
						this.isFollowed = this.userProfile.followers.some(follower => follower.userID === parseInt(this.userID));
					}
					if (this.userProfile.uploadedPhotos) {
						this.userProfile.uploadedPhotos.map(photo => {
							if (photo.likes) {
								photo.isLiked = photo.likes.some(like => like.userID === parseInt(this.userID));
								return photo;
							}
						});
					}

					try {
						let responseBan = await this.$axios.get('/users/' + this.userID + '/banned_users/' + this.userToSearchID, {
							headers: {
								Authorization: "Bearer " + this.userID
							}
						});
						this.isBanned = responseBan.data
					} catch (error) {
						console.error('Errore durante il controllo del ban status:', error);
					}

					this.$router.push({ path: '/users/' + this.$route.params.username })

				} catch (error) {
					console.error('Errore nel recupero dei dati del profilo:', error);
				}
			}
			if (this.selectedPhoto) {
				this.updateSelectedPhoto(this.selectedPhoto.photoID)
			}

			if (this.$route.params.username != this.username && !this.usernameWasModified) {
				window.location.reload();
			}
			else if (this.usernameWasModified) {
				this.usernameWasModified = false;
			}
		},

		async openDeleteModal() {
			this.showDeleteModal = true;
		},

		openPhotoPopup(photo) {
			this.selectedPhoto = photo;
			this.updateSelectedPhoto(this.selectedPhoto.photoID);
			this.isPhotoPopupOpen = true;
		},

		async postComment(photo) {
			if (!this.newComment.trim()) {
				return;
			}
			try {
				// Inviare il nuovo commento al server
				let response = await this.$axios.post('/users/' + this.userID + '/photos/' + photo.photoID + '/comments', { commentText: this.newComment }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.newComment = ''; // Resetta l'input
				photo.commentsCount += 1;
				this.loadProfileData();
			} catch (error) {
				console.error('Error posting comment:', error);
			}
		},

		async setMyUserName() {
			try {
				let response = await this.$axios.put('/users/' + this.userID, { username: this.userProfile.username }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isEditingUsername = false;
				this.loadProfileData();
				this.errormsg = null;
				this.usernameWasModified = true;
			} catch (error) {
				console.error('Errore nel salvare il nuovo username:', error);
				this.errormsg = error.response && error.response.data.message ? error.response.data.message : 'Username not valid';
			}

		},

		toggleCommentMenu(commentID) {
			this.activeCommentMenu = this.activeCommentMenu === commentID ? null : commentID;
		},

		async unbanUser() {
			try {
				let response = await this.$axios.delete('/users/' + this.userID + '/banned_users/' + this.userToSearchID, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isBanned = false;
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to unban:', error);
			}
		},

		async unfollowUser() {
			try {
				let response = await this.$axios.delete('/users/' + this.userID + '/following/' + parseInt(this.userToSearchID), {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isFollowed = false;
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to unfollow:', error);
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
				photo.likesCount -= 1;
				this.loadProfileData();

			} catch (error) {
				console.error('Error while attempting to unlike the photo:', error);
			}

		},

		async updateSelectedPhoto(photoID) {
			if (this.userProfile.uploadedPhotos) {
				const updatedPhoto = this.userProfile.uploadedPhotos.find(p => p.photoID === photoID);
				if (updatedPhoto) {
					this.selectedPhoto.likes = updatedPhoto.likes;
					this.selectedPhoto.comments = updatedPhoto.comments;
					if (this.selectedPhoto.comments) {
						this.selectedPhoto.comments.forEach(comment => {
							comment.isMyComment = comment.authorID === parseInt(this.userID);
						});
					}
				}
			}
		},

		validateUsername() {
			const usernameRegex = /^[a-zA-Z0-9]+$/;
			if (!usernameRegex.test(this.userProfile.username.trim())) {
				this.errormsg = "L'username deve contenere solo lettere e numeri.";
				return false;
			} else if (this.userProfile.username.length < 3 || this.userProfile.username.length > 16) {
				this.errormsg = "L'username deve essere compreso tra 3 e 16 caratteri.";
				return false;
			} else {
				this.errormsg = null; // Resetta il messaggio di errore se la validazione ha successo
				return true;
			}
		},

	},

}

</script>

<template>
	<div>
		<!-- Header -->
		<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
			<h1 class="app-title unselectable">WASAPhoto</h1>
		</header>

		<!-- Profile area -->
		<main class="profile-area">

			<!-- User card -->
			<div class="profile-card">
				<div class="profile-photo">
					<img :src="defaultProfilePic" class="profile-image">
				</div>
				<div class="profile-info">

					<div class="username-section">

						<!-- Username -->
						<h2 class="username" v-if="!isEditingUsername">{{ userProfile.username }}
							<button v-if="isMyProfile" class="edit-icon-button" @click="enableEditing">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
									class="bi bi-pencil-fill" viewBox="0 0 16 16">
									<path
										d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.5.5 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11z" />
								</svg>
							</button>
						</h2>

						<!-- Change username -->
						<div v-if="isEditingUsername" class="username-edit-section">
							<div class="input-and-button">
								<input type="text" v-model="userProfile.username" class="username-input"
									@input="handleUsernameInput">
								<button class="save-button" @click="setMyUserName" :disabled="errormsg">Save</button>
							</div>
							<div v-if="errormsg" class="text-danger">{{ errormsg }}</div>
						</div>

						<!-- Follow Button -->
						<button v-if="!isMyProfile && !isFollowed" class="follow-button" @click="followUser">+
							Follow</button>

						<!-- Unfollow Button -->
						<button v-if="!isMyProfile && isFollowed" class="unfollow-button"
							@click="unfollowUser">Unfollow</button>

						<!-- Ban Button -->
						<button v-if="!isMyProfile && !isBanned" class="ban-button" @click="banUser">Ban</button>

						<!-- Unban Button -->
						<button v-if="!isMyProfile && isBanned" class="unban-button" @click="unbanUser">Unban</button>

					</div>

					<div class="profile-stats">
						<p class="info">
							<span class="info-label">Followers:</span>
							<span class="info-value">{{ userProfile.followersCount }}</span>
						</p>
						<p class="info">
							<span class="info-label">Following:</span>
							<span class="info-value">{{ userProfile.followingCount }}</span>
						</p>
						<p class="info">
							<span class="info-label">Post:</span>
							<span class="info-value">{{ userProfile.uploadedPhotosCount }}</span>
						</p>
					</div>
				</div>
			</div>

			<!-- Photos area -->
			<div class="photos-grid" v-if="userProfile.uploadedPhotosCount > 0">
				<div v-for="photo in this.userProfile.uploadedPhotos" :key="photo.photoID" class="photo-card"
					@click="openPhotoPopup(photo)">
					<img :src="'data:image/jpeg;base64,' + photo.imageData" class="photo-img">
				</div>
			</div>

			<div v-if="userProfile.uploadedPhotosCount === 0" class="no-posts-container">
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" class="bi-camera"
					viewBox="0 0 16 16">
					<path
						d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4z" />
					<path
						d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5m0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7M3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0" />
				</svg>
				<p class="no-posts-text unselectable">No posts yet</p>
			</div>
		</main>

		<div v-if="showDeleteModal" class="modal">
			<div class="modal-content">
				<h4>Are you sure you want to delete this photo?</h4>
				<button @click="confirmDelete" class="delete-button">Confirm</button>
				<button @click="closeDeleteModal" class="cancel-button">Cancel</button>
			</div>
		</div>

		<!-- Photo popup -->
		<div v-if="isPhotoPopupOpen" class="photo-popup-overlay" @click="closePhotoPopup">
			<div class="photo-popup-card" @click.stop>
				<div class="photo-popup-image-container">
					<img :src="'data:image/jpeg;base64,' + selectedPhoto.imageData">
				</div>

				<!-- Photo popup infos -->
				<div class="photo-popup-info">
					<div class="photo-author-container">
						<div class="photo-author">
							<img class="author-image" :src="defaultProfilePic">{{ selectedPhoto.username }}
						</div>
						<div class="trash-icon" v-if="isMyProfile" @click="openDeleteModal()">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
								class="bi bi-trash3" viewBox="0 0 16 16">
								<path
									d="M6.5 1h3a.5.5 0 0 1 .5.5v1H6v-1a.5.5 0 0 1 .5-.5M11 2.5v-1A1.5 1.5 0 0 0 9.5 0h-3A1.5 1.5 0 0 0 5 1.5v1H1.5a.5.5 0 0 0 0 1h.538l.853 10.66A2 2 0 0 0 4.885 16h6.23a2 2 0 0 0 1.994-1.84l.853-10.66h.538a.5.5 0 0 0 0-1zm1.958 1-.846 10.58a1 1 0 0 1-.997.92h-6.23a1 1 0 0 1-.997-.92L3.042 3.5zm-7.487 1a.5.5 0 0 1 .528.47l.5 8.5a.5.5 0 0 1-.998.06L5 5.03a.5.5 0 0 1 .47-.53Zm5.058 0a.5.5 0 0 1 .47.53l-.5 8.5a.5.5 0 1 1-.998-.06l.5-8.5a.5.5 0 0 1 .528-.47M8 4.5a.5.5 0 0 1 .5.5v8.5a.5.5 0 0 1-1 0V5a.5.5 0 0 1 .5-.5" />
							</svg>
						</div>
					</div>

					<!-- Comment section -->
					<div class="photo-comments">
						<div class="comments-list">
							<div v-for="comment in selectedPhoto.comments" :key="comment.commentID"
								class="comment-container">
								<p class="comment-text">
									<strong>{{ comment.authorUsername }}:</strong> {{ comment.commentText }}
								</p>
								<div class="comment-menu-icon" v-if="comment.isMyComment"
									@click="toggleCommentMenu(comment.commentID)">
									<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
										class="bi bi-three-dots" viewBox="0 0 16 16">
										<path
											d="M3 9.5a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3m5 0a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3" />
									</svg>
									<ul v-if="activeCommentMenu === comment.commentID" class="dropdown-menu">
										<li @click="deleteComment(selectedPhoto.photoID, comment.commentID)">delete</li>
									</ul>
								</div>
							</div>
						</div>
						<div class="comment-input-container">
							<textarea v-model="newComment" placeholder="Add a comment..."
								class="comment-input"></textarea>
							<div @click="postComment(selectedPhoto)" class="send-icon">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
									class="bi bi-send-fill" viewBox="0 0 16 16">
									<path
										d="M15.964.686a.5.5 0 0 0-.65-.65L.767 5.855H.766l-.452.18a.5.5 0 0 0-.082.887l.41.26.001.002 4.995 3.178 3.178 4.995.002.002.26.41a.5.5 0 0 0 .886-.083zm-1.833 1.89L6.637 10.07l-.215-.338a.5.5 0 0 0-.154-.154l-.338-.215 7.494-7.494 1.178-.471z" />
								</svg>
							</div>
						</div>
					</div>

					<div class="photo-engagement-stats">
						<div v-if="!isMyProfile" class="photo-popup-heart-icon">
							<svg v-if="selectedPhoto.isLiked" @click="unlikePhoto(selectedPhoto)"
								xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
								class="bi bi-heart-fill" viewBox="0 0 16 16">
								<path fill-rule="evenodd"
									d="M8 1.314C12.438-3.248 23.534 4.735 8 15-7.534 4.736 3.562-3.248 8 1.314" />
							</svg>
							<svg v-else @click="likePhoto(selectedPhoto)" xmlns="http://www.w3.org/2000/svg" width="16"
								height="16" fill="currentColor" class="bi bi-heart" viewBox="0 0 16 16">
								<path
									d="m8 2.748-.717-.737C5.6.281 2.514.878 1.4 3.053c-.523 1.023-.641 2.5.314 4.385.92 1.815 2.834 3.989 6.286 6.357 3.452-2.368 5.365-4.542 6.286-6.357.955-1.886.838-3.362.314-4.385C13.486.878 10.4.28 8.717 2.01zM8 15C-7.333 4.868 3.279-3.04 7.824 1.143q.09.083.176.171a3 3 0 0 1 .176-.17C12.72-3.042 23.333 4.867 8 15" />
							</svg>
						</div>
						<div>{{ selectedPhoto.likesCount }} Likes</div>
						<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
							class="bi bi-chat" viewBox="0 0 16 16">
							<path
								d="M2.678 11.894a1 1 0 0 1 .287.801 11 11 0 0 1-.398 2c1.395-.323 2.247-.697 2.634-.893a1 1 0 0 1 .71-.074A8 8 0 0 0 8 14c3.996 0 7-2.807 7-6s-3.004-6-7-6-7 2.808-7 6c0 1.468.617 2.83 1.678 3.894m-.493 3.905a22 22 0 0 1-.713.129c-.2.032-.352-.176-.273-.362a10 10 0 0 0 .244-.637l.003-.01c.248-.72.45-1.548.524-2.319C.743 11.37 0 9.76 0 8c0-3.866 3.582-7 8-7s8 3.134 8 7-3.582 7-8 7a9 9 0 0 1-2.347-.306c-.52.263-1.639.742-3.468 1.105" />
						</svg>

						<div>{{ selectedPhoto.commentsCount }} Comments</div>
					</div>
					<!-- Upload Date -->
					<div class="photo-upload-date">{{ formatDate(selectedPhoto.uploadDate) }}</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.header {
	background-image: linear-gradient(to bottom right, #f5dd90, #b97b90, #446ca0);
	height: 70px;
}

.app-title {
	color: #00264d;
	font-weight: bold;
	font-size: 26px;
	margin-left: 17px;
	margin-top: 7px;
}

.ban-button {
	background: #df0c41;
}

.ban-button:hover {
	background: #c9302c;
}

.ban-button,
.unban-button {
	top: 20px;
	right: 20px;
	color: white;
	border: none;
	border-radius: 10px;
	padding: 8px 15px;
	cursor: pointer;
	font-weight: bold;
	height: 40px;
	margin-left: 20px;
}

.bi-camera {
	fill: #666;
	width: 50px;
	height: 50px;
}

.cancel-button {
	background-color: grey;
	color: white;
	border: none;
	padding: 10px 20px;
	margin: 5px;
	border-radius: 5px;
	cursor: pointer;
	font-weight: bold;
}

.cancel-button:hover {
	opacity: 0.75;
}

.delete-button {
	background-color: red;
	color: white;
	border: none;
	padding: 10px 20px;
	margin: 5px;
	border-radius: 5px;
	cursor: pointer;
	font-weight: bold;
}

.delete-button:hover {
	background-color: rgb(206, 4, 4);
}

.edit-icon-button {
	margin-left: 10px;
	background: none;
	border: none;
	cursor: pointer;
}

.edit-icon-button:hover {
	background: none;
}

.edit-icon-button:hover svg {
	transform: scale(1.2);
	transition: fill 0.3s ease-in-out;
}

.edit-icon-button svg {
	vertical-align: middle;
}

.follow-button {
	top: 20px;
	right: 20px;
	background: #446ca0;
	color: white;
	border: none;
	border-radius: 10px;
	padding: 8px 15px;
	cursor: pointer;
	font-weight: bold;
	height: 40px;
	margin-left: 20px;
}

.follow-button:hover {
	background: #365880;
}

.info {
	display: flex;
	width: 100%;
	font-size: 16px;
	margin: 5px 0;
}

.info-label {
	font-weight: bold;
	color: #444;
}

.info-value {
	margin-left: 5px;
}

.modal {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
}

.modal-content {
	background: white;
	padding: 20px;
	border-radius: 10px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
	text-align: center;
	width: 400px;
}

.photo-author-container {
	display: flex;
	justify-content: space-between
}

.photo-card {
	width: 380px;
	height: 380px;
	margin: 20px;
	background: #fff;
	border-radius: 15px;
	box-shadow: 0 2px 12px rgba(0, 0, 0, 0.2);
	overflow: hidden;
	position: relative;

}

.photo-card:hover {
	transform: scale(1.05);
	transition: 0.05s ease-in-out;
}

.photo-img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	/* Makes images cover the card area without distorting aspect ratio */
	position: absolute;
}

.photos-grid {
	display: flex;
	flex-wrap: wrap;
	justify-content: flex-start;
	gap: 20px;
	margin-top: 20px;
	width: 1300px;

}

.photo-popup-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background-color: rgba(0, 0, 0, 0.8);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
}

.photo-popup-card {
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

.photo-popup-heart-icon {
	margin-top: 5px;
}

.photo-popup-heart-icon:hover {
	transform: scale(1.2);
	transition: 0.1s ease-in-out;
}

.photo-popup-image-container {
	width: 550px;
	height: 550px;
	margin: 15px;
	border-radius: 15px;
	overflow: hidden;
}

.photo-popup-image-container img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.photo-popup-info {
	flex-grow: 1;
	display: flex;
	flex-direction: column;
	justify-content: space-around;
	margin-left: 10px;
	margin-right: 12px;
	width: 403px;
}

.photo-popup-info .comments-list {
	height: 350px;
}

.photo-popup-info .comment-text {
	max-width: 325px;
}

.profile-area {
	display: flex;
	flex-direction: column;
	align-items: center;
	width: calc(100%-280px);
	margin-left: 280px;
}

.profile-card {
	display: flex;
	align-items: center;
	background: #fff;
	border-radius: 15px;
	box-shadow: 0px 2px 12px rgba(0, 0, 0, 0.2);
	padding: 20px;
	width: 90%;
	margin: 20px auto;
	height: auto;
}

.profile-info {
	display: flex;
	flex-direction: column;
	flex-grow: 1;
	align-items: flex-start;
}

.profile-photo {
	width: 160px;
	height: 160px;
	background: #ddd;
	border-radius: 50%;
	margin-right: 30px;
	overflow: hidden;
	display: flex;
	align-items: center;
	justify-content: center;
}

.profile-photo img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.profile-stats {
	display: flex;
	justify-content: space-around;
	width: 100%;
	margin-left: 2px;
}

.no-posts-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-top: 225px;
}

.no-posts-text {
	font-size: 20px;
	color: #666;
	margin-top: 5px;
}

.save-button {
	background: #446ca0;
	color: white;
	border: none;
	border-radius: 10px;
	padding: 8px 15px;
	cursor: pointer;
	font-weight: bold;
	margin-left: 10px;
	width: auto;
	height: 40px;
}

.save-button:disabled {
	background-color: #7694bd;
	cursor: not-allowed;
	opacity: 0.5;
}

.save-button:hover {
	background: #365880;
}

.unban-button {
	background: rgb(3, 175, 3);
}

.unban-button:hover {
	background: rgb(3, 139, 3);
}

.unfollow-button {
	top: 20px;
	right: 20px;
	background: #df0c41;
	color: white;
	border: none;
	border-radius: 10px;
	padding: 8px 15px;
	cursor: pointer;
	font-weight: bold;
	height: 40px;
	margin-left: 20px;
}

.unfollow-button:hover {
	background: #cf0a0a;
}

.unselectable {
	user-select: none;
}

.username {
	font-size: 36px;
	font-weight: bold;
	margin-bottom: 20px;
	display: flex;
	flex-direction: row;
}

.username-edit-section {
	position: relative;
	display: flex;
	flex-direction: column;
	align-items: flex-start;
	width: 100%;
}

.username-edit-section .text-danger {
	position: absolute;
	top: 55%;
	left: 0;
	margin-top: 5px;
	margin-left: 2px;
}

.username-input {
	border: 2px solid #ccc;
	border-radius: 10px;
	padding: 8px 15px;
	width: auto;
	margin-bottom: 30px;
	margin-right: 3px;
}

.username-section {
	display: flex;
	width: 100%;
	text-align: center;
}
</style>