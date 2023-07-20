<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <ul>
          {{ nameGot }}
        </ul>

        <button class="btn btn-primary" @click="getAName">Get Name!</button>

      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'App',

  data() {
    return {
      nameGot: null
    }
  },

  methods: {
    async getAName() {
      this.nameGot = await axios.get("http://localhost:8082/api/v1/name/0")
          .catch((error) => {
            window.alert(`The api returned an error: ${error}`)
          })
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
