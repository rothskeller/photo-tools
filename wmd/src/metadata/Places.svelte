<script lang="ts">
  import { onDestroy, tick } from 'svelte'
  import HierarchyBrowser from './HierarchyBrowser.svelte'
  import Hint from './controls/Hint.svelte'
  import Label from './controls/Label.svelte'
  import TextArea from './controls/TextArea.svelte'
  import {
    places,
    image,
    prevImage,
    placeHierarchy,
    closeHierarchy,
    openHierarchyPath,
  } from '../stores'

  // top is the reference to our top-level div.  It is used to test whether the
  // focus belongs to a component descended from that div.
  let top: HTMLDivElement

  // focused is the TextArea component that has focus.  Note that when a link in
  // the hierarchy browser is clicked, it gets focus momentarily, but this
  // variable carefully ignores that and remains set to the TextArea.
  let focused: TextArea

  // inputs is the list of strings corresponding to the TextArea controls.
  // Generally, there is one more input than there are elements in $places.
  let inputs: string[] = []

  // refs is the list of references to the TextArea controls.
  let refs: TextArea[] = []

  // When $places changes, reset inputs to contains the same strings as $places
  // plus a blank string.  Also adjust refs to have the same length as inputs,
  // preserving existing values.
  onDestroy(
    places.subscribe((places) => {
      inputs = [...places, '']
      if (refs.length > inputs.length) refs = refs.slice(0, inputs.length)
      while (refs.length < inputs.length) refs = [...refs, null]
    })
  )

  // When the last input becomes non-empty, add another.
  $: if (inputs[inputs.length - 1]) {
    inputs = [...inputs, '']
    refs.push(null)
  }

  // When a TextArea indicates that it has received focus, and it's not the one
  // we think has focus already, update our idea of the focus and reset the open
  // flags in the hierarchy browser so that it's open to the value in that
  // TextArea.
  function onFocus(index: number) {
    const target = refs[index]
    if (focused === target) return
    focused = target
    closeHierarchy($placeHierarchy)
    openHierarchyPath($placeHierarchy, inputs[index])
  }

  // When any descendant of this component loses focus, check to see where the
  // focus is moving to.  If it's moving outside of this component, set focused
  // to null so that the hierarchy browser stops displaying.
  function onFocusOut(event: FocusEvent) {
    let next = event.relatedTarget as Element
    while (next) {
      if (next === top) return
      next = next.parentElement
    }
    focused = null
    places.set(inputs.filter((input) => !!input))
  }

  // When an item in the hierarchy browser is selected, set the corresponding
  // input, then re-focus on the corresponding TextArea.
  function onSelect(event: CustomEvent<string>) {
    const index = refs.findIndex((ref) => ref === focused)
    if (index >= 0) {
      inputs[index] = event.detail
      tick().then(focused.focus)
    }
  }
</script>

<div bind:this={top} on:focusout={onFocusOut}>
  <Label id="place0" label="Places" />
  {#each inputs as input, i}
    <TextArea
      id={`place${i}`}
      bind:this={refs[i]}
      bind:value={input}
      dirty={input !== (i < $image.Places.length ? $image.Places[i] : '')}
      on:focus={() => {
        onFocus(i)
      }}
    />
    {#if !input && $prevImage && $prevImage.Places.length > i}
      <Hint
        on:click={() => {
          input = $prevImage.Places[i]
          tick().then(refs[i].focus)
        }}>{$prevImage.Places[i]}</Hint
      >
    {/if}
  {/each}
  {#if focused}
    <HierarchyBrowser hierarchy={$placeHierarchy} on:select={onSelect} />
  {/if}
</div>
