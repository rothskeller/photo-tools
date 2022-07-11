<script lang="ts">
  import { tick } from 'svelte'
  import Hint from './controls/Hint.svelte'
  import Label from './controls/Label.svelte'
  import TextInput from './controls/TextInput.svelte'
  import { artist, image, images } from '../stores'

  let textinput: TextInput
  let allArtists: string[]
  $: {
    const artistSet = new Set<string>()
    Object.values($images).forEach((image) => {
      if (image.Artist) artistSet.add(image.Artist)
    })
    artistSet.add('Steven Roth')
    allArtists = [...artistSet.keys()]
    allArtists.sort()
  }
</script>

<Label id="artist" label="Artist" />
<TextInput
  id="artist"
  bind:this={textinput}
  bind:value={$artist}
  dirty={$artist !== $image.Artist}
/>
{#if !$artist}
  {#each allArtists as a}
    <Hint
      on:click={() => {
        $artist = a
        tick().then(textinput.focus)
      }}>{a}</Hint
    >
  {/each}
{/if}
