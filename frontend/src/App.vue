<script setup>
import ClientData from './components/ClientData.vue';
import { provide, ref, onMounted } from 'vue'

const info = ref([])
provide('info', info)

const showError = ref(false)

const url = "http://localhost:8080/info"

onMounted(async () => {
  updateData()
})

const toggleRotate = () => {
  const icon = document.querySelector("#icon-refresh")
  icon.classList.add("rotate")
  
  setTimeout(() => {
    icon.classList.remove("rotate")
  }, 1000)
}

const updateData = async () => {
  try {
    info.value = await (await fetch(url)).json()
    console.info(info.value)
    showError.value = false
  } catch (error) {
    console.error(error)
    showError.value = true
  }
  

}
</script>

<template>
  <nav class="navbar navbar-dark sticky-top bg-dark">
    <div class="container-fluid">
      <h1 id="title">
        <img
          src="https://cdn-icons-png.flaticon.com/512/7403/7403579.png"
          alt
          width="40"
          height="40"
          class="d-inline-block align-text-top"
        />
        FileShare
      </h1>
       <button
          type="button"
          id="refresh-button"
          class="btn btn-primary btn-md col-md-2"
          @click="updateData();toggleRotate()"
        ><i id="icon-refresh" class="fa-solid fa-arrow-rotate-right"></i> Refresh</button>
    </div>
  </nav>
  <main>
    <div id="root" class="container">
      <div id="header" class="row">
        <h2 class="col-md-5">Latest data from server</h2>
        
      </div>
      <hr />
      <div id="second-header" class="row">
        <div class="col-md-12">
          <ClientData />
        </div>
      </div>
    </div>
  </main>
  <div
    v-if="showError"
    id="error-toast"
    class="toast align-items-center text-white bg-danger border-0 show"
    role="alert"
    aria-live="assertive"
    aria-atomic="true"
  >
    <div class="d-flex">
      <div class="toast-body">Failed to fetch from server.</div>
      <button
        type="button"
        class="btn-close me-2 m-auto"
        data-bs-dismiss="toast"
        aria-label="Close"
      ></button>
    </div>
  </div>
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
  width: 100%;
}
hr {
  margin-top: 0;
}
#error-toast {
  z-index: 11;
  inset-block-end: 0;
  top: 50;
  inset-inline-end: 0;
  position: absolute;
}

.rotate {
  -moz-transition: all 0.5s linear;
  -webkit-transition: all 0.5s linear;
  transition: all 0.5s linear;
  -moz-transform:rotate(360deg);
  -webkit-transform:rotate(360deg);
  transform:rotate(360deg);
}
</style>