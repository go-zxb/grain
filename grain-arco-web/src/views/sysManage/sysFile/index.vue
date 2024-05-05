<template>
  <div class="container">
    <Breadcrumb :items="['menu.sysManage', 'menu.sysFile']" />
    <a-modal
      v-model:visible="visible"
      :title="dialogTitle"
      hide-cancel
      @ok="handleOk"
    >
      <a-upload
        draggable
        action="/"
        list-type="picture"
        :limit="100"
        :custom-request="customRequest"
      />
    </a-modal>

    <a-card class="general-card" title=" "
      ><a-row>
        <a-col :flex="1">
          <a-form
            :model="searchForm"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
            label-align="left"
          >
            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item
                  field="queryTime"
                  :label="$t('uploadTable.searchForm.queryTime')"
                >
                  <a-range-picker
                    v-model="searchForm.queryTime"
                    style="width: 100%"
                  />
                </a-form-item>
              </a-col>

              <a-col :span="8">
                <a-form-item
                  field="fileName"
                  :label="$t('uploadTable.searchForm.fileName')"
                >
                  <a-input
                    v-model="searchForm.fileName"
                    :placeholder="$t('uploadTable.searchForm.fileName.des')"
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
              {{ $t('uploadTable.searchForm.query') }}
            </a-button>
            <a-button @click="reset">
              <template #icon>
                <icon-refresh />
              </template>
              {{ $t('uploadTable.searchForm.reset') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-divider style="margin-top: 0" />
      <a-row style="margin-bottom: 16px">
        <a-col :span="12">
          <a-space>
            <a-button
              type="primary"
              @click="createHandleClick($t('uploadTable.dialog.create'), null)"
            >
              <template #icon>
                <icon-plus />
              </template>
              {{ $t('uploadTable.create') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>
      <a-table
        row-key="id"
        :loading="loading"
        :pagination="pagination"
        :columns="(cloneColumns as TableColumnData[])"
        :data="renderData"
        :size="size"
        column-resizable
        :bordered="{ cell: true }"
        @page-change="onPageChange"
      >
        <template #index="{ rowIndex }">
          {{ rowIndex + 1 + (pagination.page - 1) * pagination.pageSize }}
        </template>

          <template #fileUrl="{ record }">
              <a-image
                      width="50"
                      :src="record.fileUrl"
              />
          </template>

        <template #operations="{ record, rowIndex }">
          <a-popconfirm
            :content="$t('uploadTable.columns.operations.down.des')"
            type="info"
            @ok="downClick(record)"
          >
            <a-button type="text" size="small">
              {{ $t('uploadTable.columns.operations.down') }}
            </a-button>
          </a-popconfirm>

            <a-button type="text" size="small" @click="copyTextToClipboard(record.fileUrl)">
                {{ "复制URL" }}
            </a-button>

          <a-popconfirm
            :content="$t('uploadTable.columns.operations.delete.des')"
            type="warning"
            @ok="deleteUpload(record.id, rowIndex)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('uploadTable.columns.operations.delete') }}
            </a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, ref, reactive, watch } from 'vue';
  import { useI18n } from 'vue-i18n';
  import useLoading from '@/hooks/loading';
  import { Pagination } from '@/types/global';
  import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
  import cloneDeep from 'lodash/cloneDeep';
  import { Message } from '@arco-design/web-vue';
  import {
      uploadFile,
      GetUploadList,
      DeleteUploadById,
      SysFile,
      UploadReqData,
  } from '@/api/sysManage/sysFile';
  import { RequestOption } from '@arco-design/web-vue/es/upload/interfaces';
  import {copyTextToClipboard} from "@/utils/sys";

  type SizeProps = 'mini' | 'small' | 'medium' | 'large';
  type Column = TableColumnData & { checked?: true };

  const generateFormModel = () => {
    return {
      id: 0,
      fileName: '',
      fileUrl: '',
      queryTime: undefined,
    };
  };
  const { loading, setLoading } = useLoading(true);
  const { t } = useI18n();
  const renderData = ref<SysFile[]>([]);
  const searchForm = ref(generateFormModel());
  const cloneColumns = ref<Column[]>([]);
  const size = ref<SizeProps>('medium');

  const visible = ref(false);
  const dialogTitle = ref('');

  const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
  };
  const pagination = reactive({
    ...basePagination,
  });

  const columns = computed<TableColumnData[]>(() => [
    {
      title: t('uploadTable.columns.index'),
      dataIndex: 'index',
      slotName: 'index',
      width:80,
    },

    {
      title: t('uploadTable.columns.fileUrl'),
      dataIndex: 'fileUrl',
      slotName: 'fileUrl',
    },
      {
          title: t('uploadTable.columns.fileName'),
          dataIndex: 'fileName',
          slotName: 'fileName',
      },

    {
      title: t('uploadTable.columns.fileType'),
      dataIndex: 'fileType',
      slotName: 'fileType',
    },

    {
      title: t('uploadTable.columns.operations'),
      dataIndex: 'operations',
      slotName: 'operations',
    },
  ]);

  const fetchData = async (
    params: UploadReqData = { page: 1, pageSize: 20 }
  ) => {
    setLoading(true);
    try {
      const data = await GetUploadList(params);
      renderData.value = data.data;
      pagination.page = params.page;
      pagination.total = data.total;
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      setLoading(false);
    }
  };

  const search = () => {
    fetchData({
      ...basePagination,
      ...searchForm.value,
    } as unknown as UploadReqData);
  };
  const onPageChange = (current: number) => {
    fetchData({ ...basePagination, page: current });
  };

  fetchData();
  const reset = () => {
    searchForm.value = generateFormModel();
    fetchData();
  };

  const deleteUpload = async (val: any, index: number) => {
    const res = await DeleteUploadById(val);
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      renderData.value.splice(index, 1);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const downClick = (data: any) => {
    const url = import.meta.env.VITE_API_BASE_URL;
    window.open(data.fileUrl, '_blank');
  };
  const createHandleClick = (title: string, data: any) => {
    dialogTitle.value = title;
    visible.value = true;
  };
  const customRequest = (options: RequestOption) => {
    // docs: https://axios-http.com/docs/cancellation
    const controller = new AbortController();

    (async function requestWrap() {
      const {
        onProgress,
        onError,
        onSuccess,
        fileItem,
        name = 'file',
      } = options;
      onProgress(20);
      const formData = new FormData();
      formData.append(name as string, fileItem.file as Blob);
      const onUploadProgress = (event: ProgressEvent) => {
        let percent;
        if (event.total > 0) {
          percent = (event.loaded / event.total) * 100;
        }
        onProgress(parseInt(String(percent), 10), event);
      };

      try {
        // https://github.com/axios/axios/issues/1630
        // https://github.com/nuysoft/Mock/issues/127

        const res = await uploadFile(formData, {
          controller,
          onUploadProgress,
        });
        onSuccess(res);
      } catch (error) {
        onError(error);
      }
    })();
    return {
      abort() {
        controller.abort();
      },
    };
  };
  const handleOk = () => {
    visible.value = false;
  };

  watch(
    () => columns.value,
    (val) => {
      cloneColumns.value = cloneDeep(val);
      cloneColumns.value.forEach((item, index) => {
        item.checked = true;
      });
    },
    { deep: true, immediate: true }
  );
</script>

<script lang="ts">
  export default {
    name: 'SearchTable',
  };
</script>

<style scoped lang="less">
  .container {
    padding: 0 20px 20px 20px;
  }
  :deep(.arco-table-th) {
    &:last-child {
      .arco-table-th-item-title {
        margin-left: 16px;
      }
    }
  }
  .action-icon {
    margin-left: 12px;
    cursor: pointer;
  }
  .active {
    color: #0960bd;
    background-color: #e3f4fc;
  }
  .setting {
    display: flex;
    align-items: center;
    width: 200px;
    .title {
      margin-left: 12px;
      cursor: pointer;
    }
  }
</style>
