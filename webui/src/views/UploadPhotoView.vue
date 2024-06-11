<script>
import { RouterLink } from 'vue-router'

export default {
    data() {
        return {
            errormsg: null,
            successmsg: null,
            loading: false,
            some_data: null,
            image: null,
            previewImage: null
        }
    },
    methods: {

        // Handle the selected file upload.
        handleFileUpload() {
            // Retrieve the first file selected by the user from the file input reference.
            this.image = this.$refs.file.files[0];

            // Create a URL for the selected file to show a preview
            this.previewImage = URL.createObjectURL(this.image);
            
            this.errormsg = null;
            this.successmsg = null;
        },

        async uploadPhoto() {
            if (!this.image) {
                this.errormsg = "Please select a photo to upload";
                return;
            }
            try {
                this.loading = true;
                const userID = localStorage.getItem('userID');
                let response = await this.$axios.post("/users/" + userID + "/photos" , this.image, {
						headers: {
							Authorization: "Bearer " + userID
						}
					});
                this.successmsg = "Photo uploaded successfully!"

            } catch (error) {
                console.error('Error uploading file:', error);
                this.errormsg = "Unsupported image format. Please upload a PNG or JPG file";
            } finally {
                this.loading = false;
                this.previewImage = null;
            }
        }
    },
};
</script>


<template>
    <div>
        <!-- Header -->
        <header class="navbar sticky-top flex-md-nowrap p-0 header">
            <h1 class="app-title unselectable">WASAPhoto</h1>
        </header>

        <!-- Upload Photo area -->
        <main class="container-fluid">
            <div class="row">
                <div class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
                    <div class="mt-3">
                        <h1 class="h2 unselectable border-bottom pb-3" style="font-weight: bold; white-space: nowrap;">
                            Upload Photo</h1>

                        <input type="file" @change="handleFileUpload" accept="image/png, image/jpeg" class="form-control mt-4 mb-3" ref="file">
                        <div v-if="errormsg" class="text-danger">{{ errormsg }}</div>
                        <div v-if="successmsg" class="text-success ml-1">{{ successmsg }}</div>
                        

                        <!-- Preview image -->
                        <h5 v-if="previewImage">Preview:</h5>
                        <div v-if="previewImage" class="preview-photo-card">
                            <img :src="previewImage" alt="Preview" style="width: 100%; height: 100%; object-fit: cover;">
                        </div>

                        <!-- Upload photo button -->
                        <button v-if="previewImage" @click="uploadPhoto" class="btn btn-primary mt-3">Upload</button>
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

<style>
.header {
    background-image: none;
    background-color: #fff;
	box-shadow: 0 3px 5px rgba(0,0,0,0.2);
    height: 70px;
}


.app-title {
    color: #333333;
    font-weight: bold;
    font-size: 26px;
    margin-left: 17px;
    margin-top: 7px;
}

.preview-photo-card {
	width: 380px;
    height: 380px;
	margin: 10px;
	background: #fff;
	border-radius: 15px;
	box-shadow: 0 6px 12px rgba(0, 0, 0, 0.3);
	overflow: hidden; /* Keeps the image within the borders */
}

.unselectable {
    user-select: none;
}
</style>