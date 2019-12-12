<template>
  <div v-if="isFetched" class="dashboard-editor-container">
    <el-row style="background:#fff;padding:16px 16px 0;margin-bottom:32px;">
      <line-chart v-if="isFetched" :chart-data="chartData" />
    </el-row>

    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="8">
        <div v-if="isFetched" class="chart-wrapper">
          <pie-chart :chart-data="chartData" />
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div v-if="isFetched" class="chart-wrapper">
          <bar-chart :chart-data="chartData" />
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="8">
      <el-col
        :xs="{ span: 24 }"
        :sm="{ span: 24 }"
        :md="{ span: 24 }"
        :lg="{ span: 12 }"
        :xl="{ span: 12 }"
        style="padding-right:8px;margin-bottom:30px;"
      >
        <transaction-table :chart-data="accessData" />
      </el-col>
    </el-row>
  </div>
</template>

<script>
import LineChart from './components/LineChart'
import PieChart from './components/PieChart'
import BarChart from './components/BarChart'
import TransactionTable from './components/TransactionTable'
import axios from 'axios'

export default {
  name: 'DashboardAdmin',
  components: {
    LineChart,
    PieChart,
    BarChart,
    TransactionTable
  },
  data() {
    return {
      chartData: {},
      accessData: {},
      isFetched: false
    }
  },
  mounted() {
    axios.get('/api/logs').then(response => {
      console.log(response)
      const data = response.data
      const aggregateData = data.reduce(function(result, current) {
        const element = result.find(function(p) {
          return p.name === current.host
        })
        const date = new Date(current.time)

        if (element) {
          element.count[date.getDay()]++
        } else {
          const tmpData = {
            name: current.host,
            count: [0, 0, 0, 0, 0, 0, 0]
          }
          tmpData.count[date.getDay()]++
          result.push(tmpData)
        }
        return result
      }, [])

      const group = aggregateData.map(l => {
        const sum = l.count.reduce((r, c) => r + c)
        l['value'] = sum
        // 月曜始まりにする
        const head = l.count.shift()
        l.count.push(head)
        return l
      })

      this.chartData = group
      this.accessData = data
      this.isFetched = true
    })
  }
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
