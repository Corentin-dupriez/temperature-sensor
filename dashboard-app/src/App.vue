<template>
  <div class="container">
    <h1>My dashboard</h1>
    <TodayReadings />
    <WeekReadings />
    <LineChart />
  </div>
</template>

<script>
import LineChart from './components/LineChart.vue'
import TodayReadings from './components/TodayReadings.vue'
import WeekReadings from './components/WeekReadings.vue'
export default {
  components: {
    LineChart,
    TodayReadings,
    WeekReadings
  }
}
</script>
