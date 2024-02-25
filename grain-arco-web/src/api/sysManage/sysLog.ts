import axios from 'axios';
import { Base } from '@/api/base';

export interface SysLog {
  id: number;
  typeName: string;
  role: string;
  username: string;
  nickname: string;
  method: string;
  path: string;
  resCode: string;
  clientIP: string;
  requestAt: string;
  responseAt: string;
  latency: string;
  statusCode: string;
  bodySize: string;
}

export type SysLogs = SysLog[];

export interface PolicyParams extends Partial<SysLog> {
  page: number;
  pageSize: number;
}

export interface PolicyListRes extends Partial<Base> {
  data: SysLogs;
}

export function GetSysLogList(data: any) {
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
  return axios.get<PolicyListRes>(`/api/v1/sysLog/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function GetSysLogById(id: number) {
  return axios.get<SysLog[]>(`/api/v1/sysLog/${id}`).then((res) => {
    return res.data;
  });
}

export function AddSysLog(data: any) {
  return axios.post<Base>(`/api/v1/sysLog`, data).then((res) => {
    return res.data;
  });
}

export function UpdateSysLog(data: any) {
  return axios.put<Base>(`/api/v1/sysLog`, data).then((res) => {
    return res.data;
  });
}

export function DeleteSysLog(data: number) {
  return axios.delete<Base>(`/api/v1/sysLog?id=${data}`).then((res) => {
    return res.data;
  });
}
