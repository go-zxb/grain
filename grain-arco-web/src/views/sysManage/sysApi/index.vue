<template>
  <div class="container">
    <a-modal
      v-model:visible="dialogFormVisible"
      :title="dialogFormTitle"
      @cancel="addApiCancel"
      @before-ok="addApiConfirm"
    >
      <a-form :model="apiForm">
        <a-form-item field="path" :label="$t('apiDialog.path')">
          <a-input
            v-model="apiForm.path"
            :placeholder="$t('apiDialog.path.prompt')"
            allow-clear
          />
        </a-form-item>

        <a-form-item field="description" :label="$t('apiDialog.description')">
          <a-input
            v-model="apiForm.description"
            :placeholder="$t('apiDialog.description.prompt')"
            allow-clear
          />
        </a-form-item>

        <a-form-item field="method" :label="$t('apiDialog.method')">
          <a-select
            v-model="apiForm.method"
            :options="methodOptions"
            :placeholder="$t('apiDialog.method')"
            allow-clear
          />
        </a-form-item>

        <a-form-item field="group" :label="$t('apiDialog.group')">
          <a-select
            v-model="apiForm.group"
            :placeholder="$t('apiDialog.group')"
          >
            <a-option
              v-for="option in groupOptions"
              :key="option.group"
              :value="option.group"
              >{{ option.group }}</a-option
            >
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>
    <Breadcrumb :items="['menu.sysManage', 'menu.sysApi']" />
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
                <a-form-item field="path" :label="$t('query.form.path')">
                  <a-input
                    v-model="queryForm.path"
                    :placeholder="$t('query.form.path')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item field="filterType" :label="$t('query.form.group')">
                  <a-select
                    v-model="queryForm.group"
                    :placeholder="$t('query.form.group')"
                    allow-clear
                  >
                    <a-option
                      v-for="option in groupOptions"
                      :key="option.group"
                      :value="option.group"
                      >{{ option.group }}</a-option
                    >
                  </a-select>
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item field="status" :label="$t('query.form.method')">
                  <a-select
                    v-model="queryForm.method"
                    :placeholder="$t('query.form.method')"
                    allow-clear
                  >
                    <a-option
                      v-for="option in methodOptions"
                      :key="option.label"
                      :value="option.value as string"
                      >{{ option.label }}</a-option
                    >
                  </a-select>
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
              {{ $t('query.form.search') }}
            </a-button>
            <a-button @click="reset">
              <template #icon>
                <icon-refresh />
              </template>
              {{ $t('query.form.reset') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-divider style="margin-top: 0" />

      <a-row style="margin-bottom: 16px">
        <a-col :span="12">
          <a-space>
            <a-button type="primary" @click="createApiButtonClick()">
              <template #icon>
                <icon-plus />
              </template>
              {{ $t('addApiButton.Title') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data="apiDataList"
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
          <a-button type="text" size="small" @click="editApi(record)">
            {{ $t('apiTable.columns.operations.edit') }}
          </a-button>

          <a-popconfirm
            :content="$t('apiTable.columns.operations.delete.prompt')"
            type="warning"
            @ok="deleteApi(record.id, rowIndex)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('apiTable.columns.operations.delete') }}
            </a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import {
    ApiGroupName,
    Apis,
    GetApiList,
    GetApiGroups,
    PolicyParams,
    AddApi,
    UpdateApi,
    DeleteApi,
  } from '@/api/sysManage/sysApi';
  import { Pagination } from '@/types/global';
  import { computed, reactive, ref } from 'vue';
  import { Message, SelectOptionData } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';

  const generateFormModel = () => {
    return {
      id: 0,
      path: '',
      description: '',
      group: '',
      method: '',
    };
  };

  const isEdit = ref(false);
  const dialogFormVisible = ref(false);
  const dialogFormTitle = ref('添加Api');
  const apiForm = reactive(generateFormModel());

  const { t } = useI18n();
  const apiDataList = ref<Apis>([]);
  const loading = ref(false);
  const queryForm = ref(generateFormModel());
  const groupOptions = ref<ApiGroupName[]>([]);

  const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
  };

  const pagination = reactive({
    ...basePagination,
  });

  const getApis = async (params: PolicyParams = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    try {
      const list = await GetApiList(params);
      apiDataList.value = list.data;
      pagination.total = list.total;
      pagination.page = list.page as number;
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      loading.value = false;
    }
  };

  const search = () => {
    getApis({
      ...basePagination,
      ...queryForm.value,
    });
  };

  const reset = () => {
    queryForm.value = generateFormModel();
    getApis();
  };

  const apiGroup = async () => {
    const res = await GetApiGroups();
    groupOptions.value = res.data;
  };

  const deleteApi = async (id: number, index: number) => {
    const res = await DeleteApi(id);
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      apiDataList.value.splice(index, 1);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const editApi = (data: any) => {
    dialogFormTitle.value = '编辑Api';
    apiForm.id = data.id;
    apiForm.path = data.path;
    apiForm.description = data.description;
    apiForm.group = data.group;
    apiForm.method = data.method;
    dialogFormVisible.value = true;
    isEdit.value = true;
  };
  const onPageChange = (current: number) => {
    getApis({ ...basePagination, ...queryForm.value, page: current });
  };

  const createApiButtonClick = () => {
    clearForm();
    dialogFormTitle.value = t('addApiButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
  };
  const clearForm = () => {
    apiForm.id = 0;
    apiForm.path = '';
    apiForm.description = '';
    apiForm.group = '';
    apiForm.method = '';
  };
  const addApiCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
  };

  const addApiConfirm = async () => {
    let res;
    if (isEdit.value) {
      res = await UpdateApi(apiForm);
    } else {
      res = await AddApi(apiForm);
    }
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      onPageChange(pagination.page);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  getApis();
  apiGroup();

  const columns = [
    {
      title: '#',
      dataIndex: 'index',
      slotName: 'index',
    },
    {
      title: '分组',
      dataIndex: 'group',
    },
    {
      title: '路径',
      dataIndex: 'path',
    },
    {
      title: '方法',
      dataIndex: 'method',
    },
    {
      title: '描述',
      dataIndex: 'description',
    },
    {
      title: '操作',
      dataIndex: 'operations',
      slotName: 'operations',
    },
  ];

  const methodOptions = computed<SelectOptionData[]>(() => [
    {
      label: 'GET',
      value: 'GET',
    },
    {
      label: 'POST',
      value: 'POST',
    },
    {
      label: 'PUT',
      value: 'PUT',
    },
    {
      label: 'DELETE',
      value: 'DELETE',
    },
    {
      label: 'PATCH',
      value: 'PATCH',
    },
    {
      label: 'OPTIONS',
      value: 'OPTIONS',
    },
  ]);
</script>

<style>
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
