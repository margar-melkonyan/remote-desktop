<template>
  <v-app>
    <v-layout class="d-flex justify-center">
      <v-app-bar>
        <v-app-bar-title>
          <div
            class="animated-gradient-text"
            style="padding-right: 25rem;"
          >
            {{ $t('app.title') }}
          </div>
        </v-app-bar-title>
        <v-container
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
      <v-navigation-drawer
        v-if="auth?.user"
      >
        <v-col
          class="pa-0 my-5"
        >
          <v-list
            density="compact"
            nav
          >
            <v-list-item
              v-for="(item,key) in drawerItems"
              :key="`drawer-item-${key}`"
              class="my-2"
              :title="item.title"
              :prepend-icon="item.icon"
              :value="item.routeName"
              @click="moveTo(item.routeName)"
            />
          </v-list>
        </v-col>
      </v-navigation-drawer>
      <v-main>
        <v-container>
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
</template>

<script lang="ts" setup>
import AppFooter from './components/AppFooter.vue';
import { ref } from "vue";
import { useAuthStore } from "@/stores/auth";
import { drawerItems } from "@/plugins/drawerItems";

const router = useRouter();
const auth = useAuthStore();
const loginDialog = ref(false)
auth.currentUser()
const openLoginDialog = () => {
  loginDialog.value = true;
}
const signOut = () => {
  auth.signOut()
  router.push({name: "index"})
}
const moveTo = (routeName) => {
  router.push({
    name: routeName,
  })
}
</script>

<style scoped>
.animated-gradient-text {
  background: linear-gradient(
    90deg,
    #e8f5e9,
    #c8e6c9,
    #a5d6a7,
    #81c784,
    #66bb6a,
    #4caf50,
    #43a047,
    #388e3c,
    #2e7d32,
    #1b5e20,
    #2e7d32,
    #388e3c,
    #43a047,
    #4caf50,
    #66bb6a,
    #81c784,
    #a5d6a7,
    #c8e6c9
  );
  background-size: 400% auto;
  color: transparent;
  -webkit-background-clip: text;
  background-clip: text;
  animation: gradient 16s ease infinite;
}

@keyframes gradient {
  0% {
    background-position: 0% center;
  }
  100% {
    background-position: 400% center;
  }
}
</style>
