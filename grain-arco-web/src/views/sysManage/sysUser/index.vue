<template>
<div class="container">
    <a-modal
    v-model:visible="dialogFormVisible"
    :title="dialogFormTitle"
    @cancel="addSysUserCancel"
    @before-ok="addSysUserConfirm"
    >
    <a-form :model="sysUserForm">
        <a-form-item
        field="username"
        :label="$t('sysUserDialogForm.username')"
        :rules="[
        { required: true, message: $t('sysUserDialogForm.error.username') },
      ]"
        >
        <a-input
            v-model="sysUserForm.username"
            :placeholder="$t('sysUserDialogForm.username.prompt')"
            allow-clear
        />
        </a-form-item>
        <a-form-item
        field="password"
        :label="$t('sysUserDialogForm.password')"
        :rules="[
        { required: true, message: $t('sysUserDialogForm.error.password') },
      ]"
        >
        <a-input
            v-model="sysUserForm.password"
            :placeholder="$t('sysUserDialogForm.password.prompt')"
            allow-clear
        />
        </a-form-item>
        <a-form-item field="nickname" :label="$t('sysUserDialogForm.nickname')">
        <a-input
            v-model="sysUserForm.nickname"
            :placeholder="$t('sysUserDialogForm.nickname.prompt')"
            allow-clear
        />
        </a-form-item>
        <a-form-item field="email" :label="$t('sysUserDialogForm.email')">
        <a-input
            v-model="sysUserForm.email"
            :placeholder="$t('sysUserDialogForm.email.prompt')"
            allow-clear
        />
        </a-form-item>
        <a-form-item field="mobile" :label="$t('sysUserDialogForm.mobile')">
        <a-input
            v-model="sysUserForm.mobile"
            :placeholder="$t('sysUserDialogForm.mobile.prompt')"
            allow-clear
        />
        </a-form-item>
        <a-form-item field="roles" :label="$t('sysUserDialogForm.roles')">
        <a-cascader
            v-model="sysUserForm.roles"
            :options="roleStore.data"
            :style="{ width: '320px' }"
            placeholder="Please select ..."
            multiple
        />
        </a-form-item>

        <a-form-item field="organization" label="组织" >
            <a-select :style="{width:'200px'}" v-model="sysUserForm.organize" placeholder="请选择组织" allow-clear>
                <a-option v-for="value in organize">{{value}}</a-option>
            </a-select>
        </a-form-item>
        <a-form-item field="organization" label="部门" >
            <a-select :style="{width:'200px'}" v-model="sysUserForm.department" placeholder="请选择部门" allow-clear>
                <a-option v-for="value in department">{{value}}</a-option>
            </a-select>
        </a-form-item>
        <a-form-item field="organization" label="职位" >
            <a-select :style="{width:'200px'}" v-model="sysUserForm.position" placeholder="请选择职位" allow-clear>
                <a-option v-for="value in position">{{value}}</a-option>
            </a-select>
        </a-form-item>

        <a-form-item field="status" :label="$t('sysUserDialogForm.status')">
        <a-switch
            v-model="sysUserForm.status"
            checked-value="yes"
            unchecked-value="no"
        >
            <template #checked> ON </template>
            <template #unchecked> OFF </template>
        </a-switch>
        </a-form-item>
    </a-form>
    </a-modal>
    <Breadcrumb :items="['menu.sysUser', 'menu.sysUser']" />

    <a-layout class="layout-demo">
        <a-layout-sider
            hide-trigger
            collapsible
            :collapsed="collapsed"
        >
        <a-menu
            :defaultOpenKeys="[1]"
            :defaultSelectedKeys="[4]"
            :style="{ width: '100%' }"
            @menuItemClick="onClickMenuItem">
            <a-sub-menu v-for="menuItem in menu" :key="menuItem.id.toString()">
                <template #title v-if="menuItem.oeType==1">
                    <span><IconCalendar />{{ menuItem.name }}</span>
                 </template>
                <a-sub-menu v-for="submenu in menuItem.children" :key="submenu.id.toString()" :title="submenu.name">
                    <a-menu-item v-for="item in submenu.item" :key="menuItem.name+','+item.name">
                        {{item.name}}
                    </a-menu-item>
                </a-sub-menu>
            </a-sub-menu>
        </a-menu>
        </a-layout-sider>
<a-card :style="{ width: '100%' }">
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
                field="username"
                :label="$t('query.form.username')"
                >
                <a-input
                    v-model="queryForm.username"
                    :placeholder="$t('query.form.username.prompt')"
                    allow-clear
                />
                </a-form-item>
            </a-col>
            <a-col :span="8">
                <a-form-item field="email" :label="$t('query.form.email')">
                <a-input
                    v-model="queryForm.email"
                    :placeholder="$t('query.form.email.prompt')"
                    allow-clear
                />
                </a-form-item>
            </a-col>
            <a-col :span="8">
                <a-form-item field="mobile" :label="$t('query.form.mobile')">
                <a-input
                    v-model="queryForm.mobile"
                    :placeholder="$t('query.form.mobile.prompt')"
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
            <a-button type="primary" @click="createSysUserButtonClick()">
            <template #icon>
                <icon-plus />
            </template>
            {{ $t('addSysUserButton.Title') }}
            </a-button>
        </a-space>
        </a-col>
    </a-row>

    <a-table
        row-key="id"
        :loading="loading"
        :columns="columns"
        :data="sysUserDataList"
        :pagination="pagination"
        :bordered="{ cell: true }"
        column-resizable
        stripe
        @page-change="onPageChange"
    >
        <template #index="{ rowIndex }">
        {{ rowIndex + 1 + (pagination.page - 1) * pagination.pageSize }}
         </template>

        <template #roles="{ record }">
        <div v-for="item in record.roles" :key="item">
            <a-space wrap>
            <a-popconfirm
                :content="$t('sysUserTable.columns.setDefaultRole')"
                type="warning"
                @ok="setDefaultRole(record, item)"
            >
                <a-button
                :type="record.role === item ? 'primary' : 'outline'"
                size="small"
                >
                {{ findRoleStr(item) }}
                </a-button>
            </a-popconfirm>
            </a-space>
        </div>
        </template>

        <template #status="{ record }">
        <span v-if="record.status !== 'yes'" class="circle error"></span>
        <span v-else class="circle pass"></span>
        {{ $t(`sysUserTable.columns.status.${record.status}`) }}
    </template>

    <template #operations="{ record, rowIndex }">
        <a-button type="text" size="small" @click="editSysUser(record)">
            {{ $t('sysUserTable.columns.operations.edit') }}
        </a-button>

        <a-popconfirm
            :content="$t('sysUserTable.columns.operations.delete.prompt')"
            type="warning"
            @ok="deleteSysUser(record.id, rowIndex)"
        >
            <a-button type="text" size="small" status="danger">
            {{ $t('sysUserTable.columns.operations.delete') }}
            </a-button>
        </a-popconfirm>
    </template>
</a-table>
</a-card>
    </a-layout>
</div>
</template>

<script lang="ts" setup>
import {
    SysUsers,
    GetSysUserList,
    PolicyParams,
    AddSysUser,
    UpdateSysUser,
    DeleteSysUser,
    SetDefaultRole,
    GetOrganizeListGroup,
} from '@/api/sysManage/sysUser';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { Pagination } from '@/types/global';
import { computed, reactive, ref } from 'vue';
import { CascaderOption, Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import { useRolesStore } from '@/store';
import {Organize, UserGetOrganizeList} from "@/api/sysManage/sysOrganize";

const generateFormModel = () => {
    return {
    id: 0,
    username: '',
    password: '',
    nickname: '',
    email: '',
    mobile: '',
    roles: '',
    status: '',
    queryTime: '',
    organize:'',
    department:'',
    position:'',
    };
};

const menu = ref<Organize[]>([]);
const isEdit = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加系统用户');
const sysUserForm = reactive(generateFormModel());

const { t } = useI18n();
const sysUserDataList = ref<SysUsers>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());
const roleStore = useRolesStore();
const collapsed = ref(false);
const onCollapse = () => {
    collapsed.value = !collapsed.value;
};

const onClickMenuItem = async (key: any) => {
    console.log(key)
    queryForm.value.organize = key
    search()
};

const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const organize = ref<string[]>([]);
const department = ref<string[]>([]);
const position = ref<string[]>([]);

const getSysUsers = async (
    params: PolicyParams = { page: 1, pageSize: 20 }
) => {
    loading.value = true;
    try {
    const list = await GetSysUserList(params);
    sysUserDataList.value = list.data;
    pagination.page = params.page;
    } catch (err) {
    // you can report use errorHandler or other
    } finally {
    loading.value = false;
    }
};

const getOrganizeListGroup = async () => {
    try {
        const list = await GetOrganizeListGroup();
        console.log(list.data)
        menu.value = list.data
        menu.value.forEach((v: any)=>{
           organize.value.push(v.name)
            v.children.forEach((vv:any)=>{
                 department.value.push(vv.name);
                vv.item.forEach((vvv:any)=>{
                    position.value.push(vvv.name);
                })
            })
        })

    } catch (e) {
    } finally {
    }
}

const search = () => {
    getSysUsers({
    ...basePagination,
    ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getSysUsers();
};

const deleteSysUser = async (id: number, index: number) => {
    const res = await DeleteSysUser(id);
    if (res.success) {
    Message.success(res.message ?? '请求错误');
    sysUserDataList.value.splice(index, 1);
    } else {
    Message.warning(res.message ?? '请求错误');
    }
};

const editSysUser = (data: any) => {
    dialogFormTitle.value = '编辑系统用户';
    sysUserForm.id = data.id;
    sysUserForm.username = data.username;
    sysUserForm.password = '';
    sysUserForm.nickname = data.nickname;
    sysUserForm.email = data.email;
    sysUserForm.mobile = data.mobile;
    sysUserForm.roles = data.roles;
    sysUserForm.status = data.status;
    sysUserForm.organize = data.organize;
    sysUserForm.department = data.department;
    sysUserForm.position = data.position;

    dialogFormVisible.value = true;
    isEdit.value = true;
};
const onPageChange = (current: number) => {
    getSysUsers({ ...basePagination, page: current });
};

const createSysUserButtonClick = () => {
    dialogFormTitle.value = t('addSysUserButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    sysUserForm.id = 0;
    sysUserForm.username = '';
    sysUserForm.password = '';
    sysUserForm.nickname = '';
    sysUserForm.email = '';
    sysUserForm.mobile = '';
    sysUserForm.roles = '';
    sysUserForm.status = '';
    sysUserForm.queryTime = '';
    sysUserForm.organize = '';
    sysUserForm.department = '';
    sysUserForm.position = '';
};
const addSysUserCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const addSysUserConfirm = async () => {
    let res;
    if (isEdit.value) {
    res = await UpdateSysUser(sysUserForm);
    } else {
    res = await AddSysUser(sysUserForm);
    }
    if (res.success) {
    Message.success(res.message ?? '请求错误');
    onPageChange(pagination.page);
    } else {
    Message.warning(res.message ?? '请求错误');
    }
};

const setDefaultRole = async (record: any, newRole: string) => {
    const role = { id: record.id, role: newRole };
    const res = await SetDefaultRole(role);
    if (res.success) {
    record.role = newRole;
    Message.success(res.message ?? '请求错误');
    } else {
    Message.warning(res.message ?? '请求错误');
    }
};

const getRoles = () => {
    roleStore.getRoles();
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

getSysUsers();
getOrganizeListGroup();
getRoles();

const columns = computed<TableColumnData[]>(() => [
    {
    title: t('sysUserTable.columns.index'),
    dataIndex: 'index',
    slotName: 'index',
    },
    {
    title: t('sysUserTable.columns.username'),
    dataIndex: 'username',
    slotName: 'username',
    },
    {
    title: t('sysUserTable.columns.nickname'),
    dataIndex: 'nickname',
    slotName: 'nickname',
    },
    {
    title: t('sysUserTable.columns.email'),
    dataIndex: 'email',
    slotName: 'email',
    },
    {
    title: t('sysUserTable.columns.mobile'),
    dataIndex: 'mobile',
    slotName: 'mobile',
    },
    {
        title: t('sysUserTable.columns.roles'),
        dataIndex: 'roles',
        slotName: 'roles',
    },
    {
        title: t('sysUserTable.columns.organize'),
        dataIndex: 'organize',
        slotName: 'organize',
    },
    {
        title: t('sysUserTable.columns.department'),
        dataIndex: 'department',
        slotName: 'department',
    },
    {
        title: t('sysUserTable.columns.position'),
        dataIndex: 'position',
        slotName: 'position',
    },
    {
    title: t('sysUserTable.columns.status'),
    dataIndex: 'status',
    slotName: 'status',
    },

    {
    title: t('sysUserTable.columns.operations'),
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

<style scoped>
.layout-demo {
    height: 500px;
    background: var(--color-fill-2);
    border: 1px solid var(--color-border);
}
.layout-demo :deep(.arco-layout-sider) .logo {
    height: 32px;
    margin: 12px 8px;
    background: rgba(255, 255, 255, 0.2);
}
.layout-demo :deep(.arco-layout-sider-light) .logo{
    background: var(--color-fill-2);
}
.layout-demo :deep(.arco-layout-header)  {
    height: 64px;
    line-height: 64px;
    background: var(--color-bg-3);
}
.layout-demo :deep(.arco-layout-footer) {
    height: 48px;
    color: var(--color-text-2);
    font-weight: 400;
    font-size: 14px;
    line-height: 48px;
}
.layout-demo :deep(.arco-layout-content) {
    color: var(--color-text-2);
    font-weight: 400;
    font-size: 14px;
    background: var(--color-bg-3);
}
.layout-demo :deep(.arco-layout-footer),
.layout-demo :deep(.arco-layout-content)  {
    display: flex;
    flex-direction: column;
    justify-content: center;
    color: var(--color-white);
    font-size: 16px;
    font-stretch: condensed;
    text-align: center;
}
</style>