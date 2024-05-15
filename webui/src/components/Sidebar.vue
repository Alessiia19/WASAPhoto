<script>
	export default {
		props: {
			currentRoute: String
		},
		data() {
			return {
				sidebarRoutes: ['home', 'uploadPhoto', 'profile'], // Routes che devono mostrare la sidebar
				
			};
		},

		
		methods: {

			navigateTo(path) {
				const username = localStorage.getItem('username');
				if (!username) {
					console.error('Username non trovato nel localStorage');
					this.$router.push('/login'); // Reindirizza l'utente alla pagina di login
					return;
				}
				// Se il path contiene ":userID", sostituiscilo con l'userID effettivo
				const finalPath = path.replace(':username', username);
				this.$router.push(finalPath);
			},

			async logoutButton() {
				localStorage.removeItem('userID');
				localStorage.removeItem('username');
				this.$router.replace('/session');
			}
		} 
	};		
</script>


<template>
	<div v-if="sidebarRoutes.includes(currentRoute)">
		<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
			<div class="position-sticky pt-4 sidebar-sticky">
				<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase"></h6>
				<ul class="nav flex-column">
				
				<!-- Home button -->
				<li class="nav-item">
					<button class="nav-link large-link" @click="navigateTo('/users/:username/stream')"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-house-fill" viewBox="0 0 16 16">
						<path d="M8.707 1.5a1 1 0 0 0-1.414 0L.646 8.146a.5.5 0 0 0 .708.708L8 2.207l6.646 6.647a.5.5 0 0 0 .708-.708L13 5.793V2.5a.5.5 0 0 0-.5-.5h-1a.5.5 0 0 0-.5.5v1.293z"/>
						<path d="m8 3.293 6 6V13.5a1.5 1.5 0 0 1-1.5 1.5h-9A1.5 1.5 0 0 1 2 13.5V9.293z"/>
					</svg>
					<span class="link-text">Home</span> 
					</button>
				</li>

				<!-- Create button -->
				<li class="nav-item">
					<button class="nav-link large-link" @click="navigateTo('/link1')"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-image-fill" viewBox="0 0 16 16">
						<path d="M.002 3a2 2 0 0 1 2-2h12a2 2 0 0 1 2 2v10a2 2 0 0 1-2 2h-12a2 2 0 0 1-2-2zm1 9v1a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V9.5l-3.777-1.947a.5.5 0 0 0-.577.093l-3.71 3.71-2.66-1.772a.5.5 0 0 0-.63.062zm5-6.5a1.5 1.5 0 1 0-3 0 1.5 1.5 0 0 0 3 0"/>
					</svg>
					<span class="link-text">Create</span> 
					</button>
				</li>

				<!-- Profile button -->
				<li class="nav-item">
					<button class="nav-link large-link" @click="navigateTo('/users/:username')">
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-person-fill" viewBox="0 0 16 16">
						<path d="M3 14s-1 0-1-1 1-4 6-4 6 3 6 4-1 1-1 1zm5-6a3 3 0 1 0 0-6 3 3 0 0 0 0 6"/>
					</svg>
					<span class="link-text">Profile</span> 
					</button>
				</li>
				</ul>

				<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase border-bottom"></h6>
				<ul class="nav flex-column">
				
				<!-- Logout button -->
				<li class="nav-item">
					<button class="nav-link large-link" @click="logoutButton"> 
					<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-box-arrow-right" viewBox="0 0 16 16">
						<path fill-rule="evenodd" d="M10 12.5a.5.5 0 0 1-.5.5h-8a.5.5 0 0 1-.5-.5v-9a.5.5 0 0 1 .5-.5h8a.5.5 0 0 1 .5.5v2a.5.5 0 0 0 1 0v-2A1.5 1.5 0 0 0 9.5 2h-8A1.5 1.5 0 0 0 0 3.5v9A1.5 1.5 0 0 0 1.5 14h8a1.5 1.5 0 0 0 1.5-1.5v-2a.5.5 0 0 0-1 0z"/>
						<path fill-rule="evenodd" d="M15.854 8.354a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L14.293 7.5H5.5a.5.5 0 0 0 0 1h8.793l-2.147 2.146a.5.5 0 0 0 .708.708z"/>
					</svg>
					<span class="link-text">Logout</span> 
					</button>
				</li>
				</ul>
			</div>
		</nav>
	</div>
</template>


<style>
    .nav-item:hover {
        background-color: #f0f0f0; /* Cambia il colore di sfondo al passaggio del mouse */
    }
  
    .nav-item:hover .nav-link {
        font-weight: bold;
    }
    
    /* Regole per i link nella sidebar */
    .large-link {
        font-weight: bold; /* Testo in grassetto */
        font-size: 18px; /* Dimensione del testo più grande */
        margin-bottom: 10px; /* Spazio tra i link */
		border: none;
		background-color: transparent;
    }
    
    .large-link .link-text {
        margin-left: 10px; /* Spazio tra l'icona e il testo */
    }


</style>