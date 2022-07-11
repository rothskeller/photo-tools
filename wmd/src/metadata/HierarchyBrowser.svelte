<script lang="ts">
  import { onMount } from 'svelte'
  import HierarchyItem from './HierarchyItem.svelte'
  import type { Hier } from '../stores'

  export let hierarchy: Hier[]

  let top: HTMLDivElement

  onMount(() => {
    let available = document.body.offsetHeight - 98
    let openToLevel = countLevels(hierarchy)
    while (openToLevel > 0 && 16 * countVisible(hierarchy, openToLevel, 0) > available) {
      openToLevel--
    }
    doOpenToLevel(hierarchy, openToLevel, 0)
  })

  function countLevels(hierarchy: Hier[]): number {
    let count = 1
    hierarchy.forEach((h) => {
      if (h.Children) {
        let subcount = countLevels(h.Children)
        if (subcount >= count) count = subcount + 1
      }
    })
    return count
  }

  function countVisible(hierarchy: Hier[], openToLevel: number, level: number): number {
    let count = hierarchy.length
    hierarchy.forEach((h) => {
      if (h.Children) {
        if (h.Open || level < openToLevel) {
          count += countVisible(h.Children, openToLevel, level + 1)
        }
      }
    })
    return count
  }

  function doOpenToLevel(hierarchy: Hier[], openToLevel: number, level: number) {
    hierarchy.forEach((h) => {
      h.Open = true
      if (h.Children && level < openToLevel) {
        doOpenToLevel(h.Children, openToLevel, level + 1)
      }
    })
  }
</script>

<div bind:this={top}>
  {#each hierarchy as item}
    <HierarchyItem {item} on:select />
  {/each}
</div>

<style>
  div {
    position: fixed;
    bottom: 2rem;
    right: calc(20rem + 6px);
    max-height: calc(100% - 4rem);
    max-width: calc(100% - 22rem - 6px);
    border: 1px solid #ccc;
    background-color: #eee;
    padding: 1rem;
    overflow: auto;
  }
</style>
