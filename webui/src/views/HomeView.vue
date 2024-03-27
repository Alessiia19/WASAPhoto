<template>
	<div>
	  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Example App</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
		  <span class="navbar-toggler-icon"></span>
		</button>
	  </header>
	  
	  <div class="container-fluid">
		<div class="row">
		  <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
			<div class="position-sticky pt-3 sidebar-sticky">
			  <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
				<span>General</span>
			  </h6>
			  <ul class="nav flex-column">
				<li class="nav-item">
				  <RouterLink to="/" class="nav-link">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home"/></svg>
					Home
				  </RouterLink>
				</li>
				<li class="nav-item">
				  <RouterLink to="/link1" class="nav-link">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#layout"/></svg>
					Menu item 1
				  </RouterLink>
				</li>
				<li class="nav-item">
				  <RouterLink to="/link2" class="nav-link">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#key"/></svg>
					Menu item 2
				  </RouterLink>
				</li>
			  </ul>
  
			  <h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
				<span>Secondary menu</span>
			  </h6>
			  <ul class="nav flex-column">
				<li class="nav-item">
				  <RouterLink :to="'/some/' + 'variable_here' + '/path'" class="nav-link">
					<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#file-text"/></svg>
					Item 1
				  </RouterLink>
				</li>
			  </ul>
			</div>
		  </nav>
  
		  <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
			<!-- Contenuto della pagina Home -->
			<div>
			  <div
				class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
				<h1 class="h2">Home page</h1>
				<div class="btn-toolbar mb-2 mb-md-0">
				  <div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
					  Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
					  Export
					</button>
				  </div>
				  <div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
					  New
					</button>
				  </div>
				</div>
			  </div>
  
			  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
			</div>
		  </main>
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
  </style>
  