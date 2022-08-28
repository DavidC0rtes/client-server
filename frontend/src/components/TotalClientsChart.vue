<script setup>
import { Line } from 'vue-chartjs'
import { inject, ref, watch } from 'vue'

const biz = inject('info')
const labels = ref([])
const timelyData = []

watch(biz, (newData) => {
	labels.value.push(new Date().toLocaleTimeString())
	timelyData.push( getTotalClients())
})

const getTotalClients = () => {
	return biz.value.reduce((accum, object) => {
		return accum+Object.keys(object.Clients).length
	},0)
}

const chartData = {
			labels: labels.value,
			datasets: [{
				label: "Concurrent clients",
				data: timelyData,
			}],
}

const chartOptions = ref({
	responsive: true,
	fill: false,
	backgroundColor: [
		'rgba(68, 10, 131, 0.5)',
	],
	scales: {
		y: {
			ticks: {
				stepSize: 1
			}
		}
	}
})

</script>

<template>
<Line
	:chart-data="chartData"
	:chart-options="chartOptions"
	chart-id="line-chart"
/>
</template>