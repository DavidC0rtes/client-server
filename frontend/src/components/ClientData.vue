<!-- filename: ClientData.vue -->
<script setup>
import {ref, onMounted, provide} from 'vue'
import TotalBytesChart from './TotalBytesChart.vue'

const data = ref({0: 'a', 1:'b'})
provide('data', data)
const url = "http://localhost:8080/info"

onMounted(async () => {
    updateData()
})

const updateData = async() => {
    data.value = await ( await fetch(url)).json()
    console.log(data.value)
}
</script>

<template>
    <div id="datarow" class="row">
        <h2>Latest data from server</h2>
        <hr/>
        <h4>Listening clients</h4>
        <div v-for="(channelObj, i) in data" :key="i">On channel {{i}}:
            <ul v-for="(client, k) in channelObj.Clients" :key="k">
                <li>{{client}}</li>
            </ul>
        </div>
    </div>
    <br>
    <div class="row">
        <h4>Statistics</h4>
        <div class="col-md-12">
            <div class="col-md-6">
                <TotalBytesChart :data="data"/>
            </div>
            
        </div>

        <div class="col-md-7">
            <button type="button" class="btn btn-primary" @click="updateData">
                Refresh
            </button>
        </div>
    </div>
</template>