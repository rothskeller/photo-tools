<script lang="ts">
  export let id: string
  export let value: string
  export let dirty: boolean = false

  let textarea: HTMLTextAreaElement

  export function focus() {
    textarea.setSelectionRange(textarea.value.length, textarea.value.length)
    textarea.focus()
  }
</script>

<div>
  <textarea {id} name={id} rows="1" class:dirty bind:this={textarea} bind:value on:focus on:blur />
  <span>{value}</span>
</div>

<style>
  div {
    display: grid;
  }
  textarea {
    resize: none;
    overflow: hidden;
  }
  span {
    white-space: pre-wrap;
    visibility: hidden;
    border: 1px solid transparent;
  }
  textarea,
  span {
    grid-area: 1 / 1 / 2 / 2;
    display: block;
    width: 100%;
    font-family: Arial, Helvetica, sans-serif;
    font-size: 1rem;
    padding: 0.125rem;
  }
  .dirty {
    background-color: orange;
  }
</style>
