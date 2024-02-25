import axios from 'axios';

export interface ContentDataRecord {
  x: string;
  y: number;
}

export function queryContentData() {
  return axios.get('/api/content-data').then((res) => {
    return res.data;
  });
}

export interface PopularRecord {
  key: number;
  clickNumber: string;
  title: string;
  increases: number;
}

export function queryPopularList(params: { type: string }) {
  return axios.get('/api/popular/list', { params }).then((res) => {
    return res.data;
  });
}
