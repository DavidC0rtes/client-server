<script setup>
import { Bar } from 'vue-chartjs'
import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale,
} from 'chart.js'
import { inject, ref, computed } from 'vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

const biz = inject('data')

console.log(biz)

//let bytesTransmitted = ref(Object.keys(biz.value).map(i => biz.value[i].Total))

const chartData = computed(() => {
    return {
        labels: Object.keys(biz.value),
        datasets: [{
            label: "Amount of bytes (B) transmitted",
            data: Object.keys(biz.value).map(i => biz.value[i].Total),
        }],
    }
})

const chartOptions = ref({
    responsive: true,
    scales: {
        y: { beginAtZero: true }
    },
})
</script>

<template>
    <Bar :chart-data="chartData" :chart-options="chartOptions" />
</template>