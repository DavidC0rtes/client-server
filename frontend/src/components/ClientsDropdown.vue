<script setup>
import { inject } from 'vue';
const info = inject('info')
</script>

<template>
<details v-for="(channelObj, i) in info" :key="i">
<summary>{{channelObj.CurrFile}}({{channelObj.Filesize}}B) @ channel {{i}}</summary>
    <ul class="list-group list-group-horizontal-sm">
        <li v-for="(client, k) in channelObj.Clients" :key="k" class="list-group-item">{{client}}</li>
    </ul>
</details>
</template>