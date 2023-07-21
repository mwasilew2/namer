<template>
  <div id="app" class="container">
    <div class="row">
      <div class="col-md-6 offset-md-3 py-5">
        <ul>
          {{ this.year }}
        </ul>
        <ul>
          {{ this.currName }}
        </ul>

        <button class="btn btn-primary" @click="saveAName">Save</button>
        <button class="btn btn-primary" @click="getAName">Next</button>

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
      apiUrl: "http://localhost:8080/api/v1/name",
      year: new Date().getFullYear(),
      lastPageId: -1,
      limit: 10,

      currPage: [],
      currNameId: 0,
      currName: null,

      favourites: [],
    }
  },

  methods: {
    async getAName() {
      if (this.currPage.length == 0 || this.currNameId == this.currPage.length - 1) {
        this.lastPageId += 1
        this.currNameId = 0
        axios
            .get(this.apiUrl, {
              params: {
                year: this.year,
                limit: this.limit,
                page: this.lastPageId,
              },
            })
            .then((response) => {
              this.currPage = response.data.names
              this.currName = this.currPage[this.currNameId].name
            })
            .catch((error) => {
              window.alert(`The api returned an error: ${error}`)
            })
      } else {
        this.currNameId += 1
        this.currName = this.currPage[this.currNameId].name
      }


    },
    async saveAName() {
      this.favourites.push(this.currName)
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
