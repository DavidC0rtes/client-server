<script setup>
import ClientData from './components/ClientData.vue';
import { provide, ref, onMounted } from 'vue' 

const info = ref([])
provide('info', info)
const url = "http://localhost:8080/info"

onMounted(async () => {
    updateData()
})

const updateData = async() => {
    info.value = await ( await fetch(url)).json()
    console.log(info.value)
}
</script>

<template>
  <nav class="navbar navbar-dark sticky-top bg-dark">
    <div class="container-fluid">
      <h1 id="title">
        <img src="https://cdn-icons-png.flaticon.com/512/7403/7403579.png" alt="" width="40" height="40" class="d-inline-block align-text-top">
      Server frontend
      </h1>
    </div>
  </nav>
  <main>
    <div id="root" class="container">
      <div id="header" class="row">
        <h2 class="col-md-5">Latest data from server</h2>
        <button type="button" id="refresh-button" class="btn btn-primary btn-sm col-md-1" @click="updateData">
            Refresh
        </button>
    </div>
    <hr/>
      <div id="second-header" class="row">
        <div class="col-md-12">
          <ClientData/>
        </div>
      </div>
    </div>
  </main>
</template>

<style>
body {
  font-family: "Poppins";
}
#title {
  font-weight: 700;
  color: white;
  padding-left: 2.5em;
}

#header {
    align-items: center;
}

#second-header {
  align-items: center;
  height: 100%;
  width: 1
  00%;
}
hr {
  margin-top: 0;
}
</style>