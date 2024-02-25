import moment from 'moment';

export default function DateString(v: any) {
  return moment(v).format('YYYY-MM-DD HH:mm:ss').toString();
}
