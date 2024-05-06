## 第三次作业
实现一个基于MOM消息队列的分布式信息采集和分析系统

1. 多个数据采集点，每隔100ms采集一次数据（模拟实现），作为消息发布
2. 一个数据分析微服务，对不同设备执行计算：
    - 过去N个数据点的均值和方差
    - 历史数据最大值和最小值
    - 定时将分析结果打包发布为MOM消息
3. 数据显示微服务
    - 针对不同设备，实时绘制采集数据的折线统计图
    - 针对不同设备，实时显示数据的统计分析结果

## 实现方法
1. 采集点：使用多线程模拟多个数据采集点，每个采集点每隔100ms采集一次数据，作为消息发布
2. 数据分析微服务：使用多线程模拟一个数据分析微服务，对不同设备执行计算，定时将分析结果打包发布为MOM消息
3. 数据显示微服务：前端使用Vue3作为图形界面，使用WebSocket实时接收数据，使用ECharts绘制折线统计图，显示数据的统计分析结果；后端使用Go实现WebSocket服务器，接收MOM消息并转发给前端
