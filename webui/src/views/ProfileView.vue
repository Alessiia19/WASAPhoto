<script>
import { RouterLink } from 'vue-router'

export default {
	data: function() {
		return {
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
		await this.fetchProfileData();	
	},
	
	methods: {
		async fetchProfileData() {
			try {
				let response = await this.$axios.get('/users/' + this.userID ,{
						headers: {
							Authorization: "Bearer " + this.userID
						}
					});
				this.userProfile = response.data;
				console.log("Profilo:", response.data)
			} catch (error) {
				console.error('Errore nel recupero dei dati del profilo:', error);
			}
		},

		async setMyUserName() {},
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
	  <main class="profile-card">
		  <div class="profile-photo"></div>
		  <div class="profile-info">
			  <h2 class="username">{{ userProfile.username }}
				<button class="edit-icon-button" @click="setMyUserName">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-pencil-fill" viewBox="0 0 16 16">
						<path d="M12.854.146a.5.5 0 0 0-.707 0L10.5 1.793 14.207 5.5l1.647-1.646a.5.5 0 0 0 0-.708zm.646 6.061L9.793 2.5 3.293 9H3.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.5h.5a.5.5 0 0 1 .5.5v.207zm-7.468 7.468A.5.5 0 0 1 6 13.5V13h-.5a.5.5 0 0 1-.5-.5V12h-.5a.5.5 0 0 1-.5-.5V11h-.5a.5.5 0 0 1-.5-.5V10h-.5a.5.5 0 0 1-.175-.032l-.179.178a.5.5 0 0 0-.11.168l-2 5a.5.5 0 0 0 .65.65l5-2a.5.5 0 0 0 .168-.11z"/>
					</svg>
				</button>	
			  </h2>
			  <p class="info"><span class="info-label">Followers:</span> <span class="info-value">{{ userProfile.followersCount }}</span></p>
			  <p class="info"><span class="info-label">Following:</span> <span class="info-value">{{ userProfile.followingCount }}</span></p>
			  <p class="info"><span class="info-label">Uploads:</span> <span class="info-value">{{ userProfile.uploadedPhotosCount }}</span></p>
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
	  background: #ddd; /* Placeholder, replace with image */
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
  
  