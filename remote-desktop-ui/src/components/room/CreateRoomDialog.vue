<template>
  <v-dialog
    v-model="dialogCreateRoom"
    origin="left center"
    max-width="500"
    persistent
  >
    <v-card>
      <v-card-title class="mx-2 my-2">
        {{ $t('room.create') }}
      </v-card-title>
      <v-divider />
      <v-card-text>
        <v-text-field
          v-model="form.name"
          :error-messages="form.errors.get('name')"
          variant="outlined"
        >
          <template #label>
            {{ $t('room.fields.room_title') }}
          </template>
        </v-text-field>
        <v-row class="mb-2">
          <v-col class="py-0 d-flex align-center">
            {{ $t('room.fields.privacy') }}
          </v-col>
          <v-col class="py-0 d-flex justify-end">
            <v-switch
              v-model="form.is_private"
              :error-messages="form.errors.get('is_private')"
              inset
              base-color="red"
              color="green"
              hide-details
              @click="changeRoomPrivacy"
            />
          </v-col>
        </v-row>
        <v-text-field
          v-if="form.is_private"
          v-model="form.password"
          variant="outlined"
          :error-messages="form.errors.get('password')"
          :append-inner-icon="isHiddePassword ? 'mdi-eye' : 'mdi-eye-off'"
          :type="isHiddePassword ? 'password' : 'text'"
          @click:append-inner="showPassword"
        >
          <template #label>
            {{ $t('room.fields.password') }}
          </template>
        </v-text-field>
      </v-card-text>
      <v-divider />
      <v-card-actions>
        <v-col class="d-flex justify-start py-0">
          <v-btn
            @click="closeRoomCreationDialog"
          >
            {{ $t('close') }}
          </v-btn>
        </v-col>
        <v-col class="d-flex justify-end py-0">
          <v-btn
            :loading="form.busy"
            @click="newRoom"
          >
            {{ $t('create') }}
          </v-btn>
        </v-col>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { ref, defineEmits } from 'vue';
import Form from 'vform';
const { proxy } = getCurrentInstance();
const dialogCreateRoom = ref(false)
const isPrivateRoom = ref(false)
const isHiddePassword = ref(true);
const form = ref(new Form({
  name: '',
  is_private: false,
  password: '',
}));
const emit = defineEmits([
  "closeRoomCreateDialog",
])
function closeRoomCreationDialog() {
  emit('closeRoomCreateDialog')
  isPrivateRoom.value = false
  dialogCreateRoom.value = false
  form.value.reset()
}
function changeRoomPrivacy() {
  isPrivateRoom.value = !isPrivateRoom.value
}
function newRoom() {
  form.value.post(proxy.$api.rooms.urls.rooms())
    .then(() => {
      closeRoomCreationDialog()
      form.value.reset()
    })
    .catch(() => {})
}
function createRoom() {
  dialogCreateRoom.value = true
}
const showPassword = () => {
  isHiddePassword.value = !isHiddePassword.value
}
defineExpose({ createRoom })
</script>
