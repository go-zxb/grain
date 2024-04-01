<template>
  <div class="container">
    <a-modal
      v-model:visible="dialogFormVisible"
      :title="dialogFormTitle"
      @cancel="addRoleCancel"
      @before-ok="addRoleConfirm"
    >
      <a-form :model="roleForm">
        <a-form-item
          field="role"
          :label="$t('roleDialog.role')"
          :rules="[{ required: true, message: $t('roleForm.error.role') }]"
        >
          <a-input
            v-model="roleForm.role"
            :placeholder="$t('roleDialog.role.prompt')"
            allow-clear
          />
        </a-form-item>
        <a-form-item
          field="roleName"
          :label="$t('roleDialog.roleName')"
          :rules="[{ required: true, message: $t('roleForm.error.roleName') }]"
        >
          <a-input
            v-model="roleForm.roleName"
            :placeholder="$t('roleDialog.roleName.prompt')"
            allow-clear
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-drawer
      :width="720"
      :visible="drawerVisible"
      unmount-on-close
      @ok="setAuthOk"
      @cancel="drawerHandleCancel"
    >
      <template #title> {{ $t('roleTable.dialog.roleSet') }} </template>
      <a-tree
        v-model:checked-keys="checkedKeys"
        :checkable="true"
        :default-expand-all="false"
        :check-strictly="checkStrictly"
        :data="treeData"
      />
    </a-drawer>

      <a-drawer
              :width="720"
              :visible="drawerMenuVisible"
              unmount-on-close
              @ok="setMenuPermissionsSubmit"
              @cancel="setMenuPermissionsCancel"
      >
          <template #title> {{ '设置菜单权限' }} </template>
          <a-tree
                  v-model:checked-keys="checkedMenuKeys"
                  :checkable="true"
                  :default-expand-all="false"
                  :check-strictly="checkMenuStrictly"
                  :data="treeMenuData"
          />
      </a-drawer>

    <Breadcrumb :items="['menu.sysManage', 'menu.sysRole']" />
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
                <a-form-item field="role" :label="$t('query.form.role')">
                  <a-input
                    v-model="queryForm.role"
                    :placeholder="$t('query.form.role.prompt')"
                    allow-clear
                  />
                </a-form-item>
              </a-col>
              <a-col :span="8">
                <a-form-item
                  field="roleName"
                  :label="$t('query.form.roleName')"
                >
                  <a-input
                    v-model="queryForm.roleName"
                    :placeholder="$t('query.form.roleName.prompt')"
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
            <a-button type="primary" @click="createRoleButtonClick()">
              <template #icon>
                <icon-plus />
              </template>
              {{ $t('addRoleButton.Title') }}
            </a-button>
          </a-space>
        </a-col>
      </a-row>

      <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data="roleDataList"
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
          <a-button type="text" size="small" @click="editRole(record)">
            {{ $t('roleTable.columns.operations.edit') }}
          </a-button>

            <a-button type="text" size="small" @click="setMenu(record)">
                {{ '设置菜单权限' }}
            </a-button>

            <a-button type="text" size="small" @click="setAuth(record)">
                {{ '设置Api权限' }}
            </a-button>

          <a-popconfirm
            :content="$t('roleTable.columns.operations.delete.prompt')"
            type="warning"
            @ok="deleteRole(record.id, rowIndex)"
          >
            <a-button type="text" size="small" status="danger">
              {{ $t('roleTable.columns.operations.delete') }}
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
    Roles,
    GetRoleList,
    PolicyParams,
    AddRole,
    UpdateRole,
    DeleteRole,
    UpdateCasbin,
  } from '@/api/sysManage/sysRole';
  import { ApiAuthGroup } from '@/types/role';
  import { SysApi, GetApiAndPermissions } from '@/api/sysManage/sysApi';
  import {GetMenuAndPermission, SetMenuAndPermission} from "@/api/sysManage/sysMenu";

  const generateFormModel = () => {
    return {
      id: 0,
      role: '',
      roleName: '',
    };
  };

  const isEdit = ref(false);
  const drawerMenuVisible = ref(false);
  const checkedMenuKeys = ref<number[]>([]);
  const checkMenuStrictly = ref(false);
  const treeMenuData = ref<ApiAuthGroup[]>([]);

  const drawerVisible = ref(false);
  const dialogFormVisible = ref(false);
  const dialogFormTitle = ref('添加角色');
  const roleForm = reactive(generateFormModel());

  const { t } = useI18n();
  const roleDataList = ref<Roles>([]);
  const loading = ref(false);
  const queryForm = ref(generateFormModel());

  const role = ref(0);
  const checkedKeys = ref<number[]>([]);
  const checkStrictly = ref(false);
  const treeData = ref<ApiAuthGroup[]>([]);

  const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
  };

  const pagination = reactive({
    ...basePagination,
  });

  const getRoles = async (params: PolicyParams = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    try {
      const list = await GetRoleList(params);
      roleDataList.value = list.data;
      // pagination.total = list.total;
      // pagination.page = list.page as number;
    } catch (err) {
      // you can report use errorHandler or other
    } finally {
      loading.value = false;
    }
  };

  const search = () => {
    getRoles({
      ...basePagination,
      ...queryForm.value,
    });
  };

  const reset = () => {
    queryForm.value = generateFormModel();
    getRoles();
  };

  const deleteRole = async (id: number, index: number) => {
    const res = await DeleteRole(id);
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      roleDataList.value.splice(index, 1);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const setAuth = async (data: any) => {
    const res = await GetApiAndPermissions(data.role);
    if (res.success) {
      checkedKeys.value = [];
      checkedKeys.value = res.data.authApi;
      role.value = data.role;
      treeData.value = [];
      res.data.apiList.forEach((v1: SysApi) => {
        const c = ref<ApiAuthGroup[]>([]);
        v1.children.forEach((v2: SysApi) => {
          c.value.push({
            key: v2.id,
            title: v2.description,
            children: [],
          });
        });
        treeData.value.push({
          key: v1.id,
          title: v1.group,
          children: c.value,
        });
      });
      treeData.value.sort((a:any,b:any)=>{
          return a.key - b.key
      })
        drawerVisible.value = true;
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };

  const setAuthOk = async () => {
    const res = await UpdateCasbin({
      role: role.value.toString(),
      data: checkedKeys.value,
    });

    if (res.success) {
      drawerVisible.value = false;
      Message.success(res.message ?? 'ok');
    } else {
      Message.warning(res.message ?? '请求错误');
    }
    // console.log(checkedKeys.value)
  };
  const drawerHandleCancel = () => {
    drawerVisible.value = false;
    checkedKeys.value = [];
  };

  const editRole = (data: any) => {
    dialogFormTitle.value = t('roleTable.columns.operations.edit');
    roleForm.id = data.id;
    roleForm.role = data.role;
    roleForm.roleName = data.roleName;

    dialogFormVisible.value = true;
    isEdit.value = true;
  };

  const onPageChange = (current: number) => {
    getRoles({ ...basePagination, page: current });
  };

  const createRoleButtonClick = () => {
    dialogFormTitle.value = t('addRoleButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
  };
  const clearForm = () => {
    roleForm.id = 0;
    roleForm.role = '';
    roleForm.roleName = '';
  };
  const addRoleCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
  };

  const addRoleConfirm = async () => {
    let res;
    if (isEdit.value) {
      res = await UpdateRole(roleForm);
    } else {
      res = await AddRole(roleForm);
    }
    if (res.success) {
      Message.success(res.message ?? '请求错误');
      onPageChange(pagination.page);
    } else {
      Message.warning(res.message ?? '请求错误');
    }
  };



  const setMenuPermissionsSubmit = async () => {
      const res = await SetMenuAndPermission({"keys":checkedMenuKeys.value,"role":role.value});
      if (res.success) {
          Message.success(res.message ?? 'ok');
      }else {
          Message.warning(res.message ?? '请求错误');
      }
      drawerMenuVisible.value = false
  }

  const setMenu = async (data: any) => {
      role.value = data.role;
      const res = await GetMenuAndPermission(data.role);
      if (res.success) {
          console.log(res.data)
          treeMenuData.value = res.data
          checkedMenuKeys.value = res.data2
      }
      drawerMenuVisible.value = true
  }

  const setMenuPermissionsCancel = () => {
      drawerMenuVisible.value = false
  }

  getRoles();

  const columns = computed<TableColumnData[]>(() => [
    {
      title: t('roleTable.columns.index'),
      dataIndex: 'index',
      slotName: 'index',
    },
    {
      title: t('roleTable.columns.role'),
      dataIndex: 'role',
      slotName: 'role',
    },
    {
      title: t('roleTable.columns.roleName'),
      dataIndex: 'roleName',
      slotName: 'roleName',
    },

    {
      title: t('roleTable.columns.operations'),
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
