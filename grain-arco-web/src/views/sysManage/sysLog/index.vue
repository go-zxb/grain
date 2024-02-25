<template>
  <div class="container">
    <a-drawer
      :width="740"
      :height="340"
      :visible="visible"
      :footer="false"
      unmount-on-close
      @ok="handleCancel"
      @cancel="handleCancel"
    >
      <template #title> 日志详细数据 </template>
      <a-descriptions
        style="margin-top: 20px"
        :data="logData"
        title="User Info"
        :column="1"
      />
    </a-drawer>

    <Breadcrumb :items="['menu.sysManage', 'menu.sysLog']" />
    <a-card>
      <a-row>
        <a-col :flex="1">
          <a-form
            :model="queryForm"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
            label-align="left"
          >
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item
                  field="queryTime"
                  :label="$t('sysLogQuery.form.queryTime')"
                >
                  <a-range-picker
                    v-model="queryForm.queryTime"
                    style="width: 100%"
                  />
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item
                  field="logType"
                  :label="$t('sysLogQuery.form.logType')"
                >
                  <a-input
                    v-model="queryForm.logType"
                    :placeholder="$t('sysLogQuery.form.logType')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item
                  field="username"
                  :label="$t('sysLogQuery.form.username')"
                >
                  <a-input
                    v-model="queryForm.username"
                    :placeholder="$t('sysLogQuery.form.username.prompt')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item field="role" :label="$t('sysLogQuery.form.role')">
                  <a-input
                    v-model="queryForm.role"
                    :placeholder="$t('sysLogQuery.form.role.prompt')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item
                  field="method"
                  :label="$t('sysLogQuery.form.method')"
                >
                  <a-input
                    v-model="queryForm.method"
                    :placeholder="$t('sysLogQuery.form.method')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>
            </a-row>
          </a-form>
        </a-col>
        <a-divider style="height: 84px" direction="vertical" />
        <a-col :flex="'60px'" style="text-align: start">
          <a-space direction="vertical" :size="18">
            <a-button type="primary" @click="search">
              <template #icon>
                <icon-search />
              </template>
              {{ $t('sysLogQuery.form.search') }}
            </a-button>
            <a-button @click="reset">
              <template #icon>
                <icon-refresh />
              </template>
              {{ $t('sysLogQuery.form.reset') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-divider style="margin-top: 0" />

      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data="sysLogDataList"
        :pagination="pagination"
        :bordered="{ cell: true }"
        column-resizable
        stripe
        @page-change="onPageChange"
      >
        <template #index="{ rowIndex }">
          {{ rowIndex + 1 + (pagination.page - 1) * pagination.pageSize }}
        </template>

        <template #operations="{ record, rowIndex }">
          <a-button type="text" size="small" @click="viewLog(record)">
            {{ $t('sysLogTable.columns.operations.detailed') }}
          </a-button>

          <a-popconfirm
            :content="$t('sysLogTable.columns.operations.delete.prompt')"
            type="warning"
            @ok="deleteSysLog(record.id, rowIndex)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('sysLogTable.columns.operations.delete') }}
            </a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import { Pagination } from '@/types/global';
  import { computed, reactive, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import {
    SysLogs,
    GetSysLogList,
    PolicyParams,
    DeleteSysLog,
  } from '@/api/sysManage/sysLog';
  import DateString from '@/utils/time';

  const generateFormModel = () => {
    return {
      logType: '',
      role: '',
      username: '',
      nickname: '',
      method: '',
      path: '',
      resCode: '',
      clientIP: '',
      requestAt: '',
      responseAt: '',
      latency: '',
      statusCode: '',
      bodySize: '',
      queryTime: undefined,
    };
  };

  const { t } = useI18n();
  const sysLogDataList = ref<SysLogs>([]);
  const loading = ref(false);
  const queryForm = ref(generateFormModel());
  const visible = ref(false);
  const logData = ref<any>([]);

  const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
  };

  const pagination = reactive({
    ...basePagination,
  });

  const getSysLogs = async (
    params: PolicyParams = { page: 1, pageSize: 20 }
  ) => {
    loading.value = true;
    try {
      const list = await GetSysLogList(params);
      sysLogDataList.value = list.data;
      pagination.total = list.total;
      pagination.page = params.page;
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      loading.value = false;
    }
  };

  const search = () => {
    getSysLogs({
      ...basePagination,
      ...queryForm.value,
    });
  };

  const reset = () => {
    queryForm.value = generateFormModel();
    getSysLogs();
  };

  const deleteSysLog = async (id: number, index: number) => {
    const res = await DeleteSysLog(id);
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      sysLogDataList.value.splice(index, 1);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const onPageChange = (current: number) => {
    getSysLogs({ ...basePagination, ...queryForm.value, page: current });
  };

  const handleCancel = () => {
    visible.value = false;
  };

  getSysLogs();

  const columns = computed<TableColumnData[]>(() => [
    {
      title: t('sysLogTable.columns.index'),
      dataIndex: 'index',
      slotName: 'index',
    },
    {
      title: t('sysLogTable.columns.logType'),
      dataIndex: 'logType',
      slotName: 'logType',
    },
    {
      title: t('sysLogTable.columns.role'),
      dataIndex: 'role',
      slotName: 'role',
    },
    {
      title: t('sysLogTable.columns.username'),
      dataIndex: 'username',
      slotName: 'username',
    },
    {
      title: t('sysLogTable.columns.nickname'),
      dataIndex: 'nickname',
      slotName: 'nickname',
    },
    {
      title: t('sysLogTable.columns.method'),
      dataIndex: 'method',
      slotName: 'method',
    },
    {
      title: t('sysLogTable.columns.path'),
      dataIndex: 'path',
      slotName: 'path',
    },
    {
      title: t('sysLogTable.columns.resCode'),
      dataIndex: 'resCode',
      slotName: 'resCode',
    },
    {
      title: t('sysLogTable.columns.clientIP'),
      dataIndex: 'clientIP',
      slotName: 'clientIP',
    },
    {
      title: t('sysLogTable.columns.requestAt'),
      dataIndex: 'requestAt',
      slotName: 'requestAt',
    },
    {
      title: t('sysLogTable.columns.responseAt'),
      dataIndex: 'responseAt',
      slotName: 'responseAt',
    },
    {
      title: t('sysLogTable.columns.latency'),
      dataIndex: 'latency',
      slotName: 'latency',
    },
    {
      title: t('sysLogTable.columns.statusCode'),
      dataIndex: 'statusCode',
      slotName: 'statusCode',
    },
    {
      title: t('sysLogTable.columns.bodySize'),
      dataIndex: 'bodySize',
      slotName: 'bodySize',
    },

    {
      title: t('sysLogTable.columns.operations'),
      dataIndex: 'operations',
      slotName: 'operations',
    },
  ]);

  const viewLog = (record: any) => {
    visible.value = true;
    logData.value = [
      {
        label: t('sysLogTable.columns.logType'),
        value: record.logType,
      },
      {
        label: t('sysLogTable.columns.role'),
        value: record.role,
      },
      {
        label: t('sysLogTable.columns.username'),
        value: record.username,
      },
      {
        label: t('sysLogTable.columns.nickname'),
        value: record.nickname,
      },
      {
        label: t('sysLogTable.columns.uid'),
        value: record.uid,
      },
      {
        label: t('sysLogTable.columns.method'),
        value: record.method,
      },
      {
        label: t('sysLogTable.columns.path'),
        value: record.path,
      },
      {
        label: t('sysLogTable.columns.reqData'),
        value: record.reqData,
      },
      {
        label: t('sysLogTable.columns.resCode'),
        value: record.resCode,
      },
      {
        label: t('sysLogTable.columns.resData'),
        value: record.resData,
      },
      {
        label: t('sysLogTable.columns.clientIP'),
        value: record.clientIP,
      },
      {
        label: t('sysLogTable.columns.requestAt'),
        value: DateString(record.requestAt),
      },
      {
        label: t('sysLogTable.columns.responseAt'),
        value: DateString(record.responseAt),
      },
      {
        label: t('sysLogTable.columns.latency'),
        value: record.latency,
      },
      {
        label: t('sysLogTable.columns.errorMessage'),
        value: record.errorMessage,
      },
      {
        label: t('sysLogTable.columns.statusCode'),
        value: record.statusCode,
      },
      {
        label: t('sysLogTable.columns.bodySize'),
        value: record.bodySize,
      },
    ];
  };
</script>

<style>
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
