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
        {{ $t('connections.title') }}
      </div>
    </v-row>
    <v-divider class="my-8" />
    <v-row class="d-flex justify-center">
      <v-tabs
        v-model="currentTab"
        fixed-tabs
      >
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
          :key="`session-card-${key}`"
          :session="session"
          @update-connections="selectFetchProtocols"
          @edit-connection="openEditDialog"
        />
      </v-col>
      <v-col
        v-else
        cols="12"
      >
        <v-row class="d-flex justify-center">
          <span class="text-h6">
            {{ $t('connections.no', [currentTab.toUpperCase()]) }}
          </span>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
  <v-dialog
    v-model="editDialog"
    max-width="500"
    persistent
    scrollable
  >
    <edit-form
      :connection-info="connectionInfo"
      @close="editDialog = false"
    />
  </v-dialog>
</template>

<script lang="ts" setup>
import axios from "axios";
import { getCurrentInstance } from "vue";
const { proxy } = getCurrentInstance();
const apiSessions = proxy.$api.sessions;

const currentTab = ref<string>('all');
const sessions = ref([]);
const editDialog = ref(false);
const connectionInfo = ref({});

const selectFetchProtocols = () => {
  axios.get(`${apiSessions.urls.index()}?protocol=${currentTab.value}`, {
    headers: {
      "Guacamole-Token": window.localStorage.getItem("guac_token"),
    }
  })
    .then(({ data }) => {
      sessions.value = data.data
    })
}

const openEditDialog = async (id: number) => {
  const { data } = await  axios.get(apiSessions.urls.edit(id), {
    headers: {
      "Guacamole-Token": localStorage.getItem("guac_token"),
    },
  });
  connectionInfo.value = data.data;
  editDialog.value = true;
}

watch(currentTab, () => {
  selectFetchProtocols()
});
onMounted(() => {
  selectFetchProtocols()
})
</script>
