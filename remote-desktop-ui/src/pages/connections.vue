<template>
  <v-container max-width="800" class="mt-4">
    <v-row class="mt-4" no-gutters>
      <div class="text-h6">
        {{ $t('connections.title') }}
      </div>
    </v-row>
    <v-divider class="my-8" />
    <v-row class="d-flex justify-center">
      <v-tabs v-model="currentTab" fixed-tabs>
        <v-tab value="all">
          {{ $t('connections.all') }}
        </v-tab>
        <v-tab value="ssh">
          {{ $t('connections.ssh') }}
        </v-tab>
        <v-tab value="rdp">
          {{ $t('connections.rdp') }}
        </v-tab>
      </v-tabs>
    </v-row>
    <v-row class="mt-10">
      <v-col v-if="sessions.length > 0">
        <SessionCard
          v-for="(session, key) in sessions"
          class="my-4"
          :key="`session-card-${key}`"
          :session="session"
        />
      </v-col>
      <v-col v-else cols="12">
        <v-row class="d-flex justify-center">
          <span class="text-h6">
            {{ $t('connections.no', [currentTab.toUpperCase()]) }}
          </span>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import axios from "axios";
const currentTab = ref<string>('all');
const sessions = ref([]);
const selectFetchProtocols = () => {
  axios.get(`http://192.168.1.4:8000/api/v1/sessions?protocol=${currentTab.value}`, {
    headers: {
      "Guacamole-Token": window.localStorage.getItem("guac_token"),
    }
  })
    .then(({ data }) => {
      sessions.value = data.data
    })
}
watch(currentTab, () => {
  selectFetchProtocols()
});
onMounted(() => {
  selectFetchProtocols()
})
</script>
