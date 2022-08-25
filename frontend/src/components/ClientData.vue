<!-- filename: ClientData.vue -->
<script setup>
import {ref, onMounted} from 'vue'

const data = ref(0)

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
		<div v-for="(channel, i) in data" :key="i">On channel {{i}}:
			<ul v-for="(client, k) in channel.Clients.filter(x => x != '')" :key="k">
				<li>{{client}}</li>
			</ul>
		</div>

		<h4>Statistics</h4>

		<div class="col-md-4">
			<button type="button" class="btn btn-primary" @click="updateData">
				Refresh
			</button>
		</div>
	</div>
</template>