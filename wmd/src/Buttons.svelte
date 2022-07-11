<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { index, image, images } from './stores'

  const dispatch = createEventDispatcher()

  let backdayDisabled, backDisabled, nextDisabled, nextdayDisabled
  let backdayIndex, nextdayIndex
  $: {
    let today = $image.DateTime.substr(0, 14)
    backdayIndex = nextdayIndex = -1
    for (let i = $index - 1; i >= 0; i--) {
      let day = $images[i].DateTime.substr(0, 14)
      if (day !== today) {
        backdayIndex = i
        break
      }
    }
    for (let i = $index + 1; i < $images.length; i++) {
      let day = $images[i].DateTime.substr(0, 14)
      if (day !== today) {
        nextdayIndex = i
        break
      }
    }
    backdayDisabled = backdayIndex === -1
    backDisabled = $index === 0
    nextDisabled = $index === $images.length - 1
    nextdayDisabled = nextdayIndex === -1
  }

  function submit(newindex) {
    dispatch('submit', newindex)
  }

  function reset() {
    dispatch('reset')
  }
</script>

<div id="buttons">
  <button
    id="backday"
    disabled={backdayDisabled}
    on:click={() => {
      submit(backdayIndex)
    }}>-1d</button
  >
  <button
    id="back"
    disabled={backDisabled}
    on:click={() => {
      submit($index - 1)
    }}>&lt;</button
  >
  <button id="reset" on:click={reset}>Reset</button>
  <button
    id="next"
    disabled={nextDisabled}
    on:click={() => {
      submit($index + 1)
    }}>&gt;</button
  >
  <button
    id="nextday"
    disabled={nextdayDisabled}
    on:click={() => {
      submit(nextdayIndex)
    }}>+1d</button
  >
</div>

<style>
  #buttons {
    padding: 1rem;
    display: flex;
    justify-content: space-between;
  }
</style>
