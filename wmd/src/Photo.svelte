<script lang="ts">
  import { onMount } from 'svelte'
  import { filename } from './stores'
</script>

<div>
  <!-- svelte-ignore a11y-missing-attribute -->
  <img id="photo" src={'/' + $filename} />
</div>

<style>
  #photo {
    object-fit: contain;
    max-width: 100%;
    min-width: 100%;
    max-height: 100%;
    min-height: 100%;
  }
</style>
