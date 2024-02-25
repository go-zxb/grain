import axios from 'axios';
import { Base } from '@/api/base';
import { UserState } from '@/store/modules/user/types';
import moment from "moment";

export interface LoginData {
  username: string;
  password: string;
}

export interface LoginRes extends Partial<Base> {
  data: {
    token: string;
  };
}
export function login(data: LoginData) {
  return axios.post<LoginRes>('/api/v1/sysUser/login', data).then((res) => {
    return res.data;
  });
}

export function logout() {
  return axios.post<LoginRes>('/api/v1/sysUser/logout');
}

export function getUserInfo() {
  return axios.get('/api/v1/sysUser/info').then((res) => {
    return res.data;
  });
}

export function getMenuList() {
  return axios.get('/api/v1/sysMenu/userMenu').then((res) => {
    return res.data;
  });
}

export function switchRole(data: any) {
  return axios
      .post<LoginRes>(`/api/v1/sysUser/switchRole?role=${data}`)
      .then((res) => {
        return res.data;
      });
}