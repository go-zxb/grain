import debug from './env';

export default ({ mock, setup }: { mock?: boolean; setup: () => void }) => {
  if (mock !== false && debug) setup();
};

export const successResponseWrap = (data: unknown) => {
  return {
    data,
    success: true,
    message: '请求成功',
    code: 2000,
  };
};

export const failResponseWrap = (
  data: unknown,
  message: string,
  code = 5000
) => {
  return {
    data,
    success: false,
    message,
    code,
  };
};
