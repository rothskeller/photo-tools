<script lang="ts">
  import { filename, image } from './stores'
  import Artist from './metadata/Artist.svelte'
  import Caption from './metadata/Caption.svelte'
  import GPS from './metadata/GPS.svelte'
  import Places from './metadata/Places.svelte'
  import Title from './metadata/Title.svelte'
</script>

<div id="metadata">
  <div id="filename">{$filename}</div>
  <div id="datetime">{$image.DateTime}</div>
  <Title />
  <Caption />
  <Artist />
  <GPS />
  <Places />
</div>

<style>
  #metadata {
    padding: 1rem;
    overflow-y: auto;
  }
  #filename {
    color: #444;
    font-size: 1.25rem;
    font-weight: bold;
  }
  #datetime {
    color: #888;
  }
</style>
