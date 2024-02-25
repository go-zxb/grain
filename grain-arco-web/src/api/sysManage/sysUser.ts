import axios from 'axios';
import { Base } from '@/api/base';

export interface SysUser {
  id: number;
  userName: string;
  password: string;
  nickname: string;
  email: string;
  mobile: string;
  roles: string;
  status: string;
}

export type SysUsers = SysUser[];

export interface PolicyParams extends Partial<SysUser> {
  page: number;
  pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
  data: SysUsers;
}

export function GetSysUserList(data: any) {
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
  return axios.get<PolicyListRes>(`/api/v1/sysUser/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function GetOrganizeListGroup() {
  return axios.get(`/api/v1/organize/listGroup?qType=1`).then((res) => {
    return res.data;
  });
}

export function GetSysUserById(id: number) {
  return axios.get<SysUser[]>(`/api/v1/sysUser/${id}`).then((res) => {
    return res.data;
  });
}

export function AddSysUser(data: any) {
  return axios.post<Base>(`/api/v1/sysUser/create`, data).then((res) => {
    return res.data;
  });
}

export function UpdateSysUser(data: any) {
  return axios.put<Base>(`/api/v1/sysUser/editUserInfo`, data).then((res) => {
    return res.data;
  });
}

export function DeleteSysUser(data: number) {
  return axios.delete<Base>(`/api/v1/sysUser?id=${data}`).then((res) => {
    return res.data;
  });
}

export function SetDefaultRole(data: any) {
  return axios.put<Base>('/api/v1/sysUser/setDefaultRole', data).then((res) => {
    return res.data;
  });
}