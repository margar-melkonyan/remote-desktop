<template>
  <v-card
    class="my-4"
    rounded="xl"
  >
    <v-card-title class="mt-2">
      <v-row>
        <v-col cols="6">
          <v-tooltip
            :text="props.session.name"
            location="bottom left"
          >
            <template #activator="{ props: tooltipProps }">
              <span v-bind="tooltipProps">
                {{ props.session.name.length > 15 ? props.session.name.slice(0, 15) + '...' : props.session.name }}
              </span>
            </template>
          </v-tooltip>
        </v-col>
        <v-col cols="6">
          <div class="d-flex justify-end">
            <v-btn
              class="mr-4"
              icon="mdi-pencil"
              color="primary"
              size="small"
              @click="sendEdit"
            />
            <v-btn
              icon="mdi-trash-can"
              color="red"
              size="small"
              @click="openConfirmationDialog"
            />
          </div>
        </v-col>
      </v-row>
    </v-card-title>
    <v-divider class="my-2" />
    <v-card-actions class="mb-2">
      <v-row class="d-flex justify-space-between px-5">
        <div class="align-self-center">
          {{ $t("connections.protocol", [props.session.protocol]) }}
        </div>
        <v-btn
          color="blue"
          variant="tonal"
          rounded="lg"
          @click="openConnection"
        >
          {{ $t("connections.connect") }}
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
  <v-dialog
    v-model="confirmationDialog"
    max-width="500"
    persistent
    scrollable
  >
    <v-card
      rounded="xl"
      class="py-3"
    >
      <v-card-title>
        <span style="color: red;">
          {{ $t('connections.caution') }}
        </span>
      </v-card-title>
      <v-divider />
      <v-card-text>
        {{ $t('connections.cautions.delete') }}
      </v-card-text>
      <v-divider />
      <v-card-actions class="mt-3">
        <v-row
          class="d-flex justify-space-between px-3"
          no-gutters
        >
          <v-btn
            color="primary"
            variant="tonal"
            rounded="xl"
            @click="confirmationDialog = false"
          >
            {{ $t('close') }}
          </v-btn>
          <v-btn
            color="red"
            variant="tonal"
            rounded="xl"
            @click="confirmDeletion"
          >
            {{ $t('delete') }}
          </v-btn>
        </v-row>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import {useRouter} from "vue-router";
import {getCurrentInstance} from "vue";
import axios from "axios";

const {proxy} = getCurrentInstance();
const apiSessions = proxy.$api.sessions;
const confirmationDialog = ref(false);
const router = useRouter();

const props = defineProps({
  session: {
    identifier: String,
    connection: String,
    protocol: String,
  },
});

const emit = defineEmits([
  'updateConnections'
])

const openConnection = () => {
  router.push({
    name: "current-session",
    params: {id: props.session.identifier},
  });
};

const openConfirmationDialog = () => {
  confirmationDialog.value = true;
};
const confirmDeletion = () => {
  axios.delete(apiSessions.urls.delete(props.session.identifier), {
    headers: {
      "Guacamole-Token": localStorage.getItem("guac_token"),
    },
  })
    .then(() => {
      confirmationDialog.value = false;
      emit('updateConnections')
    });
}

const sendEdit = () => {
  axios.get(apiSessions.urls.edit(props.session.identifier), {
    headers: {
      "Guacamole-Token": localStorage.getItem("guac_token"),
    },
  })
}
</script>
