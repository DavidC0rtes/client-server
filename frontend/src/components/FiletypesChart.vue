<script setup>
import { watch,inject, ref, computed, } from 'vue'
import { Bar } from 'vue-chartjs'

const info = inject('info')
const types = ref({})
const files = ref(new Set())

// Generate random RGBA string, alpha max value is 0.6 for aesthetics and max for colors
// is 248 to avoid hard to see colors.
const randomRGBA = () => {
    let o = Math.round, r = Math.random, s = 248
    return 'rgba(' + o(r()*s)+ ',' + o(r()*s) + ',' + o(r()*s) + ',' + r().toFixed(1)*0.6 + ')'
}


watch(info, (newInfo) => {	
	newInfo.forEach(element => {
		if (element.CurrFile != "" && !files.value.has(element.CurrFile)) {
			const ext = element.CurrFile.split('.').pop() != element.CurrFile
						? element.CurrFile.split('.').pop()
						: "other"
			if (ext in types.value) {
				types.value[ext] += 1
			} else {
				types.value[ext] = 1
			}
			files.value.add(element.CurrFile)
		}
		console.log(types.value)
		console.log(files.value)
	});
})

const chartData = computed(() => {
	return {
		labels: Object.keys(types.value),
		datasets: [{
			label: "Filetypes transmitted",
			data: Object.values(types.value)
		}]
	}
})

const chartOptions = computed(() => {
	return {
	responsive: true,
	indexAxis: 'y',
	scales: {
		x: {
			ticks: {
				stepSize: 1
			}
		}
	},
	backgroundColor: Object.keys(types.value).map(x => randomRGBA())
	}
})
</script>

<template>
<Bar :chart-data="chartData" chart-id="bar-filetypes" :chart-options="chartOptions" />
</template>