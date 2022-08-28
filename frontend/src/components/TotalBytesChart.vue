<script setup>
import { Bar } from 'vue-chartjs'

import { inject, ref, computed } from 'vue'

const biz = inject('info')

console.log(biz.value)
const chartData = computed(() => {
    return {
        labels: biz.value.map((v, i) => `Channel ${i}`),
        datasets: [{
            label: "Amount of bytes (B) transmitted",
            data: biz.value.map(i => i.Total),
        }],
    }
})

const chartOptions = ref({
    responsive: true,
    scales: {
        y: {
            beginAtZero: true,
        }
    },
    backgroundColor: [
        'rgba(50, 168, 82,0.3)',
        'rgba(0, 41, 247,0.3)'
    ]
})
</script>

<template>
    <Bar :chart-data="chartData" :chart-options="chartOptions" chart-id="bar-chart" />
</template>