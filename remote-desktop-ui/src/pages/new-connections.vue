<template>
  <v-container
    max-width="800"
  >
    <v-form @submit.prevent="submit">
      <v-card
        class="ma-8"
        rounded="xl"
        :loading="form.busy"
      >
        <v-card-text>
          <v-text-field
            v-model="form.name"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.name')"
            :variant="'outlined'"
          />
          <v-select
            v-model="form.protocol"
            class="my-4"
            hide-details="auto"
            :items="items"
            item-value="value"
            item-title="name"
            :label="$t('new_connections.protocol')"
            :variant="'outlined'"
          />
          <v-divider class="my-4" />
          <v-text-field
            v-model="form.username"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.username')"
            :variant="'outlined'"
          />
          <v-text-field
            v-model="form.password"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.password')"
            :variant="'outlined'"
          />
          <v-text-field
            v-model="form.hostname"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.hostname')"
            :variant="'outlined'"
          />
          <v-text-field
            v-model="form.port"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.port')"
            :variant="'outlined'"
          />
        </v-card-text>
        <v-card-actions class="my-2 mx-2">
          <v-btn
            color="blue"
            variant="flat"
            rounded="lg"
            block
          >
            {{ $t('new_connections.add_connection') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-container>
</template>

<script setup lang="ts">

import {Form} from "vform";

const form = ref(new Form({
  name: '',
  protocol: '',
  username: '',
  password: '',
  hostname: '',
  port: '',
}))

const items = [
  {
    name: 'SSH',
    value: 'ssh'
  },
  {
    name: 'RDP',
    value: 'rdp'
  },
];

const submit = () => {
  console.log('')
}

</script>
