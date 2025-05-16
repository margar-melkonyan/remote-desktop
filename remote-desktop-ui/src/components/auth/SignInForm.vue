<template>
  <v-card-text>
    <v-col>
      <v-text-field
        v-model="form.email"
        hide-details="auto"
        variant="outlined"
        :error-messages="form.errors.get('email')"
        :label="$t('enterForm.email')"
      />
    </v-col>
    <v-col>
      <v-text-field
        v-model="form.password"
        hide-details="auto"
        variant="outlined"
        :label="$t('enterForm.password')"
        :error-messages="form.errors.get('password')"
        :append-inner-icon="isHiddePassword ? 'mdi-eye' : 'mdi-eye-off'"
        :type="isHiddePassword ? 'password' : 'text'"
        @click:append-inner="showPassword"
      />
    </v-col>
  </v-card-text>
  <v-card-actions>
    <v-btn
      :loading="form.busy"
      variant="tonal"
      block
      @click="signIn"
    >
      {{ $t('enterForm.sign_in') }}
    </v-btn>
  </v-card-actions>
</template>
<script lang="ts" setup>
import { ref, defineEmits } from 'vue';
import { Form } from 'vform';
import 'vue3-toastify/dist/index.css';
import { useAuthStore } from '@/stores/auth';
const isHiddePassword = ref(true);
const authStore = useAuthStore();
const form = ref(new Form({
  email: '',
  password: ''
}));
const emit = defineEmits([
  "close",
])
const signIn = () => {
  authStore.signIn(form.value)
    .then((value) => {
      if (value === 200) {
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
