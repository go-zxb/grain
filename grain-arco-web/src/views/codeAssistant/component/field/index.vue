<template>
    <div class="container">
        <a-modal
            v-model:visible="dialogFormVisible"
            :title="dialogFormTitle"
            @cancel="addFieldCancel"
            @before-ok="addFieldConfirm"
        >
            <a-form :model="fieldForm">

                <a-form-item field="name"
                             :label="$t('fieldDialogForm.name')"
                             :rules="[
          { required: true,
            message: $t('fieldDialogForm.error.name')}]"
                >
                    <a-input
                        v-model="fieldForm.name"
                        :placeholder="$t('fieldDialogForm.name.prompt')"
                        allow-clear
                    />
                </a-form-item>
                <a-form-item field="mysqlField"
                             :label="$t('fieldDialogForm.mysqlField')"

                >
                    <a-input
                        v-model="fieldForm.mysqlField"
                        :placeholder="$t('fieldDialogForm.mysqlField.prompt')"
                        allow-clear
                    />
                </a-form-item>

                <a-form-item field="jsonTag"
                             :label="$t('fieldDialogForm.jsonTag')"

                >
                    <a-input
                        v-model="fieldForm.jsonTag"
                        :placeholder="$t('fieldDialogForm.jsonTag.prompt')"
                        allow-clear
                    />
                </a-form-item>

                <a-form-item field="description"
                             :label="$t('fieldDialogForm.description')"

                >
                    <a-input
                        v-model="fieldForm.description"
                        :placeholder="$t('fieldDialogForm.description.prompt')"
                        allow-clear
                    />
                </a-form-item>

                <a-form-item field="type" :label="$t('fieldDialogForm.type')">
                    <a-select
                        v-model="fieldForm.type"
                        :placeholder="$t('fieldDialogForm.type.prompt')"
                    >
                        <a-option
                            v-for="option in fieldTypes"
                            :key="option.value"
                            :value="option.value"
                        >{{ option.title }}</a-option
                        >
                    </a-select>
                </a-form-item>

                <a-form-item
                    field="requiredValue"
                    :label="$t('fieldDialogForm.validationRules')"
                >
                    <a-input
                        v-model="fieldForm.validationRules"
                        allow-clear
                        :placeholder="$t('fieldDialogForm.validationRules.prompt')"
                    />
                    <a-select
                        v-model="fieldForm.validationRules"
                        :placeholder="$t('fieldDialogForm.validationRules.prompt')"
                    >
                        <a-option
                            v-for="option in requiredValueOption"
                            :key="option.value"
                            :value="option.value"
                        >{{ option.title }}</a-option
                        >
                    </a-select>
                </a-form-item>

                <a-form-item
                    field="queryCriteria"
                    :label="$t('fieldDialogForm.queryCriteria')"
                >
                    <a-select
                        v-model="fieldForm.queryCriteria"
                        :placeholder="$t('fieldDialogForm.queryCriteria.prompt')"
                        allow-clear
                    >
                        <a-option
                            v-for="option in queryType"
                            :key="option.value"
                            :value="option.value"
                        >{{ option.title }}</a-option
                        >
                    </a-select>
                </a-form-item>

                <a-form-item
                    field="required"
                    :label="$t('fieldDialogForm.required.prompt')"
                >
                    <a-switch
                        v-model="fieldForm.required"
                        checked-value="yes"
                        unchecked-value="no"
                    >
                        <template #checked> ON </template>
                        <template #unchecked> OFF </template>
                    </a-switch>
                </a-form-item>

            </a-form>
        </a-modal>

        <a-card>

            <a-row style="margin-bottom: 16px">
                <a-col :span="12">
                    <a-space>
                        <a-button type="primary" @click="createFieldButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{$t('addFieldButton.Title')}}
                        </a-button>
                    </a-space>
                </a-col>
            </a-row>

            <a-table
                row-key="id"
                :loading="loading"
                :columns="columns"
                :data="fieldDataList"
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
                    <a-button type="text" size="small" @click="editField(record)">
                        {{$t('fieldTable.columns.operations.edit')}}
                    </a-button>

                    <a-popconfirm
                        :content="$t('fieldTable.columns.operations.delete.prompt')"
                        type="warning"
                        @ok="deleteField(record.id, rowIndex)"
                    >
                        <a-button type="text" size="small" status="danger">
                            {{$t('fieldTable.columns.operations.delete')}}
                        </a-button>
                    </a-popconfirm>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import {
    Fields,
    GetFieldList,
    AddField,
    UpdateField,
    DeleteField,
} from '@/api/codeAssistant/field';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { Pagination } from '@/types/global';
import { computed, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';

const generateFormModel = () => {
    return {
        id: 0,
        parentId: 0,
        name: '',
        mysqlField: '',
        type: '',
        jsonTag: '',
        description: '',
        validationRules: '',
        queryCriteria: '',
        required: '',

    };
};

const props = defineProps({
    parentID: {
        type: Number,
        required: true,
    }
})

const isEdit = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加数据字段');
let fieldForm = reactive(generateFormModel());

const { t } = useI18n();
const fieldDataList = ref<Fields>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());

const fieldTypes = ref([
    { title: 'string', value: 'string' },
    { title: 'byte', value: '[]byte' },
    { title: 'any', value: 'any' },
    { title: 'int', value: 'int' },
    { title: 'int8', value: 'int8' },
    { title: 'int16', value: 'int16' },
    { title: 'int32', value: 'int32' },
    { title: 'int64', value: 'int64' },
    { title: 'uint', value: 'uint' },
    { title: 'uint8', value: 'uint8' },
    { title: 'uint16', value: 'uint16' },
    { title: 'uint32', value: 'uint32' },
    { title: 'uint64', value: 'uint64' },
    { title: 'bool', value: 'bool' },
    { title: 'time', value: 'time.Time' },
]);

const queryType = ref([
    { title: '等于(=)', value: 'eq' },
    { title: '小于(<)', value: 'lt' },
    { title: '大于(>)', value: 'gt' },
    { title: '小于等于(<=)', value: 'lte' },
    { title: '大于等于(>=)', value: 'gte' },
    { title: '不等于(!=)', value: 'ne' },
    { title: '模糊(LIKE)', value: 'like' },
]);

const requiredValueOption = ref([
    { title: 'email', value: 'email' },
    { title: '不等于(ne)', value: 'ne=5' },
    { title: '小于(lt)', value: 'lt=100' },
    { title: '小于等于(lte)', value: 'lte=99' },
    { title: '大于(gt)', value: 'gt=10' },
    { title: '大于等于(gte)', value: 'gte=20' },
    { title: '等于(eq)', value: 'eq=66' },
    { title: '等于(len)', value: 'len=88' },
    { title: '数字范围', value: 'min=6,max=16' },
]);

const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const getFields = async (params: Pagination = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    console.log(props.parentID)
    try {
        const list = await GetFieldList(props.parentID);
        fieldDataList.value = list.data;
        pagination.total = list.data.length;
        pagination.page = 1;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const search = () => {
    getFields({
        ...basePagination,
        ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getFields();
};

const deleteField = async (id: number, index: number) => {
    const res = await DeleteField(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        fieldDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const editField = (data: any) => {
    dialogFormTitle.value = t('fieldTable.columns.operations.edit');
    fieldForm.id = data.id;
    fieldForm.parentId = data.parentId
    fieldForm.name = data.name
    fieldForm.mysqlField = data.mysqlField
    fieldForm.type = data.type
    fieldForm.jsonTag = data.jsonTag
    fieldForm.description = data.description
    fieldForm.validationRules = data.validationRules
    fieldForm.queryCriteria = data.queryCriteria
    fieldForm.required = data.required

    dialogFormVisible.value = true;
    isEdit.value = true;
};
const onPageChange = (current: number) => {
    getFields({ ...basePagination, page: current });
};

const createFieldButtonClick = () => {
    clearForm();
    fieldForm.parentId = props.parentID
    dialogFormTitle.value = t('addFieldButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    fieldForm.id = 0;
    fieldForm.parentId = 0;
    fieldForm.name = '';
    fieldForm.mysqlField = '';
    fieldForm.type = '';
    fieldForm.jsonTag = '';
    fieldForm.description = '';
    fieldForm.validationRules = '';
    fieldForm.queryCriteria = '';
    fieldForm.required = '';

};
const addFieldCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const addFieldConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdateField(fieldForm);
    } else {
        res = await AddField(fieldForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        onPageChange(pagination.page);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

getFields();

const columns = computed<TableColumnData[]>(() => [
    {
        title: t('fieldTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },

    {
        title: t('fieldTable.columns.name'),
        dataIndex: 'name',
        slotName: 'name',
    },
    {
        title: t('fieldTable.columns.type'),
        dataIndex: 'type',
        slotName: 'type',
    },

    {
        title: t('fieldTable.columns.operations'),
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
