<script lang="ts">
  export let id: string
  export let value: string
  export let dirty: boolean = false

  let inputelm: HTMLInputElement

  export function focus() {
    inputelm.scrollTo(inputelm.scrollWidth, 0)
    inputelm.setSelectionRange(inputelm.value.length, inputelm.value.length)
    inputelm.focus()
  }
</script>

<input {id} name={id} type="text" class:dirty bind:this={inputelm} bind:value />

<style>
  input {
    display: block;
    width: 100%;
    font-family: Arial, Helvetica, sans-serif;
    font-size: 1rem;
  }
  .dirty {
    background-color: orange;
  }
</style>
