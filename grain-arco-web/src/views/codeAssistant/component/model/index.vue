<template>
    <div class="container">
        <a-modal
            v-model:visible="dialogFormVisible"
            :title="dialogFormTitle"
            @cancel="addModelCancel"
            @before-ok="addModelConfirm"
        >
            <a-form  :model="modelForm">
                <a-form-item field="structName"
                 :label="$t('modelDialogForm.structName')"
                 :rules="[{ required: true, message: $t('modelDialogForm.error.structName')}]"
                >
                    <a-input
                        v-model="modelForm.structName"
                        :placeholder="$t('modelDialogForm.structName.prompt')"
                        allow-clear
                    />
                </a-form-item>

                <a-form-item
                    field="databaseName"
                    :label="$t('modelDialogForm.databaseName')"
                    :rules="[{ required: true,message: $t('modelDialogForm.error.databaseName')}]"
                >
                    <a-select
                        v-model="modelForm.databaseName"
                        :placeholder="$t('modelDialogForm.databaseName.prompt')"
                        allow-clear
                    >
                        <a-option
                            v-for="option in database"
                            :key="option.value"
                            :value="option.value"
                        >{{ option.title }}</a-option
                        >
                    </a-select>
                </a-form-item>
                <a-form-item field="nickname"
                     :label="$t('modelDialogForm.nickname')">
                    <a-input
                        v-model="modelForm.nickname"
                        :placeholder="$t('modelDialogForm.nickname.prompt')"
                        allow-clear
                    />
                </a-form-item>

                <a-form-item
                    field="queryTime"
                    :label="$t('modelDialogForm.queryTime')"
                >
                    <a-switch
                        v-model="modelForm.queryTime"
                        checked-value="yes"
                        unchecked-value="no"
                    >
                        <template #checked> ON </template>
                        <template #unchecked> OFF </template>
                    </a-switch>
                </a-form-item>

            </a-form>
        </a-modal>

        <a-drawer :width="720" :footer="false" :visible="viewFieldlVisible" @cancel="viewFieldCancel" unmountOnClose>
            <template #title>
                字段管理
            </template>
            <FieldFile :parentID="parentID" />
        </a-drawer>

        <a-modal
            v-model:visible="viewCodeVisible"
            title="预览代码"
            hide-cancel
            fullscreen
            ok-text="关闭"
            @ok="viewCodeDialogClose"
        >
            <a-space direction="vertical" size="large">
                <a-radio-group v-model="codeType" type="button">
                    <a-radio value="model">Model</a-radio>
                    <a-radio value="router">Router</a-radio>
                    <a-radio value="handle">Handle</a-radio>
                    <a-radio value="service">Service</a-radio>
                    <a-radio value="repo">Repo</a-radio>
                    <a-radio value="flutterModel">flutterModel</a-radio>
                    <a-radio value="flutterApi">flutterApi</a-radio>
                    <a-radio value="vue">Vue</a-radio>
                    <a-radio value="api">Api</a-radio>
                    <a-radio value="zhCn">i18n-Zh-CN</a-radio>
                    <a-button type="primary" @click="copyToClipboard(computedCode)">
                        <template #icon>
                            <icon-copy />
                        </template>
                        {{$t('copy')}}
                    </a-button>
                </a-radio-group>
            </a-space>
            <codemirror v-model:value="computedCode" :options="options" />
        </a-modal>

        <a-card>

            <a-row style="margin-bottom: 16px">
                <a-col :span="12">
                    <a-space>
                        <a-button type="primary" @click="createModelButtonClick()">
                            <template #icon>
                                <icon-plus />
                            </template>
                            {{$t('addModelButton.Title')}}
                        </a-button>
                    </a-space>
                </a-col>
            </a-row>

            <a-table
                row-key="id"
                :loading="loading"
                :columns="columns"
                :data="modelDataList"
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
                    <a-button type="text" size="small" @click="editModel(record)">
                        {{$t('modelTable.columns.operations.edit')}}
                    </a-button>

                    <a-button type="text" size="small" @click="viewField(record.id)">
                        {{$t('modelTable.columns.operations.viewField')}}
                    </a-button>

                    <a-button
                        type="text"
                        size="small"
                        :loading="ViewCodeLoading"
                        @click="viewCodeClick(record)"
                    >
                        {{ $t('modelTable.columns.operations.viewCode') }}
                    </a-button>

                    <a-popconfirm
                        :content="$t('modelTable.columns.operations.delete.prompt')"
                        type="warning"
                        @ok="deleteModel(record.id, rowIndex)"
                    >
                        <a-button type="text" size="small" status="danger">
                            {{$t('modelTable.columns.operations.delete')}}
                        </a-button>
                    </a-popconfirm>
                </template>
            </a-table>
        </a-card>
    </div>
</template>

<script lang="ts" setup>
import {
    Models,
    GetModelList,
    AddModel,
    UpdateModel,
    DeleteModel,
    ViewCode,
} from '@/api/codeAssistant/model';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import { Pagination } from '@/types/global';
import { computed, reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import { useI18n } from 'vue-i18n';
import FieldFile from '../field/index.vue';
import Codemirror from 'codemirror-editor-vue3';
import 'codemirror/mode/javaScript/javaScript.js';
import 'codemirror/theme/dracula.css';


const props = defineProps({
    parentID: {
        type: Number,
        required: true,
    }
})

const generateFormModel = () => {
    return {
        id: 0,
        parentId: 0,
        structName: '',
        queryTime: '',
        databaseName: '',
        nickname: '',

    };
};



const database = ref([
    { title: 'MySQL', value: 'mysql' },
    { title: 'MongoDB', value: 'mongo' },
]);

const isEdit = ref(false);
const parentID = ref(0);
const viewFieldlVisible = ref(false);
const dialogFormVisible = ref(false);
const dialogFormTitle = ref('添加项目模型');
let modelForm = reactive(generateFormModel());

const { t } = useI18n();
const modelDataList = ref<Models>([]);
const loading = ref(false);
const queryForm = ref(generateFormModel());

const viewCodeVisible = ref(false);
const ViewCodeLoading = ref(false);
const codeType = ref('model');
const codeData = ref({
    model: 'modelCode',
    repo: 'repoCode',
    handle: 'handleCode',
    service: 'serviceCode',
    flutterModel: 'flutterModel',
    flutterApi: 'flutterApi',
    router: 'routerCode',
    vue: 'vueCode',
    api: 'apiCode',
    zhCn: 'zhCnCode',
});

const options = {
    autorefresh: true, // 是否自动刷新
    smartIndent: true, // 自动缩进
    tabSize: 4, // 缩进单元格为 4 个空格
    mode: 'application/json',
    theme: 'dracula', // 主题，使用主体需要引入相应的 css 文件
    line: true, // 是否显示行数
    viewportMargin: Infinity, // 高度自适应
    highlightDifferences: true,
    autofocus: false,
    indentUnit: 2,
    readOnly: true, // 只读
    showCursorWhenSelecting: true,
    firstLineNumber: 1,
    matchBrackets: true,
};

const viewCodeDialogClose = () => {
    viewCodeVisible.value = false;
    codeData.value = {
        model: 'modelCode',
        repo: 'repoCode',
        handle: 'handleCode',
        flutterModel: 'flutterModel',
        flutterApi: 'flutterApi',
        service: 'serviceCode',
        router: 'routerCode',
        vue: 'vueCode',
        api: 'apiCode',
        zhCn: 'zhCnCode',
    };
};

const computedCode = computed(() => {
    switch (codeType.value) {
        case 'model':
            return codeData.value.model;
        case 'router':
            return codeData.value.router;
        case 'handle':
            return codeData.value.handle;
        case 'flutterModel':
            return codeData.value.flutterModel
        case 'flutterApi':
            return codeData.value.flutterApi
        case 'service':
            return codeData.value.service;
        case 'repo':
            return codeData.value.repo;
        case 'vue':
            return codeData.value.vue;
        case 'api':
            return codeData.value.api;
        case 'zhCn':
            return codeData.value.zhCn;
        default:
            return '';
    }
});


const basePagination: Pagination = {
    page: 1,
    pageSize: 20,
};

const pagination = reactive({
    ...basePagination,
});

const getModels = async (params: Pagination = { page: 1, pageSize: 20 }) => {
    loading.value = true;
    try {
        const list = await GetModelList(params,props.parentID);
        modelDataList.value = list.data;
        pagination.total = list.total;
        pagination.page = list.page as number;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        loading.value = false;
    }
};

const search = () => {
    getModels({
        ...basePagination,
        ...queryForm.value,
    });
};

const reset = () => {
    queryForm.value = generateFormModel();
    getModels();
};

const deleteModel = async (id: number, index: number) => {
    const res = await DeleteModel(id);
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        modelDataList.value.splice(index, 1);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const viewField = (mid: any) => {
  viewFieldlVisible.value = true;
  parentID.value = mid;
}

const editModel = (data: any) => {
    dialogFormTitle.value = t('modelTable.columns.operations.edit');
    modelForm.id = data.id;
    modelForm.parentId = data.parentId
    modelForm.structName = data.structName
    modelForm.queryTime = data.queryTime
    modelForm.databaseName = data.databaseName
    modelForm.nickname = data.nickname

    dialogFormVisible.value = true;
    isEdit.value = true;
};
const onPageChange = (current: number) => {
    getModels({ ...basePagination, page: current });
};

const createModelButtonClick = () => {
    modelForm.parentId = props.parentID
    dialogFormTitle.value = t('addModelButton.Title');
    dialogFormVisible.value = true;
    isEdit.value = false;
};
const clearForm = () => {
    modelForm.id = 0;
    modelForm.parentId = 0;
    modelForm.structName = '';
    modelForm.queryTime = '';
    modelForm.databaseName = '';
    modelForm.nickname = '';

};
const addModelCancel = () => {
    dialogFormVisible.value = false;
    clearForm();
};

const viewFieldCancel = () => {
    viewFieldlVisible.value = false;
};



const addModelConfirm = async () => {
    let res;
    if (isEdit.value) {
        res = await UpdateModel(modelForm);
    } else {
        res = await AddModel(modelForm);
    }
    if (res.success) {
        Message.success(res.message ?? '请求错误');
        onPageChange(pagination.page);
    } else {
        Message.warning(res.message ?? '请求错误');
    }
};

const viewCodeClick = async (v: any) => {
    ViewCodeLoading.value = true;
    try {
        const res = await ViewCode(v.id);
        ViewCodeLoading.value = false;
        codeData.value = res.data;
        codeType.value = 'model';
        viewCodeVisible.value = true;
    } catch (e) {
        codeType.value = 'model';
        ViewCodeLoading.value = false;
    }
};

function copyToClipboard(text: string) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    Message.success("复制成功");
}

getModels();

const columns = computed<TableColumnData[]>(() => [
    {
        title: t('modelTable.columns.index'),
        dataIndex: 'index',
        slotName: 'index',
    },
    {
        title: t('modelTable.columns.databaseName'),
        dataIndex: 'databaseName',
        slotName: 'databaseName',
    },
    {
        title: t('modelTable.columns.structName'),
        dataIndex: 'structName',
        slotName: 'structName',
    },
    {
        title: t('modelTable.columns.nickname'),
        dataIndex: 'nickname',
        slotName: 'nickname',
    },

    {
        title: t('modelTable.columns.operations'),
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
