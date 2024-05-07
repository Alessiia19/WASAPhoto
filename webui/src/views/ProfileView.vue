<script>
import { RouterLink } from 'vue-router'

export default {
	data: function () {
		return {
			isEditingUsername: false,
			errormsg: null,
			loading: false,
			userID: localStorage.getItem('userID'),
			username: localStorage.getItem('username'),
			userProfile: {
				username: '',
				followersCount: 0,
				followingCount: 0,
				uploadedPhotosCount: 0,
			}
		}
	},

	async mounted() {
		await this.loadProfileData();
	},

	methods: {
		async loadProfileData() {
			try {
				let response = await this.$axios.get('/users/' + this.userID, {
					headers: {
						Authorization: "Bearer " + this.userID
					}
				});
				this.userProfile = response.data;
				this.username = response.data.username
				localStorage.setItem("username", this.username)
				this.$router.push({ path: '/users/' + this.username })

			} catch (error) {
				console.error('Errore nel recupero dei dati del profilo:', error);
			}
		},

		enableEditing() {
			this.isEditingUsername = true;
		},

		handleUsernameInput() {
			// Chiamata alla funzione di validazione durante ogni inserimento
			this.validateUsername();
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

		async setMyUserName() {
			try {
				let response = await this.$axios.put('/users/' + this.userID, { username: this.userProfile.username }, {
					headers: { Authorization: "Bearer " + this.userID }
				});
				this.isEditingUsername = false;
				this.loadProfileData();
				this.errormsg = null; // Reset dell'errore
			} catch (error) {
				console.error('Errore nel salvare il nuovo username:', error);
				this.errormsg = error.response && error.response.data.message ? error.response.data.message : 'Username not valid';
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
			<div class="profile-card">
				<div class="profile-photo"></div>
				<div class="profile-info">
					<h2 class="username" v-if="!isEditingUsername">{{ userProfile.username }}
						<button class="edit-icon-button" @click="enableEditing">
							<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
								class="bi bi-pencil-fill" viewBox="0 0 16 16">
								<path
									d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.5.5 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11z" />
							</svg>
						</button>
					</h2>
					<input v-else type="text" v-model="userProfile.username" class="username-input" required
						autocomplete="username" @input="handleUsernameInput">
					<button v-if="isEditingUsername" class="save-button" @click="setMyUserName">Save</button>
					<div v-if="errormsg" class="text-danger">{{ errormsg }}</div>

					<p class="info"><span class="info-label">Followers:</span> <span class="info-value">{{
						userProfile.followersCount }}</span></p>
					<p class="info"><span class="info-label">Following:</span> <span class="info-value">{{
						userProfile.followingCount }}</span></p>
					<p class="info"><span class="info-label">Uploads:</span> <span class="info-value">{{
						userProfile.uploadedPhotosCount }}</span></p>
				</div>
			</div>
			
			<div  v-if="userProfile.uploadedPhotosCount === 0" class="no-posts-container">
				<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" fill="currentColor" class="bi-camera" viewBox="0 0 16 16">
					<path d="M15 12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h1.172a3 3 0 0 0 2.12-.879l.83-.828A1 1 0 0 1 6.827 3h2.344a1 1 0 0 1 .707.293l.828.828A3 3 0 0 0 12.828 5H14a1 1 0 0 1 1 1zM2 4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2h-1.172a2 2 0 0 1-1.414-.586l-.828-.828A2 2 0 0 0 9.172 2H6.828a2 2 0 0 0-1.414.586l-.828.828A2 2 0 0 1 3.172 4z"/>
					<path d="M8 11a2.5 2.5 0 1 1 0-5 2.5 2.5 0 0 1 0 5m0 1a3.5 3.5 0 1 0 0-7 3.5 3.5 0 0 0 0 7M3 6.5a.5.5 0 1 1-1 0 .5.5 0 0 1 1 0"/>
				</svg>
				<p class="no-posts-text">No posts yet</p>
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

.profile-area {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.no-posts-container {
	display: flex;
	flex-direction: column;
	align-items: center;
	position: fixed;
	top: 475px;
	right: 35px;
	width: 100%;
	
}

.bi-camera {
	fill: #666; 
	width: 50px; 
	height: 50px;
}

.no-posts-text {
	font-size: 20px;
	color: #666;
	margin-top: 5px;
}

.profile-card {
	margin-left: auto;
	margin-right: 15px;
	margin-top: 15px;
	display: flex;
	flex-direction: column;
	align-items: center;
	background: #fff;
	border-radius: 15px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.2);
	padding: 20px;
	width: 340px;
}

.unselectable {
	user-select: none;
}

.profile-photo {
	width: 300px;
	height: 300px;
	background: #ddd;
	/* Placeholder, replace with image */
	border-radius: 50%;
	margin-bottom: 20px;
}

.profile-info {
	text-align: center;
}

.username {
	font-size: 36px;
	font-weight: bold;
	margin-bottom: 30px;
}

.username-input {
	border: 2px solid #ccc;
	border-radius: 10px;
	padding: 8px;
	width: 70%;
	margin-bottom: 10px;
	margin-right: 3px;
}

.save-button {
	background: #446ca0;
	color: white;
	border: none;
	border-radius: 10px;
	padding: 10px 20px;
	cursor: pointer;
	font-weight: bold;
}

.save-button:hover {
	background: #365880;
}

.edit-icon-button {
	background: none;
	border: none;
	cursor: pointer;
	padding: 0;
	display: inline;
}

.edit-icon-button:hover {
	background: none;
}

.edit-icon-button svg {
	vertical-align: middle;
}

.edit-icon-button:hover svg {
	transform: scale(1.1);
	transition: fill 0.3s ease-in-out;
}


.info {
	margin: 10px 0;
}

.info-label {
	font-weight: bold;
}

.info-value {
	margin-left: 5px;
}
</style>