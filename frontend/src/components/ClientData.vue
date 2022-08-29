<!-- filename: ClientData.vue -->
<script setup>
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    LineElement,
    PointElement,
    CategoryScale,
    LinearScale,
} from 'chart.js'


import TotalBytesChart from './TotalBytesChart.vue'
import TotalClientsChart from './TotalClientsChart.vue'
import FiletypesChart from './FiletypesChart.vue'
import ClientsDropdown from './ClientsDropdown.vue'
import { inject } from 'vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, LineElement, PointElement)

const info = inject('info')
</script>

<template>
    <div id="datarow" class="row">
        <h4>Listening clients📨</h4>
        <h5 v-if="info[0]">Maximum size: <span class="badge bg-info text-dark">{{info[0].MaxSize}}B</span></h5>
        <div class="col-md-7">
            <ClientsDropdown />
        </div>
        <img
            class="col-md-5"
            id="computer-img"
            src="https://images.unsplash.com/photo-1550745165-9bc0b252726f?ixlib=rb-1.2.1&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=370&q=80"
        />
    </div>

    <div class="row">
        <h4>Statistics📊</h4>
        <hr />
        <div class="col-md-4">
            <TotalBytesChart />
        </div>
        <div class="col-md-4">
            <TotalClientsChart />
        </div>
        <div class="col-md-4">
            <FiletypesChart />
        </div>
    </div>
</template>
<style>
#refresh-button {
    margin: 1em 0;
}
#computer-img {
    object-fit: cover;
    border-radius: 2em;
    height: 200px;
}
#datarow {
    align-items: center;
}
h4 {
    font-weight: 600;
}
</style>