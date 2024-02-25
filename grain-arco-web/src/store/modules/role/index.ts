import { defineStore } from 'pinia';
import { GetRoleList, SysRole } from '@/api/sysManage/sysRole';
import { ref } from 'vue';
import { CascaderOption } from '@arco-design/web-vue';
import { Role } from './types';

const useRolesStore = defineStore('role', {
  state: (): Role => ({
    data: [],
  }),

  getters: {
    roles(state: Role): Role {
      return { ...state };
    },
  },

  actions: {
    // setRole's information
    setRoles(partial: Partial<SysRole[]>) {
      const roles = ref<CascaderOption[]>([]);
      partial.forEach((v) => {
        roles.value.push({
          value: v?.role,
          label: v?.roleName,
        });
      });
      this.$patch({ data: roles.value });
    },

    // Reset sysUser's information
    resetInfo() {
      this.$reset();
    },

    // GetRole's information
    async getRoles() {
      if (this.data.length === 0) {
        const res = await GetRoleList({ page: 1, pageSize: 20 });
        this.setRoles(res.data);
      }
    },
  },
});

export default useRolesStore;
