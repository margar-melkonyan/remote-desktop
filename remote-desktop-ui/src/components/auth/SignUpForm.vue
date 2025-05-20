<template>
  <v-card-text>
    <v-col>
      <v-text-field
        v-model="form.name"
        hide-details="auto"
        variant="outlined"
        :label="$t('enterForm.name')"
        :error-messages="form.errors.get('name')"
        clearable
      />
    </v-col>
    <v-col>
      <v-text-field
        v-model="form.email"
        hide-details="auto"
        variant="outlined"
        :label="$t('enterForm.email')"
        :error-messages="form.errors.get('email')"
        clearable
      />
    </v-col>
    <v-col>
      <v-text-field
        v-model="form.password"
        hide-details="auto"
        variant="outlined"
        :label="$t('enterForm.password')"
        :error-messages="form.errors.get('password')"
        clearable
        :append-inner-icon="isHiddePassword ? 'mdi-eye' : 'mdi-eye-off'"
        :type="isHiddePassword ? 'password' : 'text'"
        @click:append-inner="showPassword"
      />
    </v-col>
    <v-col>
      <v-text-field
        v-model="form.password_confirmation"
        hide-details="auto"
        variant="outlined"
        :label="$t('enterForm.password_confirmation')"
        :error-messages="form.errors.get('password_confirmation')"
        clearable
        :append-inner-icon="isHiddePassword ? 'mdi-eye' : 'mdi-eye-off'"
        :type="isHiddePassword ? 'password' : 'text'"
        @click:append-inner="showPassword"
      />
    </v-col>
  </v-card-text>
  <v-card-actions>
    <v-btn
      variant="tonal"
      block
      @click="signUp"
    >
      {{ $t('enterForm.sign_up') }}
    </v-btn>
  </v-card-actions>
</template>

<script lang="ts" setup>
import { useAuthStore } from '@/stores/auth';
import Form from 'vform';
import { defineEmits } from 'vue';
import { toast } from 'vue3-toastify';
const authStore = useAuthStore();
const isHiddePassword = ref(true);
const emit = defineEmits([
  "close",
])
const form = ref(new Form({
  name: '',
  email: '',
  password: '',
  password_confirmation: null,
}));
const signUp = () => {
  authStore.signUp(form.value)
    .then((data) => {
      if (data.status === 200) {
        toast.success(data.data?.message)
        closeFormDialog()
      }
    })
}
const closeFormDialog = () => {
  emit("close");
}
const showPassword = () => {
  isHiddePassword.value = !isHiddePassword.value
}
</script>
