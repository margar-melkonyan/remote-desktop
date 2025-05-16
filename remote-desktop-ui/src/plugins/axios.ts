import axios from "axios";
import { toast } from 'vue3-toastify';

axios.interceptors.request.use(config => {

  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
}, error => {
  return Promise.reject(error);
});

axios.interceptors.response.use(response => {
  return response;
}, error => {
  if (error.response?.status === 401) {
    localStorage.removeItem('token');
  }

  if (error.response.status !== 422) {
    toast.error(error.response.data.message)
  }

  return Promise.reject(error);
});


export default { axios }
