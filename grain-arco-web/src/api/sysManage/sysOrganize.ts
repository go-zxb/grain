import axios from 'axios';
import { Base } from '@/api/base';

export interface Organize {
    id: number;
    parentId: number;
    name: string;
    leader: string;
    oeType: number;
    children:Item[];
}

export interface Item {
    name: string;
    id:   number;
    item: Item[];
}

export type Organizes = Organize[];


export interface PolicyParams extends Partial<Organize> {
    page: number;
    pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
    data: Organizes;
}

export function GetOrganizeList(data: any) {
    const uri = Object.entries(data)
        .map(
            ([key, value]) =>
                `${encodeURIComponent(key)}=${encodeURIComponent(value as string)}`
        )
        .join('&');
    return axios.get<PolicyListRes>(`/api/v1/organize/list?${uri}&qType=1`).then((res) => {
        return res.data;
    });
}

export function UserGetOrganizeList(data: any) {
    return axios.get<PolicyListRes>(`/api/v1/organize/list?id=${data}`).then((res) => {
        return res.data;
    });
}

export function AddOrganize(data: any) {
    return axios.post<Base>(`/api/v1/organize`, data).then((res) => {
        return res.data;
    });
}

export function UpdateOrganize(data: any) {
    return axios.put<Base>(`/api/v1/organize`, data).then((res) => {
        return res.data;
    });
}

export function DeleteOrganize(data: number) {
    return axios.delete<Base>(`/api/v1/organize/organizeById?id=${data}`).then((res) => {
        return res.data;
    });
}

export function GetDepartmentList(data: any) {
    return axios.get<PolicyListRes>(`/api/v1/organize/list?qType=2&id=${data}`).then((res) => {
        return res.data;
    });
}

export function GetPositionList(data: any, id: any) {
    const filteredData = Object.fromEntries(
        Object.entries(data).filter(
            ([, value]) =>
                value !== undefined && value !== null && value !== '' && value !== 0
        )
    );
    const uri = Object.entries(filteredData)
        .map(
            ([key, value]) =>
                `${encodeURIComponent(key)}=${encodeURIComponent(value as string)}`
        )
        .join('&');
    return axios.get<PolicyListRes>(`/api/v1/organize/list?${uri}&qType=3&id=${id}`).then((res) => {
        return res.data;
    });
}

export function AddPosition(data: any) {
    return axios.post<Base>(`/api/v1/organize`, data).then((res) => {
        return res.data;
    });
}

export function UpdatePosition(data: any) {
    return axios.put<Base>(`/api/v1/organize`, data).then((res) => {
        return res.data;
    });
}

export function DeletePosition(data: number) {
    return axios.delete<Base>(`/api/v1/organize/organizeById?id=${data}`).then((res) => {
        return res.data;
    });
}

