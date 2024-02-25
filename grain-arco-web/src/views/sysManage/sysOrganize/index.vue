<template>
    <div class="container">
        <a-modal
            v-model:visible="dialogFormVisible"
            :title="dialogFormTitle"
            @cancel="addOrganizeCancel"
            @before-ok="addOrganizeConfirm"
        >
            <a-form :model="organizeForm">
                <a-form-item field="name"
                 :label="$t('organizeDialogForm.name')"
                 :rules="[
                  { required: true, message: $t('organizeDialogForm.error.name')}]">
                    <a-input
                        v-model="organizeForm.name"
                        :placeholder="$t('organizeDialogForm.name.prompt')"
                        allow-clear
                    />
                </a-form-item>
            </a-form>
        </a-modal>

        <a-drawer :width="720" :footer="false" :visible="viewDepartmentVisible" @cancel="viewDepartmentCancel" unmountOnClose>
            <template #title>
                部门管理
            </template>
            <department :receivedParam="paramFromA" />
        </a-drawer>

        <Breadcrumb :items="['menu.organize', 'menu.organize']" />
        <a-card>

            <a-divider style="margin-top: 0" />

            <a-row style="margin-bottom: 16px">
                <a-col :span="12">
                    <a-space>
                        <a-button type="primary" @click="createOrganizeButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{$t('addOrganizeButton.Title')}}
                        </a-button>
                    </a-space>
                </a-col>
            </a-row>

            <a-table
                row-key="id"
                :loading="loading"
                :columns="columns"
                :data="organizeDataList"
                :pagination="pagination"
                :bordered="{ cell: true }"
                column-resizable
                stripe
                @page-change="onPageChange"
            >
                <template #index="{ rowIndex }">
                    {{rowIndex + 1 + (pagination.page - 1) * pagination.pageSize}}
                </template>

                <template #operations="{ record, rowIndex }">
                    <a-button type="text" size="small" @click="editOrganize(record)">
                        {{$t('organizeTable.columns.operations.edit')}}
                    </a-button>

                    <a-button type="text" size="small" @click="viewDepartment(record)">
                        {{$t('organizeTable.columns.operations.department')}}
                    </a-button>

                    <a-popconfirm
                        :content="$t('organizeTable.columns.operations.delete.prompt')"
                        type="warning"
                        @ok="deleteOrganize(record.id, rowIndex)"
                    >
                        <a-button type="text" size="small" status="danger">
                            {{$t('organizeTable.columns.operations.delete')}}
                        </a-button>
                    </a-popconfirm>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import {
    Organizes,
    GetOrganizeList,
    PolicyParams,
    AddOrganize,
    UpdateOrganize,
    DeleteOrganize,
} from '@/api/sysManage/sysOrganize';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { Pagination } from '@/types/global';
import { computed, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import department from './components/department.vue';

const generateFormModel = () => {
    return {
        id: 0,
        parentId: 0,
        name: '',
        oeType: 0,
    };
};

const isEdit = ref(false);
const viewDepartmentVisible = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加组织管理');
let organizeForm = reactive(generateFormModel());

const { t } = useI18n();
const organizeDataList = ref<Organizes>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());
const paramFromA = ref('');

const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const getOrganizes = async (params: PolicyParams = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    try {
        const list = await GetOrganizeList(params);
        organizeDataList.value = list.data;
        pagination.total = list.total;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const search = () => {
    getOrganizes({
        ...basePagination,
        ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getOrganizes();
};

const deleteOrganize = async (id: number, index: number) => {
    const res = await DeleteOrganize(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        organizeDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const editOrganize = (data: any) => {
    dialogFormTitle.value = t('organizeTable.columns.operations.edit');
    organizeForm.id = data.id;
    organizeForm.parentId = data.parentId
    organizeForm.name = data.name
    organizeForm.oeType = data.oeType

    dialogFormVisible.value = true;
    isEdit.value = true;
};

const viewDepartment = (data: any) => {
    paramFromA.value = data.id;
    viewDepartmentVisible.value = true;
}
const onPageChange = (current: number) => {
    getOrganizes({ ...basePagination, page: current });
};

const createOrganizeButtonClick = () => {
    organizeForm.oeType = 1;
    dialogFormTitle.value = t('addOrganizeButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    organizeForm.id = 0;
    organizeForm.parentId = 0;
    organizeForm.name = '';
    organizeForm.oeType = 0;

};

const viewDepartmentCancel = () => {
    viewDepartmentVisible.value = false;
};

const addOrganizeCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const addOrganizeConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdateOrganize(organizeForm);
    } else {
        res = await AddOrganize(organizeForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        onPageChange(pagination.page);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

getOrganizes();

const columns = computed<TableColumnData[]>(() => [
    {
        title: t('organizeTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },
    {
        title: t('organizeTable.columns.name'),
        dataIndex: 'name',
        slotName: 'name',
    },

    {
        title: t('organizeTable.columns.operations'),
        dataIndex: 'operations',
        slotName: 'operations',
    },
]);

const departmentColumns = computed<TableColumnData[]>(() => [
    {
        title: t('departmentTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },
    {
        title: t('departmentTable.columns.name'),
        dataIndex: 'name',
        slotName: 'name',
    },
    {
        title: t('departmentTable.columns.leader'),
        dataIndex: 'leader',
        slotName: 'leader',
    },

    {
        title: t('departmentTable.columns.operations'),
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
