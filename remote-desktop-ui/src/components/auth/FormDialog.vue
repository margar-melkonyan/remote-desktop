<template>
  <v-dialog
    v-model="props.loginDialog"
    max-width="500"
    persistent
    scrollable
  >
    <v-card>
      <v-tabs
        v-model="tab"
        fixed-tabs
      >
        <v-tab
          value="sign-in"
          @click="changeTab('sign-in')"
        >
          {{ $t('enterForm.sign_in') }}
        </v-tab>
        <v-tab
          value="sign-up"
          @click="changeTab('sign-up')"
        >
          {{ $t('enterForm.sign_up') }}
        </v-tab>
      </v-tabs>
      <v-divider />
      <div class="scrollable-content">
        <v-tabs-window v-model="tab">
          <v-tabs-window-item value="sign-in">
            <SignInForm @close="closeFormDialog" />
          </v-tabs-window-item>
          <v-tabs-window-item value="sign-up">
            <SignUpForm @close="changeTab('sign-in')" />
          </v-tabs-window-item>
        </v-tabs-window>
      </div>
      <v-divider />
      <v-card-actions>
        <v-col class="d-flex justify-start py-0">
          <v-btn
            block
            @click="closeFormDialog"
          >
            {{ $t('close') }}
          </v-btn>
        </v-col>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import {defineProps, defineEmits, ref} from 'vue';
const tab = ref("sign-in")
const emit = defineEmits([
  "close"
]);
const props = defineProps<{
  loginDialog: boolean;
}>();
const changeTab = (tabVal: string) => {
  tab.value = tabVal
}
const closeFormDialog = () => {
  emit("close");
}
</script>
<style>
.scrollable-content {
  max-height: 500px;
  overflow-y: auto;
}
</style>
