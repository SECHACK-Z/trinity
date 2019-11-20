<template>
  <pre>{{ result }}</pre>
</template>

<script lang="ts">
import Vue from "vue";
import axios from "axios";
export default Vue.component("Alp", {
  data() {
    return {
      result: "Hello"
    };
  },
  mounted() {
    axios.get("/api/alp").then(response => {
      // console.log(response.data);
      this.result = response.data;
    });
  }
});
</script>
