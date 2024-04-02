<template>
	<div>
	  <!-- Header -->
	  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow header">
		<h1 class="app-title unselectable">WASAPhoto</h1>
	  </header>
	  
      <!-- Upload Photo area -->
      <main class="container-fluid">
          <div class="row">
              <div class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                  <!-- Aggiungi un'area per il caricamento della foto -->
                  <div class="mt-3">
                      <h1 class="h2 unselectable border-bottom pb-3" style="font-weight: bold; white-space: nowrap;">Upload Photo</h1>
                      
                      <input type="file" @change="handleFileUpload" accept="image/*" class="form-control mt-4 mb-3">
                      
                      <!-- Aggiungi l'anteprima dell'immagine -->
                      <div v-if="previewImage" class="mt-3">
                        <h5>Preview:</h5>
                        <img :src="previewImage" alt="Preview" style="max-width: 100%;">
                    </div>
                    
                    <!-- Aggiungi il pulsante di caricamento -->
                    <button @click="uploadPhoto" class="btn btn-primary mt-3">Upload</button>
                </div>
            </div>
        </div>
      </main>
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
            selectedFile: null,
            previewImage: null
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
        // Metodo per gestire il caricamento del file selezionato
        handleFileUpload(event) {
            this.selectedFile = event.target.files[0];
            this.previewImage = URL.createObjectURL(this.selectedFile);
        },
        // Metodo per effettuare il caricamento effettivo del file
        async uploadPhoto() {
            if (!this.selectedFile) {
                alert("Please select a file to upload.");
                return;
            }
            try {
                this.loading = true;
                const formData = new FormData();
                formData.append('photo', this.selectedFile);
                // Aggiungi qui la logica per l'invio del formData al server
                // Utilizza ad esempio Axios o un metodo nativo di JavaScript per inviare la richiesta al server
                // Esempio con Axios:
                // await this.$axios.post('/upload-photo', formData);
                alert("File uploaded successfully.");
            } catch (error) {
                console.error('Error uploading file:', error);
                alert("An error occurred while uploading the file.");
            } finally {
                this.loading = false;
                // Resetta il valore dell'input del file e l'anteprima dell'immagine
                this.selectedFile = null;
                this.previewImage = null;
            }
        }
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
  
  .unselectable {
    user-select: none;
  }

  </style>
  