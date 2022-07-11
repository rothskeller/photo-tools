<script lang="ts">
  import { tick } from 'svelte'
  import Hint from './controls/Hint.svelte'
  import Label from './controls/Label.svelte'
  import TextInput from './controls/TextInput.svelte'
  import { gps, image, prevImage } from '../stores'

  let textinput: TextInput
  let xlate: string
  $: {
    xlate = ''
    const parts = $gps.split(/\s*,\s*/)
    const lat = parseFloat(parts[0])
    const long = parts.length > 1 ? parseFloat(parts[1]) : NaN
    if (!isNaN(lat) && !isNaN(long)) {
      xlate = `(${translateAngle(lat, 'N', 'S')}, ${translateAngle(long, 'E', 'W')})`
    }
  }

  function translateAngle(angle: number, pos: string, neg: string): string {
    const suff = angle < 0 ? neg : pos
    angle = Math.abs(angle)
    const deg = Math.floor(angle)
    angle = (angle - deg) * 60
    const min = Math.floor(angle)
    angle = (angle - min) * 60
    const sec = Math.round(angle * 100) / 100
    return `${deg}°${min}'${sec}"${suff}`
  }
</script>

<Label id="gps" label="GPS Coordinates" />
<TextInput id="gps" bind:this={textinput} bind:value={$gps} dirty={$gps !== $image.GPS} />
<div id="xlate">{xlate}</div>
{#if !$gps && $prevImage && $prevImage.GPS}
  <Hint
    on:click={() => {
      $gps = $prevImage.GPS
      tick().then(textinput.focus)
    }}>{$prevImage.GPS}</Hint
  >
{/if}

<style>
  #xlate {
    color: #888;
  }
</style>
