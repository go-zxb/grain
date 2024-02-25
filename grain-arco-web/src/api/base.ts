export interface Base {
  total: number;
  page: number;
  pageSize: number;
  success: boolean;
  code: number;
  message: string;
  time: number;
}

export interface Req {
  page: number;
  pageSize: number;
}
