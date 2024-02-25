import axios from 'axios';
import { Base } from '@/api/base';

export interface Model {
    id: number;
    parentId: number;
    structName: string;
    queryTime: string;
    databaseName: string;
    nickname: string;
}

export type Models = Model[];


export interface PolicyParams extends Partial<Model> {
    page: number;
    pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
    data: Models;
}

export function GetModelList(data: any,parentId: any) {
    const uri = Object.entries(data)
        .map(
            ([key, value]) =>
                `${encodeURIComponent(key)}=${encodeURIComponent(value as string)}`
        )
        .join('&');
    return axios.get<PolicyListRes>(`/api/v1/codeAssistant/models/list?${uri}&parentId=${parentId}`).then((res) => {
        return res.data;
    });
}

export function GetModelById(id: number) {
    return axios.get<Model[]>(`/api/v1/codeAssistant/models/${id}`).then((res) => {
        return res.data;
    });
}

export function AddModel(data: any) {
    return axios.post<Base>(`/api/v1/codeAssistant/models`, data).then((res) => {
        return res.data;
    });
}

export function UpdateModel(data: any) {
    return axios.put<Base>(`/api/v1/codeAssistant/models`, data).then((res) => {
        return res.data;
    });
}

export function DeleteModel(data: number) {
    return axios.delete<Base>(`/api/v1/codeAssistant/models?mid=${data}`).then((res) => {
        return res.data;
    });
}

export function ViewCode(data: any) {
    return axios.get(`/api/v1/codeAssistant/viewCode?mid=${data}`).then((res) => {
        return res.data;
    });
}
