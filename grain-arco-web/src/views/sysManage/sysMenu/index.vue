<template>
  <div class="container">
    <a-modal
      v-model:visible="dialogFormVisible"
      :title="dialogFormTitle"
      @cancel="addSysMenuCancel"
      @before-ok="addSysMenuConfirm"
    >
      <a-form :model="sysMenuForm">
        <a-form-item
          field="path"
          :label="$t('sysMenuDialogForm.path')"
          :rules="[
            {
              required: true,
              message: $t('sysMenuDialogForm.error.path'),
            },
          ]"
        >
          <a-input
            v-model="sysMenuForm.path"
            :placeholder="$t('sysMenuDialogForm.path.prompt')"
            allow-clear
          />
        </a-form-item>
        <a-form-item
          field="name"
          :label="$t('sysMenuDialogForm.name')"
          :rules="[
            { required: true, message: $t('sysMenuDialogForm.error.name') },
          ]"
        >
          <a-input
            v-model="sysMenuForm.name"
            :placeholder="$t('sysMenuDialogForm.name.prompt')"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="i18n" :label="$t('sysMenuDialogForm.i18n')">
          <a-input
            v-model="sysMenuForm.meta.i18n"
            :placeholder="$t('sysMenuDialogForm.i18n.prompt')"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="roles" :label="$t('sysMenuDialogForm.roles')">
          <a-cascader
            v-model="sysMenuForm.meta.roles"
            :options="roleStore.data"
            :style="{ width: '320px' }"
            placeholder="Please select ..."
            multiple
          />
        </a-form-item>
        <a-form-item field="icon" :label="$t('sysMenuDialogForm.icon')">
          <a-input
            v-model="sysMenuForm.meta.icon"
            :placeholder="$t('sysMenuDialogForm.icon.prompt')"
            allow-clear
          />
        </a-form-item>
        <a-form-item field="order" :label="$t('sysMenuDialogForm.order')">
          <a-input-number
            v-model="sysMenuForm.meta.order"
            :placeholder="$t('sysMenuDialogForm.order.prompt')"
            allow-clear
          />
        </a-form-item>
<!--        <a-form-item field="parentId" :label="$t('sysMenuDialogForm.parentId')">-->
<!--          <a-input-number-->
<!--            v-model="sysMenuForm.parentId"-->
<!--            :placeholder="$t('sysMenuDialogForm.parentId.prompt')"-->
<!--            allow-clear-->
<!--          />-->
<!--        </a-form-item>-->
      </a-form>
    </a-modal>

    <Breadcrumb :items="['menu.sysManage', 'menu.sysMenu']" />
    <a-card>
      <a-row style="margin-bottom: 16px">
        <a-col :span="12">
          <a-space>
            <a-button type="primary" @click="createSysMenuButtonClick()">
              <template #icon>
                <icon-plus />
              </template>
              {{ $t('addSysMenuButton.Title') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data="sysMenuDataList"
        :pagination="pagination"
        :bordered="{ cell: true }"
        column-resizable
        stripe
        @page-change="onPageChange"
      >
        <template #index="{ rowIndex }">
          {{ rowIndex + 1 + (pagination.page - 1) * pagination.pageSize }}
        </template>

        <template #i18n="{ record }">
          {{ $t(record.meta.i18n) }}
        </template>

        <template #roles="{ record }">
          <a-space wrap>
            <a-tag
              v-for="(roleID, index) of record.meta.roles"
              :key="index"
              color="green"
              >{{ findRoleStr(roleID) }}</a-tag
            >
          </a-space>
        </template>

        <template #operations="{ record, rowIndex }">

          <a-button v-if="record.parentId===0" type="text" size="small" @click="addChildMenu(record)">
              {{ $t('sysMenuTable.columns.operations.create') }}
          </a-button>

          <a-button type="text" size="small" @click="editSysMenu(record)">
            {{ $t('sysMenuTable.columns.operations.edit') }}
          </a-button>

          <a-popconfirm
            :content="$t('sysMenuTable.columns.operations.delete.prompt')"
            type="warning"
            @ok="deleteSysMenu(record.id, rowIndex)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('sysMenuTable.columns.operations.delete') }}
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
  import { CascaderOption, Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import {
    SysMenu,
    GetSysMenuList,
    PolicyParams,
    AddSysMenu,
    UpdateSysMenu,
    DeleteSysMenu,
  } from '@/api/sysManage/sysMenu';
  import { useRolesStore } from '@/store';

  const generateFormModel = () => {
    return {
      id: 0,
      parentId: 0,
      path: '',
      name: '',
      meta: {
        i18n: '',
        roles: [],
        icon: '',
        order: 0,
      },
    };
  };

  const isEdit = ref(false);
  const dialogFormVisible = ref(false);
  const dialogFormTitle = ref('添加菜单');
  const sysMenuForm = reactive(generateFormModel());

  const { t } = useI18n();
  const sysMenuDataList = ref<SysMenu[]>([]);
  const loading = ref(false);
  const roleStore = useRolesStore();

  const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
  };

  const pagination = reactive({
    ...basePagination,
  });

  const getSysMenus = async (
    params: PolicyParams = { page: 1, pageSize: 20 }
  ) => {
    loading.value = true;
    try {
      const list = await GetSysMenuList(params);
        let m: SysMenu[] = [];
        list.data.sort((a: any, b: any) => {
            return a.meta.order - b.meta.order;
        });

        list.data.forEach((item)=>{
            if(item.children!==null){
                item.children.sort((a: any, b: any) => {
                    return a.meta.order - b.meta.order;
                });
            }
            m.push(item)
        })

        sysMenuDataList.value = m;
      // pagination.total = list.total;
      pagination.page = params.page as number;
      pagination.pageSize = params.pageSize;
    } catch (err) {
      // you can report use errorHandler or other
        console.log(err)
    } finally {
      loading.value = false;
    }
  };

  const deleteSysMenu = async (id: number, index: number) => {
    const res = await DeleteSysMenu(id);
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      sysMenuDataList.value.splice(index, 1);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const addChildMenu = (data: any) => {
      dialogFormTitle.value = '创建子菜单';
      sysMenuForm.parentId = data.id;
      dialogFormVisible.value = true;
      isEdit.value = false;
  };

  const editSysMenu = (data: any) => {
    dialogFormTitle.value = '编辑菜单';
    sysMenuForm.id = data.id;
    sysMenuForm.path = data.path;
    sysMenuForm.name = data.name;
    sysMenuForm.meta.i18n = data.meta.i18n;
    sysMenuForm.meta.roles = data.meta.roles;
    sysMenuForm.meta.icon = data.meta.icon;
    sysMenuForm.meta.order = data.meta.order;
    sysMenuForm.parentId = data.parentId;
    dialogFormVisible.value = true;
    isEdit.value = true;
  };
  const onPageChange = (current: number) => {
    getSysMenus({ ...basePagination, page: current });
  };

  const createSysMenuButtonClick = () => {
      clearForm()
      dialogFormTitle.value = t('addSysMenuButton.Title');
      dialogFormVisible.value = true;
      isEdit.value = false;
  };

  const clearForm = () => {
    sysMenuForm.id = 0;
    sysMenuForm.path = '';
    sysMenuForm.name = '';
    sysMenuForm.meta.i18n = '';
    sysMenuForm.meta.roles = [];
    sysMenuForm.meta.icon = '';
    sysMenuForm.meta.order = 0;
    sysMenuForm.parentId = 0;
  };
  const addSysMenuCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
  };

  const addSysMenuConfirm = async () => {
    let res;
    if (isEdit.value) {
      res = await UpdateSysMenu(sysMenuForm);
    } else {
      res = await AddSysMenu(sysMenuForm);
    }
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      onPageChange(pagination.page);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const findRoleStr = (role: string) => {
    let str;
    roleStore.data.forEach((r: CascaderOption) => {
      if (r.value === role) {
        str = r.label;
      }
    });
    return str;
  };

  const getRoles = () => {
    roleStore.getRoles();
  };

  getSysMenus();
  getRoles();

  const columns = computed<TableColumnData[]>(() => [
    {
      title: t('sysMenuTable.columns.index'),
      dataIndex: 'index',
      slotName: 'index',
    },
    {
      title: t('sysMenuTable.columns.path'),
      dataIndex: 'path',
      slotName: 'path',
    },
    {
      title: t('sysMenuTable.columns.name'),
      dataIndex: 'name',
      slotName: 'name',
    },
    {
      title: t('sysMenuTable.columns.i18n'),
      dataIndex: 'meta.i18n',
      slotName: 'i18n',
    },
    {
      title: t('sysMenuTable.columns.roles'),
      dataIndex: 'meta.roles',
      slotName: 'roles',
    },
    {
      title: t('sysMenuTable.columns.icon'),
      dataIndex: 'meta.icon',
      slotName: 'icon',
    },
    {
      title: t('sysMenuTable.columns.order'),
      dataIndex: 'meta.order',
      slotName: 'order',
    },
    // {
    //   title: t('sysMenuTable.columns.parentId'),
    //   dataIndex: 'parentId',
    //   slotName: 'parentId',
    // },

    {
      title: t('sysMenuTable.columns.operations'),
      dataIndex: 'operations',
      slotName: 'operations',
    },
  ]);
</script>

<style>
  .container {
    padding: 0 20px 20px 20px;
  }
</style>
