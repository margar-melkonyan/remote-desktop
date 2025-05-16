// Utilities
import axios from 'axios';
import { defineStore } from 'pinia'
import Form from 'vform';
import { getCurrentInstance } from 'vue';

interface User {
  id: number;
  name: string;
  email: string;
  current_won_score?: number;
  created_at?: string;
}

export const useAuthStore = defineStore('auth', () => {
  const { proxy } = getCurrentInstance();
  const apiAuth = proxy.$api.auth;
  const apiUsers = proxy.$api.users;
  const user= ref<User|null>(null);

  async function signIn(form: Form) {
    const { data, status } = await form.post(apiAuth.urls.signIn())
    localStorage.setItem("token", data.data.token)
    await currentUser()
    return status
  }

  async function signUp(form: Form) {
    const { data, status } = await form.post(apiAuth.urls.signUp())
    return { data, status }
  }

  async function currentUser() {
    axios.get(apiUsers.urls.current())
      .then(({ data }) => {
        user.value = data.data
      })
  }

  function signOut() {
    user.value = null
    localStorage.removeItem('token')
  }

  return { user, currentUser, signIn, signUp, signOut }
})
