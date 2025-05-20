<template>
  <v-card
    class="pa-2"
    rounded="xl"
  >
    <v-card-title>
      {{ $t('connections.edit') }}
    </v-card-title>
    <v-divider class="mt-2" />
    <v-card-text>
      <v-text-field
        v-model="form.name"
        class="my-4"
        hide-details="auto"
        :label="$t('new_connections.name')"
        :variant="'outlined'"
        :error-messages="form.errors.get('name')"
      />
      <v-select
        v-model="form.protocol"
        class="my-4"
        hide-details="auto"
        :items="connectionItems"
        item-value="value"
        item-title="name"
        :label="$t('new_connections.protocol')"
        :variant="'outlined'"
        :error-messages="form.errors.get('protocol')"
      />
      <v-divider class="my-4" />
      <v-text-field
        v-model="form.username"
        class="my-4"
        hide-details="auto"
        :label="$t('new_connections.username')"
        :variant="'outlined'"
        :error-messages="form.errors.get('username')"
      />
      <v-text-field
        v-model="form.password"
        class="my-4"
        hide-details="auto"
        :label="$t('new_connections.password')"
        :variant="'outlined'"
        :error-messages="form.errors.get('password')"
      />
      <v-text-field
        v-model="form.host_name"
        class="my-4"
        hide-details="auto"
        :label="$t('new_connections.hostname')"
        :variant="'outlined'"
        :error-messages="form.errors.get('host_name')"
      />
      <v-text-field
        v-model="form.port"
        class="my-4"
        hide-details="auto"
        :label="$t('new_connections.port')"
        :variant="'outlined'"
        :error-messages="form.errors.get('port')"
      />
    </v-card-text>
    <v-divider class="mb-2" />
    <v-card-actions>
      <v-row
        class="d-flex justify-space-between px-3"
        no-gutters
      >
        <v-btn
          color="primary"
          variant="tonal"
          rounded="xl"
          @click="emit('close')"
        >
          {{ $t('close') }}
        </v-btn>
        <v-btn
          color="green"
          variant="tonal"
          rounded="xl"
          @click=""
        >
          {{ $t('update') }}
        </v-btn>
      </v-row>
    </v-card-actions>
  </v-card>
</template>

<script setup lang="ts">
import {connectionItems} from "@/plugins/common";
import {Form} from "vform";

const props = defineProps({
  connectionInfo: {
    host_name: String,
    identifier: String,
    name: String,
    password: String,
    port: String,
    protocol: String,
    username: String,
  },
})

const emit = defineEmits([
  'close'
])

const form = ref(
  new Form({
    host_name: "",
    name: "",
    protocol: "ssh",
    username: "",
    password: "",
    port: "",
  }),
);

onMounted(() => {
  form.value.fill(props.connectionInfo)
})
</script>


<style scoped>

</style>
