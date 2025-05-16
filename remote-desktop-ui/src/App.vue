<template>
  <v-app>
    <v-layout class="d-flex justify-center">
      <v-app-bar>
        <v-container
          max-width="800"
          class="justify-center"
        >
          <v-col
            v-if="auth.user == null"
            class="d-flex justify-end"
          >
            <v-btn
              @click="openLoginDialog"
            >
              {{ $t('menu.sign_in') }}
            </v-btn>
          </v-col>
          <v-col
            v-else
            class="d-flex justify-end"
            cols="12"
          >
            <v-btn
              icon
            >
              <v-icon>
                mdi-account-circle
              </v-icon>
              <v-menu
                location="left"
                origin="top"
                transition="scale-transition"
                activator="parent"
              >
                <v-list style="border: #ff7fea 0.05rem solid">
                  <v-list-item style="max-width: 400px">
                    <v-list-item-title>
                      {{ auth.user?.name }} / {{ auth.user?.email }}
                    </v-list-item-title>
                    <v-list-item-subtitle class="d-flex justify-end my-2">
                      {{ $t('navbar.auth.won', [auth.user?.current_won_score]) }}
                    </v-list-item-subtitle>
                  </v-list-item>
                  <v-list-item>
                    <v-list-item-title>
                      <v-tooltip
                        text="Статистика за последние 50 игр"
                        location="bottom"
                      >
                        <template #activator="{ props }">
                          <v-btn
                            v-bind="props"
                            block
                            color="green"
                            variant="tonal"
                            @click="openStatistic"
                          >
                            {{ $t('statistics.open') }}
                          </v-btn>
                        </template>
                      </v-tooltip>
                    </v-list-item-title>
                  </v-list-item>
                  <v-divider class="my-4" />
                  <v-list-item>
                    <v-btn
                      color="#ff7fea"
                      block
                      prepend-icon="mdi-exit-to-app"
                      variant="tonal"
                      @click="signOut"
                    >
                      {{ $t('menu.sign_out') }}
                    </v-btn>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-btn>
          </v-col>
        </v-container>
      </v-app-bar>
      <v-main>
        <v-container max-width="800">
          <router-view @open-login-dialog="openLoginDialog" />
        </v-container>
      </v-main>
      <AppFooter />
    </v-layout>
  </v-app>
  <FormDialog
    :login-dialog="loginDialog"
    @close="loginDialog = false"
  />
  <v-dialog
    v-model="statisticDialog"
    origin="left center"
    max-width="500"
    max-height="500"
    scrollable
    persistent
  >
    <UserStatistic
      @close-dialog="statisticDialog = false"
    />
  </v-dialog>
</template>

<script lang="ts" setup>
import AppFooter from './components/AppFooter.vue';
import { ref } from "vue";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();
const loginDialog = ref(false)
const statisticDialog = ref(false);
auth.currentUser()
const openLoginDialog = () => {
  loginDialog.value = true;
}
const openStatistic = () => {
  statisticDialog.value = true
}
const signOut = () => {
  auth.signOut()
  router.push({name: "index"})
}
</script>
