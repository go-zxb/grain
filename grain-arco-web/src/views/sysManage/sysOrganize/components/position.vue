<template>
    <div class="container">
        <a-modal
            v-model:visible="dialogFormVisible"
            :title="dialogFormTitle"
            @cancel="addPositionCancel"
            @before-ok="addPositionConfirm"
        >
            <a-form :model="positionForm">
                <a-form-item
                    field="name"
                    :label="$t('positionDialogForm.name')"
                    :rules="[
            { required: true, message: $t('positionDialogForm.error.name') },
          ]"
                >
                    <a-input
                        v-model="positionForm.name"
                        :placeholder="$t('positionDialogForm.name.prompt')"
                        allow-clear
                    />
                </a-form-item>
            </a-form>
        </a-modal>

        <Breadcrumb :items="['menu.position', 'menu.position']" />
        <a-card>
            <a-row style="margin-bottom: 16px">
                <a-col :span="12">
                    <a-space>
                        <a-button type="primary" @click="createPositionButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{ $t('addPositionButton.Title') }}
                        </a-button>
                    </a-space>
                </a-col>
            </a-row>

            <a-table
                row-key="id"
                :loading="loading"
                :columns="columns"
                :data="positionDataList"
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
                    <a-button type="text" size="small" @click="editPosition(record)">
                        {{ $t('positionTable.columns.operations.edit') }}
                    </a-button>

                    <a-popconfirm
                        :content="$t('positionTable.columns.operations.delete.prompt')"
                        type="warning"
                        @ok="deletePosition(record.id, rowIndex)"
                    >
                        <a-button type="text" size="small" status="danger">
                            {{ $t('positionTable.columns.operations.delete') }}
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
import { computed, defineProps, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import * as child_process from 'child_process';
import {
    Organizes,
    GetPositionList,
    PolicyParams,
    AddPosition,
    UpdatePosition,
    DeletePosition,
} from '@/api/sysManage/sysOrganize';

const props = defineProps({
    receivedParam: {
        type: String,
        required: true,
    },
});

const generateFormModel = () => {
    return {
        id: 0,
        name: '',
        parentId: 0,
        oeType: 0,
    };
};

const isEdit = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加职位');
const positionForm = reactive(generateFormModel());

const { t } = useI18n();
const positionDataList = ref<Organizes>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());

const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const getPositions = async (
    params: Pagination = { page: 1, pageSize: 20 }
) => {
    loading.value = true;
    try {
        const list = await GetPositionList(params, props.receivedParam);
        positionDataList.value = list.data;
        pagination.total = list.total;
        pagination.page = list.page as number;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const search = () => {
    getPositions({
        ...basePagination,
        ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getPositions();
};

const deletePosition = async (id: number, index: number) => {
    const res = await DeletePosition(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        positionDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const editPosition = (data: any) => {
    dialogFormTitle.value = t('positionTable.columns.operations.edit');
    positionForm.id = data.id;
    positionForm.name = data.name;
    positionForm.parentId = data.parentId;
    positionForm.oeType = data.oeType;

    dialogFormVisible.value = true;
    isEdit.value = true;
};
const onPageChange = (current: number) => {
    getPositions({ ...basePagination, page: current });
};

const createPositionButtonClick = () => {
    positionForm.oeType = 3;
    positionForm.parentId = props.receivedParam as unknown as number;
    dialogFormTitle.value = t('addPositionButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    positionForm.id = 0;
    positionForm.name = '';
    positionForm.parentId = 0;
    positionForm.oeType = 0;
};
const addPositionCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const addPositionConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdatePosition(positionForm);
    } else {
        res = await AddPosition(positionForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        onPageChange(pagination.page);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

getPositions();

const columns = computed<TableColumnData[]>(() => [
    {
        title: t('positionTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },
    {
        title: t('positionTable.columns.name'),
        dataIndex: 'name',
        slotName: 'name',
    },

    {
        title: t('positionTable.columns.operations'),
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
