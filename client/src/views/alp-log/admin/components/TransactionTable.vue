<template>
  <el-table :data="list" style="width: 100%;padding-top: 15px;">
    <el-table-column label="URI" min-width="200">
      <template slot-scope="scope">
        <!-- {{ scope.row.uri | orderNoFilter }} -->
        {{ scope.row.uri }}
      </template>
    </el-table-column>
    <el-table-column label="Method" width="195" align="center">
      <template slot-scope="scope">
        <!-- ¥{{ scope.row.method | toThousandFilter }} -->
        {{ scope.row.method }}
      </template>
    </el-table-column>
    <el-table-column label="Status" width="100" align="center">
      <template slot-scope="{row}">
        <el-tag :type="row.status | statusFilter">
          {{ row.status }}
        </el-tag>
      </template>
    </el-table-column>
  </el-table>
</template>

<script>
// import { transactionList } from '@/api/remote-search'
import axios from 'axios'
export default {
  filters: {
    statusFilter(status) {
      const statusMap = {
        200: 'success',
        404: 'danger'
      }
      return statusMap[status]
    },
    orderNoFilter(str) {
      return str.substring(0, 30)
    }
  },
  props: {
    chartData: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      list: null
    }
  },
  watch: {
    chartData: {
      deep: true,
      handler(val) {
        this.setOptions(val)
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      axios.get('/api/logs').then(response => {
        console.log('accessLog', response)
        this.list = response.data.reverse().slice(0, 8)
      })
      // transactionList().then(response => {
      //   console.log("tran", "response.data)
      //   this.list = response.data.slice(0, 8)
      // })
    }
  }
}
</script>
