<script lang="ts">
  import { tick } from 'svelte'
  import Hint from './controls/Hint.svelte'
  import Label from './controls/Label.svelte'
  import TextArea from './controls/TextArea.svelte'
  import { caption, image, prevImage } from '../stores'

  let textarea: TextArea
</script>

<Label id="caption" label="Caption" />
<TextArea
  id="caption"
  bind:this={textarea}
  bind:value={$caption}
  dirty={$caption !== $image.Caption}
/>
{#if !$caption && $prevImage && $prevImage.Caption}
  <Hint
    on:click={() => {
      $caption = $prevImage.Caption
      tick().then(textarea.focus)
    }}>{$prevImage.Caption}</Hint
  >
{/if}
