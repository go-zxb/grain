<template>
    <div class="container">
        <a-modal
            v-model:visible="dialogFormVisible"
            :title="dialogFormTitle"
            @cancel="addDepartmentCancel"
            @before-ok="addDepartmentConfirm"
        >
            <a-form :model="departmentForm">
                <a-form-item
                    field="name"
                    :label="$t('departmentDialogForm.name')"
                    :rules="[
            { required: true, message: $t('departmentDialogForm.error.name') },
          ]">
                    <a-input
                        v-model="departmentForm.name"
                        :placeholder="$t('departmentDialogForm.name.prompt')"
                        allow-clear
                    />
                </a-form-item>
                <a-form-item
                    field="leader"
                    :label="$t('departmentDialogForm.leader')"
                    :rules="[
            {
              required: true,
              message: $t('departmentDialogForm.error.leader'),
            },
          ]">
                    <a-input
                        v-model="departmentForm.leader"
                        :placeholder="$t('departmentDialogForm.leader.prompt')"
                        allow-clear
                    />
                </a-form-item>
            </a-form>
        </a-modal>

        <a-drawer
            :width="720"
            :footer="false"
            :visible="viewPositionVisible"
            unmount-on-close
            @cancel="viewOrganizeCancel"
        >
            <template #title> 职位管理 </template>
            <position :received-param="paramFromA" />
        </a-drawer>

        <Breadcrumb :items="['menu.department', 'menu.department']" />
        <a-card>
            <a-divider style="margin-top: 0" />
            <a-row style="margin-bottom: 16px">
                <a-col :span="12">
                    <a-space>
                        <a-button type="primary" @click="createDepartmentButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{ $t('addDepartmentButton.Title') }}
                        </a-button>
                    </a-space>
                </a-col>
            </a-row>

            <a-table
                row-key="id"
                :loading="loading"
                :columns="columns"
                :data="departmentDataList"
                :bordered="{ cell: true }"
                column-resizable
                stripe
            >
                <template #index="{ rowIndex }">
                    {{ rowIndex + 1 }}
                </template>

                <template #operations="{ record, rowIndex }">
                    <a-button type="text" size="small" @click="editDepartment(record)">
                        {{ $t('departmentTable.columns.operations.edit') }}
                    </a-button>

                    <a-button type="text" size="small" @click="viewPosition(record)">
                        {{ $t('departmentTable.columns.operations.position') }}
                    </a-button>

                    <a-popconfirm
                        :content="$t('departmentTable.columns.operations.delete.prompt')"
                        type="warning"
                        @ok="deleteDepartment(record.id, rowIndex)"
                    >
                        <a-button type="text" size="small" status="danger">
                            {{ $t('departmentTable.columns.operations.delete') }}
                        </a-button>
                    </a-popconfirm>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { computed, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import { defineProps } from 'vue/dist/vue';
import {
    Organizes,
    GetDepartmentList,
    AddOrganize,
    UpdateOrganize,
    DeleteOrganize,
} from '@/api/sysManage/sysOrganize';
import position from './position.vue';

const generateFormModel = () => {
    return {
        id: 0,
        parentId: 0,
        name: '',
        leader: '',
        oeType: 0,
    };
};

const isEdit = ref(false);
const viewPositionVisible = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加部门管理');
const departmentForm = reactive(generateFormModel());
const paramFromA = ref('');

const { t } = useI18n();
const departmentDataList = ref<Organizes>([]);
const loading = ref(false);

const props = defineProps({
    receivedParam: {
        type: String,
        required: true,
    },
});

const getDepartments = async () => {
    loading.value = true;
    try {
        const list = await GetDepartmentList(props.receivedParam);
        departmentDataList.value = list.data;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const deleteDepartment = async (id: number, index: number) => {
    const res = await DeleteOrganize(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        departmentDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const editDepartment = (data: any) => {
    dialogFormTitle.value = t('departmentTable.columns.operations.edit2');
    departmentForm.id = data.id;
    departmentForm.parentId = data.parentId;
    departmentForm.name = data.name;
    departmentForm.leader = data.leader;
    departmentForm.oeType = data.oeType;

    dialogFormVisible.value = true;
    isEdit.value = true;
};

const createDepartmentButtonClick = () => {
    departmentForm.parentId = props.receivedParam as unknown as number;
    departmentForm.oeType = 2;
    dialogFormTitle.value = t('addDepartmentButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    departmentForm.id = 0;
    departmentForm.parentId = 0;
    departmentForm.name = '';
    departmentForm.leader = '';
    departmentForm.oeType = 0;
};

const viewOrganizeCancel = () => {
    viewPositionVisible.value = false;
};

const viewPosition = (data: any) => {
    paramFromA.value = data.id;
    viewPositionVisible.value = true;
};

const addDepartmentCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const addDepartmentConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdateOrganize(departmentForm);
    } else {
        res = await AddOrganize(departmentForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        await getDepartments();
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

getDepartments();

const columns = computed<TableColumnData[]>(() => [
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
