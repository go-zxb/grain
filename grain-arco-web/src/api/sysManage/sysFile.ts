import axios from 'axios';
import { Base } from '@/api/base';

export interface SysFile {
  id: number;
  fileName: string;
  fileUrl: string;
}

export interface UploadReqData extends Partial<SysFile> {
  page: number;
  pageSize: number;
}

export interface UploadRes extends Partial<Base> {
  data: SysFile[];
}

export function DeleteUploadByIds(data: number[]) {
  return axios
    .delete<Base>(`/api/v1/upload/ids`, {
      data: {
        ids: data,
      },
    })
    .then((res) => {
      return res.data;
    });
}

export function DeleteUploadById(data: number) {
  return axios.delete<Base>(`/api/v1/upload?id=${data}`, {}).then((res) => {
    return res.data;
  });
}

export function DeleteUploadByIdList(data: number) {
  return axios
    .delete<Base>(`/api/v1/upload/list`, {
      data: {
        userIds: data,
      },
    })
    .then((res) => {
      return res.data;
    });
}

export function GetUploadList(data: UploadReqData) {
  const uri = Object.entries(data)
    .map(([key, value]) => {
      if (value !== null && value !== undefined && value !== '') {
        return `${encodeURIComponent(key)}=${encodeURIComponent(
          value.toString()
        )}`;
      }
      return '';
    })
    .filter((item) => item !== '')
    .join('&');

  return axios.get<UploadRes>(`/api/v1/upload/list?${uri}`).then((res) => {
    return res.data;
  });
}

export function uploadFile(
  data: FormData,
  config: {
    controller: AbortController;
    onUploadProgress?: (progressEvent: any) => void;
  }
) {
  // const controller = new AbortController();
  return axios.post('/api/v1/upload', data);
}
