
<template id="mytmp">
  <div>
    <Network id="mynetwork"
    :nodes="nodes"
    :edges="edges"
    :options="options" />
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import { Network } from 'vue2vis';

export default Vue.component('Log',{
    template: "#mytmp",
    components: {
        Network
    },
    data(){
        return {
        network: null,
        nodes: [
            { id: 1, value: 5, label: "" },
            { id: 2, value: 5, label: "/Service/1" },
            { id: 3, value: 5, label: "/Service/2" },
            { id: 4, value: 5, label: "/Service/3" },
            { id: 5, value: 5, label: "/Service/4" },
        ],
        edges: [
            { from: 1, to: 2, value: 7, title: "" },
            { from: 1, to: 3, value: 10, title: "" },
            { from: 1, to: 4, value: 8, title: "" },
            { from: 1, to: 5, value: 9, title: "" },
        ],
        options: {
            nodes: {
            shape: "dot",
            scaling: {
                label: {
                min: 8,
                max: 20
                }
            }
            }
        },
        container: null
        }
    },
    mounted(){
        // let container = document.getElementById('mynetwork');
        // var data = {
        //     nodes: this.nodes,
        //     edges: this.edges
        // };
        // console.log(Network)
        // let network = Network(container, data, this.options);
    }
})
</script>

<style scoped>
#mynetwork {
  width: 600px;
  height: 600px;
  border: 1px solid lightgray;
}
</style>