<template>
  <v-container
    max-width="800"
    class="mt-4"
  >
    <v-row
      class="mt-4"
      no-gutters
    >
      <div class="text-h6">
        {{ $t('sessions.title') }}
      </div>
    </v-row>
    <v-divider class="my-8" />
    <v-row class="d-flex justify-center">
      <v-tabs
        v-model="currentTab"
        :disabled="sessions.length === 0"
        fixed-tabs
      >
        <v-tab
          value="all"
        >
          {{ $t('sessions.all') }}
        </v-tab>
        <v-tab
          value="ssh"
        >
          {{ $t('sessions.ssh') }}
        </v-tab>
        <v-tab
          value="rdp"
        >
          {{ $t('sessions.rdp') }}
        </v-tab>
      </v-tabs>
    </v-row>
    <v-row class="mt-10">
      <v-col v-if="sessions.length > 0">
        <SessionCard
          v-for="(session, key) in sessions"
          :key="`session-card-${key}`"
          :session="session"
        />
      </v-col>
      <v-col
        v-else
        cols="12"
      >
        <v-row class="d-flex justify-center">
          <span class="text-h6">
            {{ $t('sessions.no_connections', [currentTab.toUpperCase()]) }}
          </span>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import axios from "axios";
const currentTab = ref<string>('all');
const sessions = [];
const selectFetchProtocols = () => {
  axios.get(`http://localhost:8080/api/v1/protocols?protocol=${currentTab.value}`)
}

const openConnection = () => {

}

watch(currentTab, () => {
  selectFetchProtocols()
});
</script>
