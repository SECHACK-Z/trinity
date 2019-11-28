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
    axios.get('/api/alp').then(response => {
      this.lineChartData = {
        expectedData: [100, 120, 161, 134, 105, 160, 165],
        actualData: [120, 82, 91, 154, 162, 140, 145]
      }
      this.pieChartData = {
        hosts: ['a.sechack-z.org', 'b.sechack-z.org'],
        pieChartData: [
          { value: 320, name: 'a.sechack-z.org' },
          { value: 240, name: 'b.sechack-z.org' }
        ]
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
