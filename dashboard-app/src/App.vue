<template>
  <div class="container">
    <h1>My dashboard</h1>
    <TodayReadings />
    <LineChart />
  </div>
</template>

<script>
import LineChart from './components/LineChart.vue'
import TodayReadings from './components/TodayReadings.vue'
export default {
  components: {
    LineChart,
    TodayReadings
  }
}
</script>
