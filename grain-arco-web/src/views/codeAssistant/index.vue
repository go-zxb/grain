<template>
    <div class="container">
        <a-modal
                v-model:visible="dialogFormVisible"
                :title="dialogFormTitle"
                @cancel="addCodeAssistantCancel"
                @before-ok="addCodeAssistantConfirm"
        >
            <a-form :model="codeAssistantForm">
                <a-form-item field="projectName"
                             :label="$t('codeAssistantDialogForm.projectName')"
                             :rules="[
          { required: true,
            message: $t('codeAssistantDialogForm.error.projectName')}]"
                >
                    <a-input
                            v-model="codeAssistantForm.projectName"
                            :placeholder="$t('codeAssistantDialogForm.projectName.prompt')"
                            allow-clear
                    />
                </a-form-item>
                <a-form-item field="projectPath"
                             :label="$t('codeAssistantDialogForm.projectPath')"

                >
                    <a-input
                            v-model="codeAssistantForm.projectPath"
                            :placeholder="$t('codeAssistantDialogForm.projectPath.prompt')"
                            allow-clear
                    />
                </a-form-item>
                <a-form-item field="webProjectPath"
                             :label="$t('codeAssistantDialogForm.webProjectPath')"

                >
                    <a-input
                            v-model="codeAssistantForm.webProjectPath"
                            :placeholder="$t('codeAssistantDialogForm.webProjectPath.prompt')"
                            allow-clear
                    />
                </a-form-item>
                <a-form-item field="description"
                             :label="$t('codeAssistantDialogForm.description')"

                >
                    <a-input
                            v-model="codeAssistantForm.description"
                            :placeholder="$t('codeAssistantDialogForm.description.prompt')"
                            allow-clear
                    />
                </a-form-item>

            </a-form>
        </a-modal>

        <div v-if="modelType==2">
            <Breadcrumb :items="['menu.codeAssistant', 'menu.codeAssistant','menu.codeAssistantModel']" />
            <a-card>
                <modelFile :parentID="parentID"
                           @backProject="backProject"/>
            </a-card>
        </div>

        <div v-if="modelType==1">
            <Breadcrumb :items="['menu.codeAssistant', 'menu.codeAssistant']" />
            <a-card>

                <a-row style="margin-bottom: 16px">
                    <a-row :span="12">
                        <a-button type="primary" @click="createCodeAssistantButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{$t('addCodeAssistantButton.Title')}}
                        </a-button>
                    </a-row>
                </a-row>

                <a-table
                    row-key="id"
                    :loading="loading"
                    :columns="columns"
                    :data="codeAssistantDataList"
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
                        <a-button type="text" size="small" @click="editCodeAssistant(record)">
                            {{$t('codeAssistantTable.columns.operations.edit')}}
                        </a-button>

                        <a-button type="text" size="small" @click="viewModel(record.id)">
                            {{$t('codeAssistantTable.columns.operations.viewModel')}}
                        </a-button>

                        <a-popconfirm
                            :content="$t('codeAssistantTable.columns.operations.delete.prompt')"
                            type="warning"
                            @ok="deleteCodeAssistant(record.id, rowIndex)"
                        >
                            <a-button type="text" size="small" status="danger">
                                {{$t('codeAssistantTable.columns.operations.delete')}}
                            </a-button>
                        </a-popconfirm>
                    </template>
                </a-table>
            </a-card>
        </div>
    </div>
</template>

<script lang="ts" setup>
import {
    CodeAssistants,
    GetCodeAssistantList,
    AddCodeAssistant,
    UpdateCodeAssistant,
    DeleteCodeAssistant,
} from '@/api/codeAssistant/project';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { Pagination } from '@/types/global';
import { computed, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import modelFile from './component/model/index.vue';

const generateFormModel = () => {
    return {
        id: 0,
        projectName: '',
        projectPath: '',
        webProjectPath: '',
        description: '',

    };
};

const isEdit = ref(false);
const viewModelVisible = ref(false);
const modelType = ref(1);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加代码助手');
const parentID = ref(0);
let codeAssistantForm = reactive(generateFormModel());

const { t } = useI18n();
const codeAssistantDataList = ref<CodeAssistants>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());
const emit = defineEmits(['backProject']);

const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const getCodeAssistants = async (params: Pagination = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    try {
        const list = await GetCodeAssistantList(params);
        codeAssistantDataList.value = list.data;
        pagination.total = list.total;
        pagination.page = list.page as number;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const search = () => {
    getCodeAssistants({
        ...basePagination,
        ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getCodeAssistants();
};

const deleteCodeAssistant = async (id: number, index: number) => {
    const res = await DeleteCodeAssistant(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        codeAssistantDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const editCodeAssistant = (data: any) => {
    dialogFormTitle.value = t('codeAssistantTable.columns.operations.edit');
    codeAssistantForm.id = data.id;
    codeAssistantForm.projectName = data.projectName
    codeAssistantForm.projectPath = data.projectPath
    codeAssistantForm.webProjectPath = data.webProjectPath
    codeAssistantForm.description = data.description

    dialogFormVisible.value = true;
    isEdit.value = true;
};
const onPageChange = (current: number) => {
    getCodeAssistants({ ...basePagination, page: current });
};

const createCodeAssistantButtonClick = () => {
    dialogFormTitle.value = t('addCodeAssistantButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};

const backProject = () => {
  modelType.value = 1;
  emit('backProject', modelType.value);
}
const clearForm = () => {
    codeAssistantForm.id = 0;
    codeAssistantForm.projectName = '';
    codeAssistantForm.projectPath = '';
    codeAssistantForm.webProjectPath = '';
    codeAssistantForm.description = '';

};
const addCodeAssistantCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const viewModelCancel = () => {
    viewModelVisible.value = false;
};

const viewModel = (p5ID: any) => {
    viewModelVisible.value = true;
    parentID.value = p5ID;
    modelType.value=2;
};

const addCodeAssistantConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdateCodeAssistant(codeAssistantForm);
    } else {
        res = await AddCodeAssistant(codeAssistantForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        onPageChange(pagination.page);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

getCodeAssistants();

const columns = computed<TableColumnData[]>(() => [
    {
        title: t('codeAssistantTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },
    {
        title: t('codeAssistantTable.columns.projectName'),
        dataIndex: 'projectName',
        slotName: 'projectName',
    },
    {
        title: t('codeAssistantTable.columns.projectPath'),
        dataIndex: 'projectPath',
        slotName: 'projectPath',
    },
    {
        title: t('codeAssistantTable.columns.webProjectPath'),
        dataIndex: 'webProjectPath',
        slotName: 'webProjectPath',
    },
    {
        title: t('codeAssistantTable.columns.description'),
        dataIndex: 'description',
        slotName: 'description',
    },

    {
        title: t('codeAssistantTable.columns.operations'),
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
