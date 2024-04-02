<script>
import { RouterLink } from 'vue-router'

export default {
data() {
	return {
	errormsg: null,
	username: "",
	followersCount: 0,
	followingCount: 0,
	photos: 0,
	}
},

methods: {

	async getUserInfo(){
		this.username = response.data.username;
		this.followersCount = response.data.followersCount;
		this.followingCount = response.data.followingCount;
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
		
		<!-- Profile area -->
		<main class="container-fluid">
			<div class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<div class="profile-container">
					<div class="profile-photo"></div>
					<div class="profile-info">
						<p class="username">this.username</p>
						<p><span class="info-label">Followers:</span> <span class="info-value">this.followersCount</span></p>
						<p><span class="info-label">Following:</span> <span class="info-value">this.followingCount</span></p>
						<p><span class="info-label">Photos:</span> <span class="info-value">this.uploadedPhotos</span></p>
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
    color: #00264d; /* Colore del titolo dell'applicazione */
    font-weight: bold; /* Rende il titolo in grassetto */
    font-size: 26px; /* Dimensione del titolo dell'applicazione */
	margin-left: 17px;
	margin-top: 7px;
  }
  
  .unselectable {
    user-select: none;
  }
  
  .profile-container {
	display: flex;
	padding: 20px;
  }

  .profile-photo {
	width: 300px;
	height: 300px;
	border: 2px solid #0000002f;
	margin-right: 20px;
  }
  
  .profile-info {
	display: flex;
	flex-direction: column;
	justify-content: center;
  }

  .username {
	font-size: 24px;
	font-weight: bold;
	text-decoration: underline;
	margin-bottom: 10px;
  }

  .info-label {
	font-weight: bold;
  }

  .info-value {
	margin-left: 5px;
  }
  
</style>
  