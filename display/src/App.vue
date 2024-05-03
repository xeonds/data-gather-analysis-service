<template>
    <el-main>
        <el-button @click="id++">+</el-button>
        <el-button @click="id--">-</el-button>
        <p>{{ content }}</p>
    </el-main>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'

const content = ref('')
const id = ref(0)
const url = ref('ws://localhost:8001/ws')
const initWebSocket = (addr: string, callback: Function) => {
    const socket = new WebSocket(addr)
    socket.addEventListener('open', () => {
        console.log('connected')
    })
    socket.addEventListener('close', () => {
        console.log('disconnected')
    })
    socket.addEventListener('message', (event) => {
        console.log(event.data);
        callback(event.data);
    })
    return socket
}

onMounted(() =>
    initWebSocket(url.value, (data: string) =>
        content.value += data
    )
);
</script>
