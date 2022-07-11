import { writable, derived, Writable } from 'svelte/store'

export interface Metadata {
  Filename: string
  Artist: string
  Caption: string
  DateTime: string
  GPS: string
  Groups: string[]
  Keywords: string[]
  Location: string
  People: string[]
  Places: string[]
  Title: string
  Topics: string[]
}

export interface Hier {
  Name: string
  Open?: boolean
  Children: Hier[]
}

export const index = writable(0)
export const images: Writable<Metadata[]> = writable([])
export const image = derived([index, images], ([index, images]) => images[index])
export const filename = derived([image], ([image]) => image?.Filename)
export const prevImage = derived([index, images], ([index, images]) =>
  index > 0 ? images[index - 1] : null
)

// The server sends us pre-built hierarchies for the various hierarchical tags.
// These are writable because we update them as new tags are added.
export const placeHierarchy: Writable<Hier[]> = writable([])

// For each modifiable metadata item, we have a store.
export const artist = writable('')
export const caption = writable('')
export const gps = writable('')
export const groups: Writable<string[]> = writable([])
export const keywords: Writable<string[]> = writable([])
export const location = writable('')
export const people: Writable<string[]> = writable([])
export const places: Writable<string[]> = writable([])
export const title = writable('')
export const topics: Writable<string[]> = writable([])

// closeHiererchy is a recursive function that sets the Open flag for all items
// in a hierarchy to false.
export function closeHierarchy(hier: Hier[]) {
  hier.forEach((h) => {
    h.Open = false
    if (h.Children) closeHierarchy(h.Children)
  })
}

// openHierarchyPath sets the Open flag to true for the items in the specified
// path.
export function openHierarchyPath(hier: Hier[], path: string) {
  const parts = path.split(/\s*\/\s*/)
  for (let i = 0; i < parts.length; i++) {
    const match = hier.find((h) => h.Name === parts[i])
    if (!match) return
    match.Open = true
    if (match.Children) hier = match.Children
    else return
  }
}
