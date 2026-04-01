<template>
  <div>
    <Line
      :data="chartData"
      :options="chartOptions"
    />
  </div>
</template>

<script>
import {
  Chart as ChartJS,
  Title, Tooltip, Legend,
  LineElement, CategoryScale, LinearScale, PointElement
} from 'chart.js'
import { Line } from 'vue-chartjs'

ChartJS.register(
  Title, Tooltip, Legend,
  LineElement, PointElement,
  CategoryScale, LinearScale
)

export default {
  components: { Line },

  data() {
    return {
      connection: null,
      chartData: {
        labels: [],
        datasets: [
          {
            label: 'Temperature',
            backgroundColor: '#42A5F5',
            borderColor: '#42A5F5',
            data: []
          }
        ]
      },
      chartOptions: {
        responsive: true,
        animation: false, // ✅ important for real-time feel
        plugins: {
          legend: {
            position: 'top'
          },
          title: {
            display: true,
            text: 'Live Temperature'
          }
        }
      }
    }
  },

  mounted() {
  this.connection = new WebSocket("ws://127.0.0.1:8000/ws")

  this.connection.onopen = () => {
    console.log("WebSocket connected")
  }
this.connection.onmessage = (event) => {
  const payload = JSON.parse(event.data).data
  console.log("Payload:", payload)
  console.log("Temp:", payload.temp, typeof payload.temp)
  const newLabels = [
    ...this.chartData.labels,
    new Date().toLocaleTimeString()
  ]

  const newData = [
    ...this.chartData.datasets[0].data,
    payload.temp
  ]

  // limit size
  if (newLabels.length > 20) {
    newLabels.shift()
    newData.shift()
  }

  // ✅ completely new object (NO mutation)
  this.chartData = {
    labels: newLabels,
    datasets: [
      {
        ...this.chartData.datasets[0],
        data: newData
      }
    ]
  }
}  }
}
</script>
