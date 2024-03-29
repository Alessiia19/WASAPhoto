<script>
export default {
  data: function() {
    return {
      username: "",
      errormsg: null
    };
  },

  methods: {
    async login() {
      try {
        let response = await this.$axios.post("/session", {
          username: this.username.trim()
        });

        // Dopo aver ricevuto l'user id dal backend
        localStorage.setItem("userID", response.data.userID);
        
        // Redirect alla home dopo il login
        this.$router.replace("/home");
        this.$emit('changeLoginStatus', true)

      } catch (error) {
        // Gestione degli errori
        this.errormsg = error.response.data.message;
      }
    },

    mounted(){
      if (localStorage.getItem('userID')){
        this.$router.replace("/home")
      }
    },

  }
};
</script>

<template>
  <div class="login-container">
    <div class="app-name">
      <h1 class="app-title">WASAPhoto</h1>
      <div class="underline"></div> <!-- Riga sottostante -->
    </div>
    <div class="login-box">
      <h2 class="login-title">Login</h2>
      <form @submit.prevent="login">
        <div class="form-group">
          <label for="username" class="form-label">Username</label>
          <input type="text" class="form-control" id="username" v-model="username" required>
        </div>
        <div class="text-center">
          <button type="submit" class="btn btn-primary">Login</button>
        </div>
      </form>
    </div>
  </div>
</template>


<style scoped>
.login-container {
  background: linear-gradient(to bottom right, #f5dd90, #b97b90, #446ca0);
  background-size: cover;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.app-name {
  position: absolute; /* Posizionamento assoluto */
  top: 160px; /* Distanza dal top */
  left: 690px; /* Distanza dalla sinistra */
}

.app-title {
  color: #00264d; /* Colore del titolo dell'applicazione */
  font-weight: bold; /* Rende il titolo in grassetto */
  font-size: 50px; /* Dimensione del titolo dell'applicazione */
}

.underline {
  width: 300px; /* Lunghezza della riga sottostante */
  height: 2px; /* Spessore della riga sottostante */
  background-color: #00264d; /* Colore della riga sottostante */
  margin-top: 5px; /* Spazio sopra la riga sottostante */
}

.login-box {
  background-color: #ffffff; /* Colore del rettangolo del box di login */
  padding: 60px;  /* Spaziatura interna del box di login */
  border-radius: 15px; /* Arrotondamento dei bordi del box di login */
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.3);  /* Ombra del box di login */
}

.login-title {
  color: #00264d; /* Colore del titolo */
  text-align: center;
  margin-bottom: 40px;
  font-weight: bold;
  font-size: 34px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  color: #00264d; /* Colore delle etichette dei campi di input */
}

.form-control {
  background-color: #385273; /* Colore di sfondo del campo di input */
  color: #ffffff; /* Colore del testo del campo di input */
}

.btn-primary {
  background-color: #00264d; /* Colore di sfondo del bottone */
  color: #ffffff; /* Colore del testo del bottone */
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
}

.btn-primary:hover {
  background-color: #001a33; /* Cambio di colore al passaggio del mouse sul bottone */
}
</style>
