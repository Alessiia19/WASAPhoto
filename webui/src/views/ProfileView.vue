<template>
	<div>
	  <!-- Header -->
	  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
		<h1 class="app-title unselectable">WASAPhoto</h1>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
		  <span class="navbar-toggler-icon"></span>
		</button>
	  </header>
	  
	  <!-- Sidebar -->
	  <div class="container-fluid">
		<div class="row">
		  <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
			<div class="position-sticky pt-4 sidebar-sticky">
			  <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase"></h6>
			  <ul class="nav flex-column">
				<li class="nav-item">
				  <RouterLink to="/home" class="nav-link large-link"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-house-fill" viewBox="0 0 16 16">
						<path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L8 2.207l6.646 6.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293z"/>
						<path d="m8 3.293 6 6V13.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5V9.293z"/>
					</svg>
					<span class="link-text">Home</span> 
				  </RouterLink>
				</li>
				<li class="nav-item">
				  <RouterLink to="/link1" class="nav-link large-link"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-image-fill" viewBox="0 0 16 16">
						<path d="M.002 3a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2h-12a2 2 0 0 1-2-2zm1 9v1a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9.5l-3.777-1.947a.5.5 0 0 0-.577.093l-3.71 3.71-2.66-1.772a.5.5 0 0 0-.63.062zm5-6.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0"/>
					</svg>
					<span class="link-text">Create</span> 
				  </RouterLink>
				</li>
				<li class="nav-item">
				  <RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link large-link"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-person-fill" viewBox="0 0 16 16">
						<path d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6"/>
					</svg>
					<span class="link-text">Profile</span> 
				  </RouterLink>
				</li>
			  </ul>
  
			  <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase border-bottom"></h6>
			  <ul class="nav flex-column">
				<li class="nav-item">
				  <RouterLink to="/" class="nav-link large-link"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-box-arrow-right" viewBox="0 0 16 16">
						<path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0z"/>
						<path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z"/>
					</svg>
					<span class="link-text">Logout</span> 
				  </RouterLink>
				</li>
			  </ul>
			</div>
		  </nav>

		</div>
	  </div>
	</div>
  </template>
  
  <script>
  import { RouterLink } from 'vue-router'
  
  export default {
	data() {
	  return {
		errormsg: null,
		loading: false,
		some_data: null,
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
	},
  };
  </script>
  
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
  
  /* Regole per i link nella sidebar */
  .large-link {
	font-size: 18px; /* Dimensione del testo più grande */
	font-weight: bold; /* Testo in grassetto */
	margin-bottom: 10px; /* Spazio tra i link */
  }
  
  .large-link .link-text {
	margin-left: 10px; /* Spazio tra l'icona e il testo */
  }

  /* Regole per gestire l'effetto di hover sui link nella sidebar */
  .nav-item:hover {
    background-color: #f0f0f0; /* Cambia il colore di sfondo al passaggio del mouse */
  }
  
  .nav-item:hover .nav-link {
    font-weight: bold;
  }
  
  .unselectable {
    user-select: none;
  }

  </style>
  