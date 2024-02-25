import { Base } from '@/api/base';
import axios from 'axios';

export interface SysApi {
  id: number;
  path: string;
  group: string;
  method: string;
  description: string;
  children: SysApi[];
}

export type Apis = SysApi[];

export interface ApiGroupName {
  group: string;
}

export interface PolicyParams extends Partial<SysApi> {
  page: number;
  pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
  data: Apis;
}

export interface ApiGroup extends Partial<Base> {
  data: ApiGroupName[];
}

export interface ApiAndPermissions extends Partial<Base> {
  data: {
    authApi: number[];
    apiList: SysApi[];
  };
}

export function GetApiList(data: any) {
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
  return axios.get<PolicyListRes>(`/api/v1/sysApi/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function GetApiGroups() {
  return axios.get<ApiGroup>(`/api/v1/sysApi/apiGroups`).then((res) => {
    return res.data;
  });
}

export function AddApi(data: any) {
  return axios.post<Base>(`/api/v1/sysApi`, data).then((res) => {
    return res.data;
  });
}

export function UpdateApi(data: any) {
  return axios.put<Base>(`/api/v1/sysApi`, data).then((res) => {
    return res.data;
  });
}

export function DeleteApi(data: number) {
  return axios.delete<Base>(`/api/v1/sysApi?id=${data}`).then((res) => {
    return res.data;
  });
}

export function GetApiAndPermissions(role: any) {
  return axios
    .get<ApiAndPermissions>(`/api/v1/sysApi/apiAndPermissions?role=${role}`)
    .then((res) => {
      return res.data;
    });
}
