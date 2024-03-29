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

<template>
	<div>
	  <!-- Header -->
	  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
		<h1 class="app-title unselectable">WASAPhoto</h1>
	  </header>
	  
	  <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
		  <!-- Home -->
		  <div>
			  <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
				  <h1 class="h2 unselectable" style="font-weight: bold; white-space: nowrap;">Home page</h1>
				  <div class="input-group input-group-sm" style="width: 1150px; margin-right: 10px;">
					  <div class="input-group-prepend">
						  <span class="input-group-text" style="border-top-right-radius: 0; border-bottom-right-radius: 0;">
							  <svg xmlns="http://www.w3.org/2000/svg" width="16" height="17" fill="currentColor" class="bi bi-search" viewBox="0 0 16 16">
								  <path d="M11.742 10.344a6.5 6.5 0 1 0-1.397 1.398h-.001q.044.06.098.115l3.85 3.85a1 1 0 0 0 1.415-1.414l-3.85-3.85a1 1 0 0 0-.115-.1zM12 6.5a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0"/>
							  </svg>
						  </span>
					  </div>
					  <input type="text" class="form-control" placeholder="Search" aria-label="Search" aria-describedby="button-search">
					  <button class="btn btn-outline-secondary" type="button" id="button-search">Search</button>
				  </div>
			  </div>
  
			  <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
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

  </style>
  