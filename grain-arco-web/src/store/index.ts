import { createPinia } from 'pinia';
import useAppStore from './modules/app';
import useUserStore from './modules/user';
import useRolesStore from './modules/role';
import useTabBarStore from './modules/tab-bar';

const pinia = createPinia();

export { useAppStore, useUserStore, useTabBarStore, useRolesStore };
export default pinia;
