<script setup lang="ts">
import '@wangeditor/editor/dist/css/style.css';

import { onBeforeUnmount, shallowRef } from 'vue';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';

defineOptions({
  name: 'RichEditor'
});

const model = defineModel('value', { type: String });

const props = defineProps<{ placeholder: string }>();

const editorRef = shallowRef();

const mode = 'simple';

const toolbarConfig = {};
const editorConfig = { placeholder: props.placeholder };

onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor === null) {
    return;
  }
  editor.destroy();
});

const handleCreated = (editor: any) => {
  editorRef.value = editor;
};
</script>

<template>
  <div class="div-border">
    <Toolbar class="tool-bar-border" :editor="editorRef" :default-config="toolbarConfig" :mode="mode" />
    <Editor
      v-model="model"
      class="rich-editor"
      :default-config="editorConfig"
      :mode="mode"
      @on-created="handleCreated"
    />
  </div>
</template>

<style lang="css">
.div-border {
  border: 1px solid #ccc;
}
.tool-bar-border {
  border-bottom: 1px solid #ccc;
}
.rich-editor {
  height: 200px !important;
  overflow-y: hidden;
}
</style>
