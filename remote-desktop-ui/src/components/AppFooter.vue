<template>
  <v-footer
    height="40"
    app
    class="d-flex justify-center"
  >
    <v-container
      max-width="800"
      class="d-flex align-center"
    >
      <v-row
        class="text-caption text-disabled justify-end"
      >
        &copy;
        <span>
          RemoteDesktop GateWay By
          <a
            target="_blank"
            href="https://github.com/margar-melkonyan"
            class="social-link"
          >
            Melkonyan M.
          </a>
        </span>
      </v-row>
    </v-container>
  </v-footer>
</template>

<script setup lang="ts">
  //
</script>

<style scoped lang="sass">
  .social-link
    color: rgba(var(--v-theme-on-background), var(--v-disabled-opacity))
    text-decoration: underline #ff7fea
    transition: .2s ease-in-out

    &:hover
      color: #7fff94
</style>
