import axios from 'axios';
import { Base } from '@/api/base';

export interface SysRole {
  id: number;
  role: string;
  roleName: string;
}

export type Roles = SysRole[];

export interface PolicyParams extends Partial<SysRole> {
  page: number;
  pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
  data: Roles;
}

export function GetRoleList(data: any) {
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
  return axios.get<PolicyListRes>(`/api/v1/sysRole/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function GetRoleById(id: number) {
  return axios.get<SysRole[]>(`/api/v1/sysRole/${id}`).then((res) => {
    return res.data;
  });
}

export function AddRole(data: any) {
  return axios.post<Base>(`/api/v1/sysRole`, data).then((res) => {
    return res.data;
  });
}

export function UpdateRole(data: any) {
  return axios.put<Base>(`/api/v1/sysRole`, data).then((res) => {
    return res.data;
  });
}

export function DeleteRole(data: number) {
  return axios.delete<Base>(`/api/v1/sysRole?id=${data}`).then((res) => {
    return res.data;
  });
}

export function UpdateCasbin(data: any) {
  return axios.put<Base>('/api/v1/casbin', data).then((res) => {
    return res.data;
  });
}
