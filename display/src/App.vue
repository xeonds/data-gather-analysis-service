<template>
    <div>
        <div>
            <label for="device-select">选择设备：</label>
            <select id="device-select" v-model="selectedDevice" @change="switchDisplayDevice">
                <option v-for="device in deviceList" :key="device.id" :value="device.id">{{ device.name }}</option>
            </select>
        </div>
        <div>
            <div id="chart" style="width: 600px; height: 400px;"></div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, Ref } from 'vue';
import * as echarts from 'echarts';

const selectedDevice: Ref<any> = ref(0);
const deviceList: Ref<any[]> = ref([]);
const analysisData: Ref<any> = ref({});
const chart: Ref<any> = ref(null);
const ws: Ref<WebSocket> = ref({} as WebSocket)
const url = ref('ws://localhost:8001/ws')
const data: Ref<any> = ref({})

const fetchDeviceList = () => {
    fetch('/count')
        .then(response => response.json())
        .then(data => deviceList.value = Array.from({ length: data.count }, (_, i) => ({ id: i, name: `设备${i}` })))
        .catch(error => console.error('Error fetching device list:', error));
};

const switchDisplayDevice = () => {
    analysisData.value = data.value[selectedDevice.value] || {};
};

watch(selectedDevice, () => switchDisplayDevice());
watch(analysisData, () => {
    if (analysisData.value) {
        chart.value.setOption({
            title: { text: '数据分析图表' },
            xAxis: {
                type: 'category',
                data: ['Max', 'Min', 'Avg', 'Variance']
            },
            yAxis: { type: 'value' },
            series: [{
                data: [analysisData.value.Max, analysisData.value.Min, analysisData.value.Avg, analysisData.value.Variance],
                type: 'bar'
            }]
        });
    }
}, { deep: true });

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
    ws.value = initWebSocket(
        url.value,
        (e: string) => {
            const parsed = JSON.parse(e)
            data.value[parsed.id] = parsed.data
            switchDisplayDevice()
        })
}

onMounted(() => {
    fetchDeviceList();
    chart.value = echarts.init(document.getElementById('chart') as HTMLDivElement);
    conn();
});
</script>
