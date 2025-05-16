<template>
  <v-card
    color="#7fff94"
    variant="outlined"
  >
    <v-card-title>
      {{ $t('room.title', [props.room.name]) }}
    </v-card-title>
    <v-card-text>
      {{ $t('room.players', [props.room.player_in, props.room.capacity]) }}
    </v-card-text>
    <v-card-actions>
      <v-row>
        <v-col
          cols="12"
          class="d-flex
          justify-end"
        >
          <v-btn
            variant="outlined"
            :color="props.room.player_in === props.room.capacity ? 'white' : '#ff7fea'"
            density="comfortable"
            @click="openRoom"
          >
            {{ $t('enter') }}
          </v-btn>
        </v-col>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts" setup>
import { ref, defineEmits } from 'vue';
import { useAuthStore } from '@/stores/auth';
const router = useRouter()
const emit = defineEmits([
  "openLoginDialog"
])
const authStore = useAuthStore()
const props = defineProps<{
  room: {
    id: string;
    name: string;
    player_in: number;
    capacity: number;
    is_private: boolean;
    created_at: string;
  },
}>()
const enterRoom = ref(false);
function openRoom() {
  if (authStore.user == null) {
    emit("openLoginDialog")
    return
  }
  router.push({
    name: 'rooms.game',
    params: { id: props.room.id, room: props.room},
  })
  enterRoom.value = true
}
function closeRoom() {
  enterRoom.value = false
}
</script>
