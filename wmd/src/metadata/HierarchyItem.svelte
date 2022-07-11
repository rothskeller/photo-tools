<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { Hier } from '../stores'
  export let item: Hier
  export let path: string = ''

  let dispatch = createEventDispatcher()
  let itempath = path ? `${path} / ${item.Name}` : item.Name

  function onClick() {
    item.Open = true
    dispatch('select', itempath)
  }
</script>

<!-- svelte-ignore a11y-invalid-attribute -->
<a class="item" href="#" on:click|preventDefault={onClick}>{item.Name}</a>
{#if item.Open && item.Children && item.Children.length}
  <div class="children">
    {#each item.Children as child}
      <svelte:self item={child} path={itempath} on:select />
    {/each}
  </div>
{/if}

<style>
  .item {
    display: block;
    white-space: nowrap;
    cursor: pointer;
    user-select: none;
    margin: 0;
    padding: 0;
    line-height: 1;
    text-decoration: none;
  }
  .children {
    margin-left: 1.5rem;
  }
</style>
