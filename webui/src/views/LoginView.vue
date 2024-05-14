<script>
export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      username: "",
      userToSearchID: 0,
      user: { userID: 0, username: "" }
    };
  },

  methods: {
    async login() {

      // Username validation
      if (!this.validateUsername()) {
        return;
      }

      try {
        let response = await this.$axios.post("/session", { username: this.username.trim() })

        // Dopo aver ricevuto l'user id dal backend setta i valori 
        this.user = response.data
        localStorage.setItem("userID", this.user.userID)
        localStorage.setItem("username", this.user.username)
        localStorage.setItem("userToSearchID", this.userToSearchID)
        console.log("User loggato:", this.user)
        console.log("Localstorage:", localStorage)

        // Redirect alla home dopo il login
        this.$router.replace({ path: '/users/' + this.user.username + '/stream' })

      } catch (error) {
        // Gestione degli errori
        this.errormsg = error.response.data.message
      }
    },

    handleUsernameInput() {
      // Chiamata alla funzione di validazione durante ogni inserimento
      this.validateUsername();
    },

    validateUsername() {
      const usernameRegex = /^[a-zA-Z0-9]+$/;
      if (!usernameRegex.test(this.username.trim())) {
        this.errormsg = "L'username deve contenere solo lettere e numeri.";
        return false;
      } else if (this.username.length < 3 || this.username.length > 16) {
        this.errormsg = "L'username deve essere compreso tra 3 e 16 caratteri.";
        return false;
      } else {
        this.errormsg = null; // Resetta il messaggio di errore se la validazione ha successo
        return true;
      }
    },

    mounted() {
      if (localStorage.getItem('username')) {
        this.$router.replace({ path: '/users/' + localStorage.getItem('username') + '/stream' });
      }
    },

  }
};
</script>

<template>
  <div class="login-container">
    <div class="app-name">
      <h1 class="app-title">WASAPhoto</h1>
      <div class="underline"></div>
    </div>
    <div class="login-wrapper">
      <div class="login-box">
        <h2 class="login-title">Login</h2>
        <form @submit.prevent="login">
          <div class="form-group">
            <label for="username" class="form-label">Username</label>
            <input type="text" class="form-control" id="username" v-model="username" required autocomplete="username"
              @input="handleUsernameInput">
            <div v-if="errormsg" class="text-danger">{{ errormsg }}</div>
          </div>
          <div class="text-center">
            <button type="submit" class="btn btn-primary">Login</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>


<style scoped>
.login-wrapper {
  position: relative;
  /* Posizionamento relativo per il contenitore */
}

.login-container {
  background: linear-gradient(to bottom right, #f5dd90, #b97b90, #446ca0);
  background-size: cover;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.app-name {
  position: absolute;
  /* Posizionamento assoluto */
  top: 160px;
  /* Distanza dal top */
  left: 690px;
  /* Distanza dalla sinistra */
}

.app-title {
  color: #00264d;
  /* Colore del titolo dell'applicazione */
  font-weight: bold;
  /* Rende il titolo in grassetto */
  font-size: 50px;
  /* Dimensione del titolo dell'applicazione */
}

.underline {
  width: 300px;
  /* Lunghezza della riga sottostante */
  height: 2px;
  /* Spessore della riga sottostante */
  background-color: #00264d;
  /* Colore della riga sottostante */
  margin-top: 5px;
  /* Spazio sopra la riga sottostante */
}

.login-box {
  background-color: #ffffff;
  /* Colore del rettangolo del box di login */
  padding: 60px;
  /* Spaziatura interna del box di login */
  border-radius: 15px;
  /* Arrotondamento dei bordi del box di login */
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.3);
  /* Ombra del box di login */
  height: 320px;
  /* Altezza fissa del riquadro di login */
  width: 400px;
  /* Larghezza fissa del riquadro di login */
}

.login-title {
  color: #00264d;
  /* Colore del titolo */
  text-align: center;
  margin-bottom: 40px;
  font-weight: bold;
  font-size: 34px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  color: #00264d;
  /* Colore delle etichette dei campi di input */
}

.form-control {
  background-color: #385273;
  /* Colore di sfondo del campo di input */
  color: #ffffff;
  /* Colore del testo del campo di input */
}

.btn-primary {
  background-color: #00264d;
  /* Colore di sfondo del bottone */
  color: #ffffff;
  /* Colore del testo del bottone */
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
  position: absolute;
  /* Posizionamento assoluto per il pulsante */
  bottom: 20px;
  /* Distanza dal fondo */
  left: 50%;
  /* Allineamento al centro */
  transform: translateX(-50%);
  /* Per centrare orizzontalmente */
}

.btn-primary:hover {
  background-color: #001a33;
  /* Cambio di colore al passaggio del mouse sul bottone */
}
</style>
