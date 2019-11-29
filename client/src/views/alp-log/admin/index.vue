<template>
  <div class="dashboard-editor-container">
    <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <line-chart v-if="isFetched" :chart-data="lineChartData" />
    </el-row>

    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="8">
        <div v-if="isFetched" class="chart-wrapper">
          <pie-chart :chart-data="pieChartData" />
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div v-if="isFetched" class="chart-wrapper">
          <bar-chart :chart-data="barChartData" />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import LineChart from './components/LineChart'
import PieChart from './components/PieChart'
import BarChart from './components/BarChart'
import axios from 'axios'

export default {
  name: 'DashboardAdmin',
  components: {
    LineChart,
    PieChart,
    BarChart
  },
  data() {
    return {
      lineChartData: {},
      pieChartData: {},
      barChartData: {},

      isFetched: false
    }
  },
  mounted() {
    axios.get('/api/rawLog').then(response => {
      console.log(response)
      // const nameOfDay = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat']
      const group = response.data.reduce(function(result, current) {
        const element = result.find(function(p) {
          return p.name === current.host
        })
        if (element) {
          element.value++ // count
        } else {
          result.push({
            name: current.host,
            value: 1
          })
        }
        return result
      }, [])
      const hosts = group.map(l => l.name)
      console.log(hosts)
      console.log(group)
      this.lineChartData = { hosts: ['a.sechack-z.org', 'b.sechack-z.org'],
        expectedData: [79, 52, 200, 334, 390, 330, 220],
        actualData: [30, 100, 150, 450, 250, 100, 120]
      }
      this.pieChartData = {
        hosts: hosts,
        pieChartData: group
      }
      this.barChartData = [
        { host: 'a.sechack-z.org', count: [79, 52, 200, 334, 390, 330, 220] },
        { host: 'b.sechack-z.org', count: [80, 52, 200, 334, 390, 330, 220] }
      ]

      this.isFetched = true
    })
  }
  // methods: {
  //   handleSetLineChartData(type) {
  //     this.lineChartData = lineChartData
  //   }
  // }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width: 1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
