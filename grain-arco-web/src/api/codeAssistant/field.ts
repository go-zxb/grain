import axios from 'axios';
import { Base } from '@/api/base';

export interface Field {
    id: number;
    parentId: number;
    name: string;
    mysqlField: string;
    type: string;
    jsonTag: string;
    description: string;
    validationRules: string;
    queryCriteria: string;
    required: string;
}

export type Fields = Field[];


export interface PolicyParams extends Partial<Field> {
    page: number;
    pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
    data: Fields;
}

export function GetFieldList(parentID: number) {
    return axios.get<PolicyListRes>(`/api/v1/codeAssistant/fields/list?parentId=${parentID}`).then((res) => {
        return res.data;
    });
}

export function AddField(data: any) {
    return axios.post<Base>(`/api/v1/codeAssistant/fields`, data).then((res) => {
        return res.data;
    });
}

export function UpdateField(data: any) {
    return axios.put<Base>(`/api/v1/codeAssistant/fields`, data).then((res) => {
        return res.data;
    });
}

export function DeleteField(data: number) {
    return axios.delete<Base>(`/api/v1/codeAssistant/fields?fid=${data}`).then((res) => {
        return res.data;
    });
}
