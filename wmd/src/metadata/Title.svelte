<script lang="ts">
  import { tick } from 'svelte'
  import Hint from './controls/Hint.svelte'
  import Label from './controls/Label.svelte'
  import TextInput from './controls/TextInput.svelte'
  import { title, image, prevImage } from '../stores'

  let textinput: TextInput
</script>

<Label id="title" label="Title" />
<TextInput id="title" bind:this={textinput} bind:value={$title} dirty={$title !== $image.Title} />
{#if !$title && $prevImage && $prevImage.Title}
  <Hint
    on:click={() => {
      $title = $prevImage.Title
      tick().then(textinput.focus)
    }}>{$prevImage.Title}</Hint
  >
{/if}
