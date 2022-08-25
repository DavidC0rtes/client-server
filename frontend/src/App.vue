<script setup>
import ClientData from './components/ClientData.vue';
</script>

<template>
  <main>
    <div id="root" class="container">
      <div class="row">
        <div class="col-md-6 offset-md-3 py-5">
          <h1>Server frontend</h1>
          <ClientData/>
        </div>
      </div>
    </div>
  </main>
</template>