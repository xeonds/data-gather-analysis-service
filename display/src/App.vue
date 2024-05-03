<template>
    <el-main>
        <div>
            <el-button @click="id++">+</el-button>
            <span>{{ id }}</span>
            <el-button @click="id--">-</el-button>
        </div>
        <div>
            <el-input :model="url" placeholder="ws://localhost:8001/ws" />
            <el-button @click="conn()">conn</el-button>
        </div>

        <p>{{ content }}</p>
    </el-main>
</template>

<script lang="ts" setup>
import { Ref, onMounted, ref } from 'vue'

const content = ref('')
const id = ref(0)
const ws: Ref<WebSocket> = ref({} as WebSocket)
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
const conn = () => {
    ws.value = initWebSocket(url.value, (data: string) =>
        content.value += data
    )
}

onMounted(() => conn());
</script>
