<template>
  <div>
    <p v-if="data && data.length">{{data[0].temp}}C, {{data[0].humidity}}%</p>
    <p v-else>No readings for this week</p>
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

      fetch(`http://localhost:8000/week`, {
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
