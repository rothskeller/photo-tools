<script lang="ts">
  import Photo from './Photo.svelte'
  import Map from './Map.svelte'
  import Metadata from './Metadata.svelte'
  import Buttons from './Buttons.svelte'
  import {
    Hier,
    artist,
    caption,
    filename,
    gps,
    groups,
    image,
    images,
    index,
    keywords,
    location,
    people,
    places,
    placeHierarchy,
    title,
    topics,
  } from './stores'

  const dataload = fetch('/metadata.json')
    .then((resp) => resp.json())
    .then((data) => {
      images.set(data.images)
      placeHierarchy.set(data.placeHierarchy)
    })

  function reset() {
    // sets all metadata stores based on selected image
    if (!$image) return
    $artist = $image.Artist
    $caption = $image.Caption
    $gps = $image.GPS
    $groups = [...$image.Groups]
    $keywords = [...$image.Keywords]
    $location = $image.Location
    $people = [...$image.People]
    $places = [...$image.Places]
    $title = $image.Title
    $topics = [...$image.Topics]
  }
  $: {
    $image // reference $image so that this will get triggered when it changes
    reset()
  }

  function stringListsEqual(a, b) {
    // true if two arrays contain same strings
    if (a.length != b.length) return false
    for (let i = 0; i < a.length; i++) if (a[i] !== b[i]) return false
    return true
  }

  function dirty() {
    // true if any metadata have been changed
    return (
      $artist !== $image.Artist ||
      $caption !== $image.Caption ||
      $gps !== $image.GPS ||
      !stringListsEqual($groups, $image.Groups) ||
      !stringListsEqual($keywords, $image.Keywords) ||
      $location !== $image.Location ||
      !stringListsEqual($people, $image.People) ||
      !stringListsEqual($places, $image.Places) ||
      $title !== $image.Title ||
      !stringListsEqual($topics, $image.Topics)
    )
  }

  async function submit(event) {
    if (dirty()) {
      await save()
      updateHierarchies()
    }
    $index = event.detail
  }

  async function save() {
    const body = new FormData()
    body.append('artist', $artist)
    body.append('caption', $caption)
    body.append('title', $title)
    body.append('gps', $gps)
    $places.forEach((p) => {
      body.append('places', p)
    })
    if (!$places.length) body.append('places', '')
    const resp = await fetch(`/${$filename}`, { method: 'POST', body })
    const result = await resp.json()
    $images[$index] = result
  }

  function updateHierarchies() {
    $image.Places.forEach((place) => {
      updateHierarchy($placeHierarchy, place)
    })
    console.log($placeHierarchy)
  }

  function updateHierarchy(hier: Hier[], value: string) {
    const parts = value.split(/\s*\/\s*/)
    let sub = hier.find((n) => n.Name === parts[0])
    if (!sub) {
      sub = { Name: parts[0], Children: null }
      hier.push(sub)
      hier.sort((a: Hier, b: Hier) => {
        if (a.Name < b.Name) return -1
        if (a.Name > b.Name) return +1
        return 0
      })
    }
    if (parts.length > 1) {
      if (!sub.Children) sub.Children = []
      updateHierarchy(sub.Children, parts.slice(1).join('/'))
    }
  }
</script>

<main>
  {#await dataload then unused}
    <div id="photo-map">
      <Photo />
      <div class="divider" />
      <Map />
    </div>
    <div class="divider" />
    <form id="metadata-buttons" on:submit|preventDefault={() => {}}>
      <Metadata />
      <div class="divider" />
      <Buttons on:submit={submit} on:reset={reset} />
    </form>
  {/await}
</main>

<style>
  :global(body) {
    margin: 0 !important;
    font-family: Arial, Helvetica, sans-serif;
  }

  main {
    width: 100vw;
    display: grid;
    grid: 100vh / 1fr 6px 20rem;
  }

  #photo-map {
    position: relative;
    display: grid;
    grid: calc(50% - 3px) 6px calc(50% - 3px) / 100%;
  }

  #metadata-buttons {
    display: grid;
    grid: 1fr 6px max-content / 100%;
  }

  .divider {
    background-color: #ccc;
  }
</style>
