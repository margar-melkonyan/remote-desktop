<template>
  <v-container max-width="800">
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
            clearable
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
            clearable
            :variant="'outlined'"
            :error-messages="form.errors.get('username')"
          />
          <v-text-field
            v-model="form.password"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.password')"
            clearable
            :variant="'outlined'"
            :error-messages="form.errors.get('password')"
            :append-inner-icon="isHiddePassword ? 'mdi-eye' : 'mdi-eye-off'"
            :type="isHiddePassword ? 'password' : 'text'"
            @click:append-inner="showPassword"
          />
          <v-text-field
            v-model="form.host_name"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.hostname')"
            clearable
            :variant="'outlined'"
            :error-messages="form.errors.get('host_name')"
          />
          <v-text-field
            v-model="form.port"
            class="my-4"
            hide-details="auto"
            :label="$t('new_connections.port')"
            clearable
            :variant="'outlined'"
            :error-messages="form.errors.get('port')"
          />
        </v-card-text>
        <v-card-actions class="my-2 mx-2">
          <v-btn
            type="submit"
            color="blue"
            variant="flat"
            rounded="lg"
            block
            :loading="form.busy"
          >
            {{ $t("new_connections.add_connection") }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-container>
</template>

<script setup lang="ts">
import { Form } from "vform";
import { getCurrentInstance } from 'vue';
import { connectionItems } from "@/plugins/common";
const { proxy } = getCurrentInstance();
const apiSessions = proxy.$api.sessions;
const router = useRouter();
const isHiddePassword = ref(true);
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

const showPassword = () => {
  isHiddePassword.value = !isHiddePassword.value
}

const submit = () => {
  form.value.post(apiSessions.urls.store(), {
    headers: {
      "Guacamole-Token": localStorage.getItem('guac_token')
    }
  }).then(() => {
    router.push({name: "connections"})
  })
};
</script>
