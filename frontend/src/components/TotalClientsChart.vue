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
import { inject, ref } from 'vue'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)


const props = defineProps({
	data: Object
})
console.log(inject('data'))
//const bytesTransmitted = Object.keys(props.data).filter(i => props.data[i].Total)

const chartData = ref({
	labels: ['0', '1'],
	datasets: [{
		label: "Connected clients since last update",
		data: [256, 4096]
	}],
})

const chartOptions = ref({
	responsive: true,
	scales: {
		y: { beginAtZero: true }
	},
})

</script>

<template>
<Bar
	:chart-data="chartData"
	:chart-options="chartOptions"
/>
</template>