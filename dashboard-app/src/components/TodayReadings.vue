<template>
  <div>
    <p v-if="data">{{data}}</p>
    <p v-else>No readings for today</p>
  </div>
</template>
<script>
export default {
  data () {
    return {
      data : null
    }
  },
  mounted() {

      fetch(`http://localhost:8000/today`, {
        method: "GET",
      })
      .then(response => response.json())
      .then(data => {
        this.data = data
      })
      .catch(err => {
      console.log(err)
      }
      )
    }
}
</script>
