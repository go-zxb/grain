import axios from 'axios';
import { Base } from '@/api/base';

export interface SysMenu {
  id: number;
  ParentId: number;
  path: string;
  name: string;
  meta: {
    locale: string;
    requiresAuth: boolean;
    icon: string;
    order: string;
    roles: string[];
    rolesStr: string[];
  };
  children: SysMenu[];
}

export type SysMenus = SysMenu[];

export interface PolicyParams extends Partial<SysMenu> {
  page: number;
  pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
  data: SysMenus;
}

export function GetSysMenuList(data: any) {
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
  return axios.get<PolicyListRes>(`/api/v1/sysMenu/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function GetSysMenuById(id: number) {
  return axios.get<SysMenu[]>(`/api/v1/sysMenu/${id}`).then((res) => {
    return res.data;
  });
}

export function AddSysMenu(data: any) {
  return axios.post<Base>(`/api/v1/sysMenu`, data).then((res) => {
    return res.data;
  });
}

export function UpdateSysMenu(data: any) {
  return axios.put<Base>(`/api/v1/sysMenu`, data).then((res) => {
    return res.data;
  });
}

export function DeleteSysMenu(data: number) {
  return axios.delete<Base>(`/api/v1/sysMenu?id=${data}`).then((res) => {
    return res.data;
  });
}
