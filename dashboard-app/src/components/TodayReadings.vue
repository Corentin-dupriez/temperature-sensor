<template>
  <div>
    <button @click="getToday">Get Today's readings</button>
  </div>
</template>
<script>
export default {
  methods: {
    getToday() {
      fetch(`http://localhost:8000/today`, {
        method: "GET",
      })
      .then(response => response.json())
      .then(data => console.log(data))
    }
  },
  mounted(){}
}
</script>
