<template id="mytmp">
  <v-container>
    <Network :nodes="nodes" :edges="edges" :options="options" />
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import { Network } from "vue2vis";
import axios from "axios";

type logData = {
  method: string;
  uri: string;
};
type resultType = {
  uri: string;
  count: number;
};
type nodeType = {
  id: number;
  value: number; //nodeの大きさ
  label: string;
  x: number;
  y: number;
};
type edgeType = {
  from: number;
  to: number;
  value: number;
};
export default Vue.component("Log", {
  template: "#mytmp",
  components: {
    Network
  },
  data() {
    return {
      network: null,
      nodes: [{ id: 0, value: 5, label: "", x: 100, y: 200 }],
      edges: [] as edgeType[],
      options: {
        width: "800px",
        height: "800px",
        nodes: {
          shape: "dot",
          scaling: {
            label: {
              min: 8,
              max: 20
            }
          }
        },
        edges: {
          smooth: false
        },
        physics: false,
        interaction: {
          dragNodes: false, // do not allow dragging nodes
          zoomView: false, // do not allow zooming
          dragView: false // do not allow dragging
        }
      },
      container: null
    };
  },
  mounted() {
    const getRandomInt = (max: number) => {
      return Math.floor(Math.random() * Math.floor(max));
    };
    axios.get("/api/log").then(response => {
      // console.log(response.data);
      const data = JSON.parse(response.data);
      const results: Array<resultType> = [];
      data.forEach((d: any) => {
        const result = results.find(r => r.uri === d.uri);
        if (result) {
          result.count++; // count
        } else {
          results.push({
            uri: d.uri,
            count: 1
          });
        }
      });

      // console.log(results);

      const edges: edgeType[] = [];
      const nodes: nodeType[] = [
        { id: 0, label: " ", value: 5, x: -100, y: 240 }
      ];
      results.map((r, i) => {
        nodes.push({ id: i + 1, label: r.uri, value: 5, x: 150, y: 80 * i });
        edges.push({ from: i + 1, to: 0, value: r.count });
      });
      this.nodes = nodes;
      this.edges = edges;
    });

    // let container = document.getElementById('mynetwork');
    // var data = {
    //     nodes: this.nodes,
    //     edges: this.edges
    // };
    // console.log(Network)
    // let network = Network(container, data, this.options);
  }
});
</script>

<style scoped>
#mynetwork {
  width: 600px;
  height: 600px;
  border: 1px solid lightgray;
}
</style>
