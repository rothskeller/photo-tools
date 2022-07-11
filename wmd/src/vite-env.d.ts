/// <reference types="svelte" />
/// <reference types="vite/client" />

interface WindowOrWorkerGlobalScope {
  structuredClone(value: any, options?: StructuredSerializeOptions): any
}
declare function structuredClone(value: any, options?: StructuredSerializeOptions): any
