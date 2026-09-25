<script setup lang="ts">
import { computed } from 'vue';

defineOptions({
  name: 'ChangeLogs'
});

interface LogItem {
  id: number;
  content: string;
  version: string;
  time: string;
}

const logs = computed<LogItem[]>(() => [
  {
    id: 1,
    content: 'AmiDesk v2.1: project rebranding across the Go module, CLI, Docker image and administration panel.',
    version: 'v2.1',
    time: '2026-09-25'
  },
  {
    id: 2,
    content: 'Amitronic v2.0: anonymous HTTPS bootstrap, admin TOTP hardening, account lockout and Docker deployment improvements.',
    version: 'v2.0',
    time: '2026-08-28'
  },
  { id: 3, content: 'RustDesk client compatibility update.', version: 'v1.2.7', time: '2024-08-29' }
]);
</script>

<template>
  <NCard :title="$t('page.home.changeLogs')" :bordered="false" size="small" segmented class="card-wrapper">
    <NList>
      <NListItem v-for="item in logs" :key="item.id">
        <template #prefix>
          <SoybeanAvatar class="size-48px!" />
        </template>
        <NThing :title="item.content" :title-extra="item.version" :description="item.time" />
      </NListItem>
    </NList>
  </NCard>
</template>

<style scoped></style>
