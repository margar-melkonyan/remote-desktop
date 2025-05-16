<template>
  <v-container>
    <v-row class="pa-0 mt-6">
      <v-col cols="6">
        <div class="text-h5">
          {{ $t('rooms.title') }}
        </div>
      </v-col>
      <v-col
        v-if="auth.user !== null"
        class="d-flex justify-end"
        cols="6"
      >
        <v-btn @click="openCreateRoomDialog">
          {{ $t('rooms.create') }}
        </v-btn>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-tabs
          v-if="auth.user !== null"
          color="#ff7fea"
          fixed-tabs
        >
          <v-tab
            fixed
            @click="changeTab('all')"
          >
            {{ $t('rooms.all') }}
          </v-tab>
          <v-tab
            fixed
            @click="changeTab('my')"
          >
            {{ $t('rooms.my') }}
          </v-tab>
        </v-tabs>
      </v-col>
    </v-row>
    <v-divider class="my-8" />
    <RoomList
      class="my-8"
      :rooms="rooms"
      @open-login-dialog="openLoginDialog"
    />
    <v-divider class="my-4" />
    <CreateRoomDialog
      ref="newRoomDialog"
      @close-room-create-dialog="fetchRooms"
    />
  </v-container>
</template>

<script lang="ts" setup>
import type CreateRoomDialog from "@/components/room/CreateRoomDialog.vue";
import {ref} from "vue";
import {useAuthStore} from "@/stores/auth";
import axios from "axios";
const {proxy} = getCurrentInstance();

const auth = useAuthStore();
const apiRooms = proxy.$api.rooms;
const newRoomDialog = ref<InstanceType<typeof CreateRoomDialog> | null>(null);
const rooms = ref([]);
const currentTab = ref<string>("all")
let intervalId: number;
const emit = defineEmits([
  'openLoginDialog'
])
const openLoginDialog = () => {
  emit('openLoginDialog')
}
const openCreateRoomDialog = () => {
  if (newRoomDialog.value) {
    newRoomDialog.value.createRoom();
  }
};
const fetchRooms = () => {
  let url;
  switch (currentTab.value) {
    case "all":
      url = apiRooms.urls.rooms();
      break;
    case "my":
      url = apiRooms.urls.my();
      break;
  }
  axios.get(url)
    .then(({data}) => {
      rooms.value = data.data ?? []
    })
}
const changeTab = (tab: string) => {
  currentTab.value = tab
  fetchRooms()
}
onMounted(() => {
  intervalId = setInterval(fetchRooms, 2000)
  fetchRooms()
})
onBeforeUnmount(() => {
  clearInterval(intervalId)
})
</script>
