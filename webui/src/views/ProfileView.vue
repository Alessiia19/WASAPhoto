<script>
import { RouterLink } from 'vue-router'

export default {
	data: function () {
		return {
			isEditingUsername: false,
			usernameWasModified: false,
			isMyProfile: false,
			isFollowed: false,
			errormsg: null,
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
						commentsCount: 0,
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

		async checkIfOwnProfile() {
			const routeUsername = this.$route.params.username;
			this.isMyProfile = (routeUsername === this.username);

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
				console.log("Following user:", this.userToSearchID);
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to follow:', error);
			}
		},

		async handleRouteChange() {
			console.log("Username modified: ", this.usernameWasModified);
			if (!this.usernameWasModified) {
				this.checkIfOwnProfile();
				this.loadProfileData();
				console.log("OK")
			}
			else {
				this.usernameWasModified = false;
			}
		},

		handleUsernameInput() {
			// Chiamata alla funzione di validazione durante ogni inserimento
			this.validateUsername();
		},

		async loadProfileData() {
			console.log("Inizio caricamento dati profilo, isMyProfile:", this.isMyProfile);
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
					let response = await this.$axios.get('/users/' + this.userToSearchID, {
						headers: {
							Authorization: "Bearer " + this.userID
						}
					});
					this.userProfile = response.data;
					this.username = response.data.username;
					console.log(this.userID)
					console.log(this.userProfile)
					if (this.userProfile.followers){
						this.isFollowed = this.userProfile.followers.some(follower => follower.userID === parseInt(this.userID));
					}
					console.log(this.isFollowed)
					this.$router.push({ path: '/users/' + this.$route.params.username })

				} catch (error) {
					console.error('Errore nel recupero dei dati del profilo:', error);
				}
			}
			if (this.$route.params.username != this.username && !this.usernameWasModified) {
				window.location.reload();
			}
			else if (this.usernameWasModified) {
				this.usernameWasModified = false;
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

		async unfollowUser() {
			try {
				let response = await this.$axios.delete('/users/' + this.userID + '/following/' + parseInt(this.userToSearchID), {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isFollowed = false;
				console.log("Unfollowing user:", this.userToSearchID);
				this.loadProfileData();
			} catch (error) {
				console.error('Error while attempting to unfollow:', error);
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
				<div class="profile-photo"></div>
				<div class="profile-info">
					<div class="username-section">
						<h2 class="username" v-if="!isEditingUsername">{{ userProfile.username }}
							<button v-if="isMyProfile" class="edit-icon-button" @click="enableEditing">
								<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
									class="bi bi-pencil-fill" viewBox="0 0 16 16">
									<path
										d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.5.5 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11z" />
								</svg>
							</button>
						</h2>
						<input v-else type="text" v-model="userProfile.username" class="username-input"
							@input="handleUsernameInput">
						<button v-if="isEditingUsername" class="save-button" @click="setMyUserName">Save</button>
						<div v-if="errormsg" class="text-danger">{{ errormsg }}</div>

						<!-- Follow Button -->
						<button v-if="!isMyProfile && !isFollowed" class="follow-button" @click="followUser">+ Follow</button>

						<!-- Unfollow Button -->
						<button v-if="!isMyProfile && isFollowed" class="unfollow-button" @click="unfollowUser">Unfollow</button>
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
				<div v-for="photo in this.userProfile.uploadedPhotos" :key="photo.photoID" class="photo-card">
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


.bi-camera {
	fill: #666;
	width: 50px;
	height: 50px;
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
}

.profile-stats {
	display: flex;
	justify-content: space-around;
	width: 100%;
	margin-top: 10px;
	margin-left: 2px;
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
	transition: fill 0.05s ease-in-out;	
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

.save-button:hover {
	background: #365880;
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