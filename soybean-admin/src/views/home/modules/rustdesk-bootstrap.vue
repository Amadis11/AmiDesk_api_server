<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { fetchRustdeskBootstrap } from '@/service/api/system';

defineOptions({
  name: 'RustdeskBootstrap'
});

const bootstrap = ref<Api.System.RustdeskBootstrap>();

onMounted(async () => {
  const response = await fetchRustdeskBootstrap();
  if (response.data !== null) {
    bootstrap.value = response.data;
  }
});
</script>

<template>
  <NCard title="RustDesk Server / Bootstrap" :bordered="false" size="small" class="card-wrapper">
    <NDescriptions v-if="bootstrap" label-placement="left" :column="2" bordered size="small">
      <NDescriptionsItem label="Bootstrap">
        <NTag :type="bootstrap.enabled ? 'success' : 'default'" size="small">
          {{ bootstrap.enabled ? 'Enabled' : 'Disabled' }}
        </NTag>
      </NDescriptionsItem>
      <NDescriptionsItem label="Public Key status">
        <NTag :type="bootstrap.key_loaded ? 'success' : 'warning'" size="small">
          {{ bootstrap.key_status }}
        </NTag>
      </NDescriptionsItem>
      <NDescriptionsItem label="ID Server">{{ bootstrap.id_server || '-' }}</NDescriptionsItem>
      <NDescriptionsItem label="Relay Server">{{ bootstrap.relay_server || '-' }}</NDescriptionsItem>
      <NDescriptionsItem label="API Server">{{ bootstrap.api_server || '-' }}</NDescriptionsItem>
      <NDescriptionsItem label="Public Key file">{{ bootstrap.key_file || '-' }}</NDescriptionsItem>
      <NDescriptionsItem label="Public Key fingerprint">{{ bootstrap.key_fingerprint || '-' }}</NDescriptionsItem>
      <NDescriptionsItem label="Bootstrap revision">{{ bootstrap.revision }}</NDescriptionsItem>
      <NDescriptionsItem label="Push to anonymous">{{ bootstrap.push_to_anonymous ? 'Yes' : 'No' }}</NDescriptionsItem>
      <NDescriptionsItem label="Push to authenticated">{{ bootstrap.push_to_authenticated ? 'Yes' : 'No' }}</NDescriptionsItem>
    </NDescriptions>
    <NEmpty v-else description="Bootstrap status unavailable" size="small" />
  </NCard>
</template>