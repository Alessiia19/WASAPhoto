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
      if (!usernameRegex.test(this.username)) {
        this.errormsg = "L'username deve contenere solo lettere e numeri.";
        return false;
      } else if (this.username.trim().length < 3 || this.username.trim().length > 16) {
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
      <h1 class="app-title unselectable">WASAPhoto</h1>
      <div class="underline ms-2"></div>
    </div>
    <div class="login-wrapper">
      <div class="login-box">
        <h2 class="login-title">Login</h2>
        <form @submit.prevent="login">
          <div class="form-group">
            <label for="username" class="form-label">Username</label>
            <input type="text" class="form-control" id="username" v-model="username" required autocomplete="off"
              @input="handleUsernameInput">
            <div v-if="errormsg" class="text-danger">{{ errormsg }}</div>
          </div>
          <div class="text-center">
            <button type="submit" class="btn btn-primary" :disabled="errormsg">Login</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>


<style scoped>
.login-wrapper {
  position: relative;
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
  top: 160px;
  left: 680px;
}

.app-title {
  color: #00264d;
  font-weight: bold;
  font-size: 50px;
}

.underline {
  width: 300px;
  height: 2px;
  background-color: #00264d;
  margin-top: 5px;
}

.login-box {
  background-color: #ffffff;
  padding: 60px;
  border-radius: 15px;
  box-shadow: 0 0 15px rgba(0, 0, 0, 0.3);
  height: 320px;
  width: 400px;
}

.login-title {
  color: #00264d;
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
}

.form-control {
  background-color: #385273;
  color: #ffffff;
}

.btn-primary {
  background-color: #00264d;
  color: #ffffff;
  border: none;
  padding: 10px 20px;
  border-radius: 5px;
  cursor: pointer;
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
}

.btn-primary:hover {
  background-color: #001a33;
}

.unselectable {
	user-select: none;
}
</style>
